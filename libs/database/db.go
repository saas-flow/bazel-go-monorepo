package database

import (
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var Module = fx.Module("database",
	fx.Provide(
		New,
	),
)

func New(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
	// Initialize the GORM DB connection
	db, err := gorm.Open(dialector, opts...)
	if err != nil {
		zap.L().Fatal("❌ Failed to connect to database", zap.Error(err))
		return nil, err
	}

	// Register the OpenTelemetry plugin with GORM
	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		zap.L().Fatal("❌ Failed to use OpenTelemetry plugin", zap.Error(err))
		return nil, err
	}

	if os.Getenv("ENV") != "production" {
		db = db.Debug()
		zap.L().Info("🔍 Database is running in DEBUG mode")
	}

	zap.L().Info("✅ Database connection successfully configured with connection pooling.")

	return db, nil
}
