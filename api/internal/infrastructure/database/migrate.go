package database

import (
	"github.com/unitechio/oss-monorepo/api/internal/domain"
	"github.com/unitechio/oss-monorepo/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	logger.Info("Running database migrations...")

	if err := db.AutoMigrate(
		&domain.Page{},
		&domain.Block{},
	); err != nil {
		logger.Error("Failed to migrate page builder tables", zap.Error(err))
		return err
	}
	logger.Info("Database migrations completed successfully")
	return nil
}
