package model

import (
	"time"

	"gorm.io/gorm"
)

// Article 是数据库里的文章表模型。
type Article struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"index;not null"` // 发布人id
	User   User `gorm:"foreignKey:UserID"`

	Title     string         `gorm:"type:varchar(100);not null"` // 标题
	Content   string         `gorm:"type:text;not null"`         // 内容
	Location  string         `gorm:"type:varchar(255)"`          // 位置
	IsPublic  bool           `gorm:"default:true;not null"`      // 是否公开
	Images    []ArticleImage `gorm:"foreignKey:ArticleID"`
	Topics    []Topic        `gorm:"many2many:article_topics;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
