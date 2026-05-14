package repository

import (
	"context"
	"errors"
	"go_redbook/internal/model"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(ctx context.Context, article *model.Article) error
	DeleteByID(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*model.Article, error)
	FindOrCreateTopic(ctx context.Context, name string) (*model.Topic, error)
	List(ctx context.Context, offset, limit int) ([]model.Article, int64, error)
	Update(ctx context.Context, article *model.Article) error
	ReplaceImages(ctx context.Context, articleID uint, images []model.ArticleImage) error
	ReplaceTopics(ctx context.Context, article *model.Article, topics []model.Topic) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

func (r *articleRepository) DeleteByID(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&model.Article{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *articleRepository) FindByID(ctx context.Context, id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Images").
		Preload("Topics").
		First(&article, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) FindOrCreateTopic(ctx context.Context, name string) (*model.Topic, error) {
	var topic model.Topic

	err := r.db.WithContext(ctx).
		Where("name = ?", name).
		First(&topic).Error

	if err == nil {
		return &topic, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	topic = model.Topic{Name: name}
	if err := r.db.WithContext(ctx).Create(&topic).Error; err != nil {
		return nil, err
	}
	return &topic, nil
}

func (r *articleRepository) List(ctx context.Context, offset, limit int) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64
	db := r.db.WithContext(ctx).Model(&model.Article{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Preload("User").Preload("Images").Preload("Topics").Order("id DESC").Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

func (r *articleRepository) Update(ctx context.Context, article *model.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

func (r *articleRepository) ReplaceImages(ctx context.Context, articleID uint, images []model.ArticleImage) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("article_id = ?", articleID).Delete(&model.ArticleImage{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if len(images) > 0 {
		for i := range images {
			images[i].ArticleID = articleID
		}
		if err := tx.Create(&images).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (r *articleRepository) ReplaceTopics(ctx context.Context, article *model.Article, topics []model.Topic) error {
	return r.db.WithContext(ctx).Model(article).Association("Topics").Replace(topics)
}
