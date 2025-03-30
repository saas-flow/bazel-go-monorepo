package logger

import (
	"regexp"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func maskName(name string) string {
	if len(name) <= 2 {
		return name[:1] + "*"
	}
	return name[:1] + strings.Repeat("*", len(name)-1)
}

// Masking function for sensitive data
func maskValue(key, value string) string {

	if key == "user.name.first" || key == "user.name.last" {
		return maskName(value)
	}

	if key == "user.email" {
		emailRegex := regexp.MustCompile(`([a-zA-Z0-9._%+-]+)@([a-zA-Z0-9.-]+\.[a-zA-Z]{2,})`)
		return emailRegex.ReplaceAllString(value, "***@***.com")
	}

	if key == "user.phone" {
		phoneRegex := regexp.MustCompile(`(\+\d{2,3})\d{4,7}(\d{3})`)
		return phoneRegex.ReplaceAllString(value, "$1*******$2")
	}

	if key == "credit_card" {
		ccRegex := regexp.MustCompile(`\b(\d{4})[-\s]?(\d{4})[-\s]?(\d{4})[-\s]?(\d{4})\b`)
		return ccRegex.ReplaceAllString(value, "**** **** **** $4")
	}

	sensitiveFields := []string{"user.password", "token", "cvv", "session_id"}
	for _, field := range sensitiveFields {
		if key == field {
			return "******"
		}
	}
	return value
}

type MaskingCore struct {
	zapcore.Core
}

func (c *MaskingCore) With(fields []zapcore.Field) zapcore.Core {
	maskedFields := make([]zapcore.Field, len(fields))
	for i, field := range fields {
		switch field.Type {
		case zapcore.StringType:
			maskedFields[i] = zap.String(field.Key, maskValue(field.Key, field.String))
		// case zapcore.ArrayMarshalerType, zapcore.ObjectMarshalerType:
		// 	maskedFields[i] = zap.Any(field.Key, "******") // Masking untuk objek/array
		default:
			maskedFields[i] = field
		}
	}
	return &MaskingCore{c.Core.With(maskedFields)}
}
