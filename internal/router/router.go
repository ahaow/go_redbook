package router

import (
	"time"

	"go_redbook/config"
	"go_redbook/internal/handler"
	"go_redbook/internal/repository"
	"go_redbook/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitRouter 初始化 Gin 引擎和全局中间件。
func InitRouter(db *gorm.DB, jwtCfg config.JwtConfig) *gin.Engine {
	gin.SetMode("debug")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000"}, // 你前端的地址，*表示全部允许
		AllowOrigins:     []string{"*"}, // 你前端的地址，*表示全部允许
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/api/v1")
	registerUserRoutes(api, db, jwtCfg)

	return r
}

// registerUserRoutes 组装用户模块的依赖并注册路由。
// 依赖方向：handler -> service -> repository -> db。
func registerUserRoutes(api *gin.RouterGroup, db *gorm.DB, jwtCfg config.JwtConfig) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtCfg)
	userHandler := handler.NewUserHandler(userService)

	userHandler.RegisterRoutes(api)
}
