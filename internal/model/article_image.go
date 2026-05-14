package model

type ArticleImage struct {
	BaseModel
	ArticleID uint `gorm:"index;not null"`

	URL       string `gorm:"type:varchar(500);not null"`
	SortOrder int    `gorm:"default:0"` // 多图排序
}
