package model

// 文章点赞表
type ArticleLike struct {
	BaseModel
	UserID    uint `gorm:"not null;index:idx_user_article_like,unique"`
	ArticleID uint `gorm:"not null;index:idx_user_article_like,unique"`
}
