package model

// 文章收藏表
type ArticleFavorite struct {
	BaseModel
	UserID    uint `gorm:"not null;index:idx_user_article_favorite,unique"`
	ArticleID uint `gorm:"not null;index:idx_user_article_favorite,unique"`
}
