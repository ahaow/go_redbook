package model

import "time"

type Topic struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(50);uniqueIndex;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
