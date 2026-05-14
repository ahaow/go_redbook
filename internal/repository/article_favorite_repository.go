package repository

import (
	"context"
	"go_redbook/internal/model"

	"gorm.io/gorm"
)

type ArticleFavoriteRepository interface {
	Create(ctx context.Context, favorite *model.ArticleFavorite) error
	Delete(ctx context.Context, userID, articleID uint) error
	Exists(ctx context.Context, userID, articleID uint) (bool, error)
	CountByArticleID(ctx context.Context, articleID uint) (int64, error)
}

type articleFavoriteRepository struct {
	db *gorm.DB
}

func NewArticleFavoriteRepository(db *gorm.DB) ArticleFavoriteRepository {
	return &articleFavoriteRepository{db: db}
}

func (r *articleFavoriteRepository) Create(ctx context.Context, favorite *model.ArticleFavorite) error {
	return r.db.WithContext(ctx).Create(favorite).Error
}

func (r *articleFavoriteRepository) Delete(ctx context.Context, userID, articleID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Delete(&model.ArticleFavorite{}).Error
}

func (r *articleFavoriteRepository) Exists(ctx context.Context, userID, articleID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.ArticleFavorite{}).
		Where("user_id = ? AND article_id = ?", userID, articleID).
		Count(&count).Error
	return count > 0, err
}

func (r *articleFavoriteRepository) CountByArticleID(ctx context.Context, articleID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.ArticleFavorite{}).
		Where("article_id = ?", articleID).
		Count(&count).Error
	return count, err
}
