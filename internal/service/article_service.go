package service

import (
	"context"
	"errors"
	"fmt"
	"go_redbook/internal/model"
	"go_redbook/internal/repository"

	"gorm.io/gorm"
)

type CreateArticleInput struct {
	UserID    uint
	Title     string
	Content   string
	ImageURLs []string
	Topics    []string
	Location  string
	IsPublic  bool
}

type UpdateArticleInput struct {
	CurrentUserID uint
	ArticleID     uint
	Title         *string
	Content       *string
	Location      *string
	IsPublic      *bool
	Images        *[]string
	Topics        *[]string
}

type ArticleList struct {
	Items []model.Article
	Total int64
}

type ArticleService interface {
	Create(ctx context.Context, input CreateArticleInput) (*model.Article, error)
	GetByID(ctx context.Context, id uint) (*model.Article, error)
	List(ctx context.Context, page, pageSize int) (*ArticleList, error)
	DeleteById(ctx context.Context, currentUserID, articleID uint) error
	Update(ctx context.Context, input UpdateArticleInput) (*model.Article, error)
	Like(ctx context.Context, userID, articleID uint) error
	Unlike(ctx context.Context, userID, articleID uint) error
	Favorite(ctx context.Context, userID, articleID uint) error
	Unfavorite(ctx context.Context, userID, articleID uint) error
}

type articleService struct {
	articleRepo         repository.ArticleRepository
	articleLikeRepo     repository.ArticleLikeRepository
	articleFavoriteRepo repository.ArticleFavoriteRepository
}

func NewArticleService(
	articleRepo repository.ArticleRepository,
	articleLikeRepo repository.ArticleLikeRepository,
	articleFavoriteRepo repository.ArticleFavoriteRepository,
) ArticleService {
	return &articleService{
		articleRepo:         articleRepo,
		articleLikeRepo:     articleLikeRepo,
		articleFavoriteRepo: articleFavoriteRepo,
	}
}

// 创建文章
func (s *articleService) Create(ctx context.Context, input CreateArticleInput) (*model.Article, error) {
	article := &model.Article{
		UserID:   input.UserID,
		Title:    input.Title,
		Content:  input.Content,
		Location: input.Location,
		IsPublic: input.IsPublic,
	}

	for i, url := range input.ImageURLs {
		article.Images = append(article.Images, model.ArticleImage{
			URL:       url,
			SortOrder: i,
		})
	}

	var topics []model.Topic
	for _, name := range input.Topics {
		if name == "" {
			continue
		}

		topic, err := s.articleRepo.FindOrCreateTopic(ctx, name)
		if err != nil {
			return nil, err
		}
		topics = append(topics, *topic)
	}

	article.Topics = topics

	if err := s.articleRepo.Create(ctx, article); err != nil {
		return nil, err
	}

	return article, nil
}

// 文章列表
func (s *articleService) List(ctx context.Context, page, pageSize int) (*ArticleList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	articles, total, err := s.articleRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}
	return &ArticleList{Items: articles, Total: total}, nil
}

// 根据id获取某个文章
func (s *articleService) GetByID(ctx context.Context, id uint) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, fmt.Errorf("article: %w", ErrNotFound)
	}
	return article, nil
}

// 根据id删除某个文章
func (s *articleService) DeleteById(ctx context.Context, currentUserId, articleID uint) error {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		return err
	}
	if article == nil {
		return fmt.Errorf("article: %w", ErrNotFound)
	}

	if article.UserID != currentUserId {
		return fmt.Errorf("article: %w", ErrForbidden)
	}

	err = s.articleRepo.DeleteByID(ctx, articleID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("article: %w", ErrNotFound)
		}
		return err
	}
	return nil
}

// 更新某个文章
func (s *articleService) Update(ctx context.Context, input UpdateArticleInput) (*model.Article, error) {
	article, err := s.articleRepo.FindByID(ctx, input.ArticleID)
	if err != nil {
		return nil, err
	}
	if article == nil {
		return nil, fmt.Errorf("article: %w", ErrNotFound)
	}
	if article.UserID != input.CurrentUserID {
		return nil, fmt.Errorf("article: %w", ErrForbidden)
	}

	if input.Title != nil {
		article.Title = *input.Title
	}
	if input.Content != nil {
		article.Content = *input.Content
	}
	if input.Location != nil {
		article.Location = *input.Location
	}
	if input.IsPublic != nil {
		article.IsPublic = *input.IsPublic
	}
	if err := s.articleRepo.Update(ctx, article); err != nil {
		return nil, err
	}
	if input.Images != nil {
		var images []model.ArticleImage
		for i, url := range *input.Images {
			images = append(images, model.ArticleImage{
				ArticleID: article.ID,
				URL:       url,
				SortOrder: i,
			})
		}
		if err := s.articleRepo.ReplaceImages(ctx, article.ID, images); err != nil {
			return nil, err
		}
	}

	if input.Topics != nil {
		var topics []model.Topic
		for _, name := range *input.Topics {
			if name == "" {
				continue
			}
			topic, err := s.articleRepo.FindOrCreateTopic(ctx, name)
			if err != nil {
				return nil, err
			}
			topics = append(topics, *topic)
		}
		if err := s.articleRepo.ReplaceTopics(ctx, article, topics); err != nil {
			return nil, err
		}
	}

	return article, nil
}

// 点赞
func (s *articleService) Like(ctx context.Context, userID, articleID uint) error {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		return err
	}
	if article == nil {
		return fmt.Errorf("article: %w", ErrNotFound)
	}

	exists, err := s.articleLikeRepo.Exists(ctx, userID, articleID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return s.articleLikeRepo.Create(ctx, &model.ArticleLike{
		UserID:    userID,
		ArticleID: articleID,
	})
}

// 取消点赞
func (s *articleService) Unlike(ctx context.Context, userID, articleID uint) error {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		return err
	}
	if article == nil {
		return fmt.Errorf("article: %w", ErrNotFound)
	}
	return s.articleLikeRepo.Delete(ctx, userID, articleID)
}

// 收藏
func (s *articleService) Favorite(ctx context.Context, userID, articleID uint) error {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		return err
	}
	if article == nil {
		return fmt.Errorf("article: %w", ErrNotFound)
	}

	exists, err := s.articleFavoriteRepo.Exists(ctx, userID, articleID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	return s.articleFavoriteRepo.Create(ctx, &model.ArticleFavorite{
		UserID:    userID,
		ArticleID: articleID,
	})
}

// 取消收藏
func (s *articleService) Unfavorite(ctx context.Context, userID, articleID uint) error {
	article, err := s.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		return err
	}
	if article == nil {
		return fmt.Errorf("article: %w", ErrNotFound)
	}
	return s.articleFavoriteRepo.Delete(ctx, userID, articleID)
}
