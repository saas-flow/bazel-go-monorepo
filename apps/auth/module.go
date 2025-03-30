package main

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisSession "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
	delivery_http "github.com/saas-flow/monorepo/apps/auth/internal/user/delivery/http"
	"github.com/saas-flow/monorepo/apps/auth/internal/user/repository"
	"github.com/saas-flow/monorepo/apps/auth/internal/user/usecase"
	"github.com/saas-flow/monorepo/libs/config"
	"github.com/saas-flow/monorepo/libs/database"
	"github.com/saas-flow/monorepo/libs/httpserver"
	"github.com/saas-flow/monorepo/libs/logger"
	"github.com/saas-flow/monorepo/libs/middleware"
	"github.com/saas-flow/monorepo/libs/otelcol"
	redisClient "github.com/saas-flow/monorepo/libs/redis"
	"github.com/saas-flow/monorepo/libs/validator"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func RegisterDatabaseDriver() gorm.Dialector {
	return postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.GetString("DATABASE.HOST"),
		config.GetString("DATABASE.USER"),
		config.GetString("DATABASE.PASSWORD"),
		config.GetString("DATABASE.DBNAME"),
		config.GetString("DATABASE.PORT"),
		config.GetString("DATABASE.SSLMODE"),
		config.GetString("DATABASE.TIMEZONE"),
	))
}

func RegisterFxLogger(logger *zap.Logger) fxevent.Logger {
	return &fxevent.ZapLogger{Logger: logger}
}

func RegisterRedisOption() *redis.Options {
	return &redis.Options{
		Addr: config.GetString("REDIS.ADDR"),
		DB:   config.GetInt("REDIS.DB"),
	}
}

func RegisterGin(
	translate ut.Translator,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(otelgin.Middleware(config.GetString("SERVICE_NAME")))
	r.Use(middleware.Logging())
	r.Use(middleware.MiddlewareError(translate))

	store, err := redisSession.NewStore(10, "tcp", config.GetString("REDIS.ADDR"), "", []byte(config.GetString("SESSION.SECRET")))
	if err != nil {
		zap.L().Panic(err.Error())
		panic(err)
	}

	r.Use(sessions.Sessions(config.GetString("SESSION.NAME"), store))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetStringSlice("CORS_ALLOW_ORIGINS"), // Ganti dengan domain frontend kamu
		AllowMethods:     config.GetStringSlice("CORS_ALLOW_METHODS"),
		AllowHeaders:     config.GetStringSlice("CORS_ALLOW_HEADERS"),
		ExposeHeaders:    []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	return r
}

func RegisterHttpServerConfig(r *gin.Engine) *httpserver.Config {
	return &httpserver.Config{
		Addr:    config.GetString("APP.PORT"),
		Handler: r.Handler(),
	}
}

var Module = fx.Module("auth",
	logger.Module,
	otelcol.Resource,
	otelcol.MetricProvider,
	otelcol.TraceProvider,
	fx.Provide(RegisterDatabaseDriver, RegisterRedisOption),
	database.Module,
	redisClient.Module,
	validator.Module,
	validator.TranslationModule,
	fx.Provide(
		RegisterGin,
		RegisterHttpServerConfig,
	),
	repository.Module,
	usecase.AuthUsecaseModule,
	delivery_http.AuthRouter,
	httpserver.Module,
	fx.WithLogger(RegisterFxLogger),
)
