package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 放公共数据库字段，业务模型直接内嵌即可。
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
