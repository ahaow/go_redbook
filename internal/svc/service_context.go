package svc

import (
	"go_redbook/config"
	"go_redbook/internal/repository"
	"go_redbook/internal/service"

	"gorm.io/gorm"
)

// ServiceContext 集中管理项目依赖。
// 新增模块时，在这里完成 repo -> service 的接线。
type ServiceContext struct {
	Config *config.Config
	db     *gorm.DB

	UserService    service.UserService
	ArticleService service.ArticleService
}

// NewServiceContext 根据配置和基础设施创建依赖容器。
func NewServiceContext(cfg *config.Config, db *gorm.DB) *ServiceContext {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.Jwt)

	articleRepo := repository.NewArticleRepository(db)
	articleLikeRepo := repository.NewArticleLikeRepository(db)
	articleFavoriteRepo := repository.NewArticleFavoriteRepository(db)
	articleService := service.NewArticleService(articleRepo, articleLikeRepo, articleFavoriteRepo)

	return &ServiceContext{
		Config:         cfg,
		UserService:    userService,
		ArticleService: articleService,
		db:             db,
	}
}
