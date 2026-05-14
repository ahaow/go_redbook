package router

import (
	"time"

	"go_redbook/internal/handler"
	"go_redbook/internal/middleware"
	"go_redbook/internal/svc"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化 Gin 引擎和全局中间件。
func InitRouter(svcCtx *svc.ServiceContext) *gin.Engine {
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
	registerUserRoutes(api, svcCtx)
	registerPublicArticleRoutes(api, svcCtx)

	auth := api.Group("")
	auth.Use(middleware.JWTAuth(svcCtx.Config.Jwt))
	registerPrivateArticleRoutes(auth, svcCtx)

	return r
}

// registerUserRoutes 组装用户模块的依赖并注册路由。
// 依赖方向：handler -> service -> repository -> db。
func registerUserRoutes(api *gin.RouterGroup, svcCtx *svc.ServiceContext) {
	userHandler := handler.NewUserHandler(svcCtx.UserService)
	userHandler.RegisterRoutes(api)
}

func registerPublicArticleRoutes(api *gin.RouterGroup, svcCtx *svc.ServiceContext) {
	articleHandler := handler.NewArticleHandler(svcCtx.ArticleService)
	articleHandler.RegisterPublicRoutes(api)
}

func registerPrivateArticleRoutes(api *gin.RouterGroup, svcCtx *svc.ServiceContext) {
	articleHandler := handler.NewArticleHandler(svcCtx.ArticleService)
	articleHandler.RegisterCreateRoutes(api)
	articleHandler.RegisterPrivateRoutes(api)
}
