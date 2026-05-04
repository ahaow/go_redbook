package database

import (
	"log"
	"os"
	"time"

	"go_redbook/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(cfg config.DatabaseConfig) *gorm.DB {
	dsn := cfg.Dsn

	// 初始化 GORM 日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值，可从 cfg.SlowSQL 设置
			LogLevel:      logger.Info, // 日志等级，可从 cfg.LogLevel 设置
			Colorful:      true,
		},
	)
	// 打开数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// 获取底层 *sql.DB 并设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns) // 最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenCons)  // 最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour)     // 每个连接最大生命周期，可配置化

	return db
}
