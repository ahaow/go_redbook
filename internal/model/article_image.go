package model

import (
	"time"
)

type ArticleImage struct {
	ID        uint `gorm:"primaryKey"`
	ArticleID uint `gorm:"index;not null"`

	URL       string `gorm:"type:varchar(500);not null"`
	SortOrder int    `gorm:"default:0"` // 多图排序

	CreatedAt time.Time
}
