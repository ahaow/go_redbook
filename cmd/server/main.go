package main

import (
	"go_redbook/config"
	"go_redbook/global"
	"go_redbook/internal/database"
	"go_redbook/internal/logger"
	"go_redbook/internal/router"
	"go_redbook/internal/svc"
)

func main() {

	// 1. 获取配置信息
	cfg, err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	global.Config = cfg

	// 2. 初始化日志
	log, err := logger.NewLogger("development", "logs", "[myApp] ")
	if err != nil {
		panic(err)
	}
	global.Log = log
	defer log.Close()

	// 3. 初始化gorm
	global.DB = database.InitDB(cfg.Database)
	if err := database.AutoMigrate(global.DB); err != nil {
		panic(err)
	}

	// 4. 初始化依赖容器和路由
	svcCtx := svc.NewServiceContext(cfg, global.DB)
	r := router.InitRouter(svcCtx)

	port := cfg.App.Port
	if port == "" {
		port = ":3000"
	}
	if err := r.Run(port); err != nil {
		panic(err)
	}
}
