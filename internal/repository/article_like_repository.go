package repository

import (
	"context"
	"go_redbook/internal/model"

	"gorm.io/gorm"
)

type ArticleLikeRepository interface {
	Create(ctx context.Context, like *model.ArticleLike) error
	Delete(ctx context.Context, userID, articleID uint) error
	Exists(ctx context.Context, userID, articleID uint) (bool, error)
	CountByArticleID(ctx context.Context, articleID uint) (int64, error)
}

type articleLikeRepository struct {
	db *gorm.DB
}

func NewArticleLikeRepository(db *gorm.DB) ArticleLikeRepository {
	return &articleLikeRepository{db: db}
}

// 创建点赞
func (r *articleLikeRepository) Create(ctx context.Context, like *model.ArticleLike) error {
	return r.db.WithContext(ctx).Create(like).Error
}

// 判断是否已点赞
func (r *articleLikeRepository) Exists(ctx context.Context, userID, articleID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.ArticleLike{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error
	return count > 0, err
}

// 取消点赞
func (r *articleLikeRepository) Delete(ctx context.Context, userID, articleID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Delete(&model.ArticleLike{}).Error
}

// 统计点赞数
func (r *articleLikeRepository) CountByArticleID(ctx context.Context, articleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.ArticleLike{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}
