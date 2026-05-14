package database

import (
	"go_redbook/internal/model"

	"gorm.io/gorm"
)

// AutoMigrate 统一维护需要自动迁移的模型。
// 新增 model 后，在这里追加即可。
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Article{},
		&model.ArticleImage{},
		&model.ArticleLike{},
		&model.ArticleFavorite{},
		&model.Topic{},
	)
}
