package model

type Topic struct {
	BaseModel
	Name string `gorm:"type:varchar(50);uniqueIndex;not null"`
}
