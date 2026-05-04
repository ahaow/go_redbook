package global

import (
	"go_redbook/config"
	"go_redbook/internal/logger"

	"gorm.io/gorm"
)

var (
	Config *config.Config
	Log    *logger.Logger
	DB     *gorm.DB
)
