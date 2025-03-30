package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// List of endpoints to be excluded from logging
	excludedEndpoints = []string{
		"/metrics",
		"/health/liveness",
		"/health/readiness",
	}
	sensitiveFields = []string{"user", "organization", "credit_card", "cvv", "ssn"}
)

// Function to check if the path is in the list of excluded endpoints
func isExcludedPath(path string) bool {
	for _, endpoint := range excludedEndpoints {
		if strings.HasPrefix(path, endpoint) {
			return true
		}
	}
	return false
}

// Mask email (j***@example.com)
func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***@***.com"
	}
	return parts[0][:1] + "***@" + parts[1]
}

// Mask name (A*** B***)
func maskName(name string) string {
	if len(name) <= 2 {
		return name[:1] + "*"
	}
	return name[:1] + strings.Repeat("*", len(name)-1)
}

// Masking data sensitif
func maskSensitiveData(data map[string]interface{}) map[string]interface{} {
	for _, field := range sensitiveFields {
		if value, exists := data[field]; exists {
			v := reflect.ValueOf(value)
			switch v.Kind() {
			case reflect.Map:
				maskedData := make(map[string]interface{})
				for _, key := range v.MapKeys() {
					keyStr := key.String()
					value := v.MapIndex(key).Interface()
					if keyStr == "email" {
						maskedData[keyStr] = maskEmail(fmt.Sprintf("%v", value))
					} else if keyStr == "first_name" || keyStr == "last_name" || keyStr == "full_name" {
						maskedData[keyStr] = maskName(value.(string))
					} else {
						maskedData[keyStr] = "******"
					}
				}
				data[field] = maskedData
			case reflect.Slice:
				data[field] = value
			default:
				data[field] = "******"
			}
		}
	}
	return data
}

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {

		if isExcludedPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		ctx := c.Request.Context()
		header := c.Request.Header
		carrier := propagation.HeaderCarrier(header)

		propgator := propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
		propgator.Inject(c.Request.Context(), carrier)
		c.Request = c.Request.WithContext(propgator.Extract(ctx, carrier))

		start := time.Now()
		fields := []zapcore.Field{}
		if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
			fields = append(fields, zap.String("request_id", requestID))
		}

		if span := trace.SpanContextFromContext(c.Request.Context()); span.IsValid() {
			fields = append(fields, zap.String("trace_id", span.TraceID().String()))
			fields = append(fields, zap.String("span_id", span.SpanID().String()))
		}

		fields = append(fields, zap.String("http_method", c.Request.Method))
		fields = append(fields, zap.String("http_url", c.Request.URL.Path))

		if userId := c.Writer.Header().Get("X-User-ID"); userId != "" {
			fields = append(fields, zap.String("user_id", userId))
		}

		if roles := c.Writer.Header().Get("X-Roles"); roles != "" {
			fields = append(fields, zap.String("roles", roles))
		}

		// Baca & Masking Body Request
		var body []byte
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		body, _ = io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf) // Pastikan body bisa dibaca ulang

		// Handle "application/x-www-form-urlencoded"
		if c.ContentType() == "application/x-www-form-urlencoded" {
			if err := c.Request.ParseForm(); err != nil {
				fields = append(fields, zap.String("error", "failed to parse form data"))
			} else {
				formData := make(map[string]interface{})
				for key, values := range c.Request.PostForm {
					if len(values) > 1 {
						formData[key] = values // Simpan sebagai array jika ada banyak nilai
					} else {
						formData[key] = values[0] // Simpan sebagai string jika hanya satu nilai
					}
				}

				maskedBody := maskSensitiveData(formData)
				fields = append(fields, zap.Any("http_request", maskedBody))
			}
		}

		// Handle "application/json"
		if c.ContentType() == "application/json" {
			tee := io.TeeReader(c.Request.Body, &buf)
			body, _ = io.ReadAll(tee)
			c.Request.Body = io.NopCloser(&buf) // Reset body agar bisa dibaca ulang

			if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
				var newb map[string]interface{}
				if err := json.Unmarshal(body, &newb); err != nil {
					fields = append(fields, zap.String("error", "failed to parse request body"))
				} else {
					maskedBody := maskSensitiveData(newb)
					fields = append(fields, zap.Any("http_request", maskedBody))
				}
			}
		}

		c.Next()

		// Tambahkan durasi request & status HTTP
		fields = append(fields, zap.Duration("http_duration", time.Since(start)))
		fields = append(fields, zap.Int("http_status", c.Writer.Status()))

		// Tangkap error jika ada
		if err := c.Errors.Last(); err != nil {
			fields = append(fields, zap.Any("error", err.Error()))
		}

		if span := trace.SpanContextFromContext(c.Request.Context()); span.IsValid() {
			c.Header("X-Trace-ID", span.TraceID().String())
		}

		zap.L().With(fields...).Info("Incoming Request")
	}
}
