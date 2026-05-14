package model

// User 是数据库里的用户表模型。
// model 层只描述数据结构和 GORM 映射关系，不写业务逻辑。
type User struct {
	BaseModel
	Username string `gorm:"size:64;not null;uniqueIndex" json:"username"`
	Email    string `gorm:"size:128;not null;uniqueIndex" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	Nickname string `gorm:"size:64" json:"nickname"`
	Status   int    `gorm:"not null;default:1" json:"status"`
}
