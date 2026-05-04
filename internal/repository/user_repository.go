package repository

import (
	"context"
	"errors"

	"go_redbook/internal/model"

	"gorm.io/gorm"
)

// UserRepository 定义用户表的数据访问能力。
// 上层 service 只依赖这个接口，不直接操作 GORM，后面做单元测试会更方便。
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uint) (*model.User, error)
	FindByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error)
	FindByAccount(ctx context.Context, account string) (*model.User, error)
	List(ctx context.Context, offset, limit int) ([]model.User, int64, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例。
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 只负责把用户写入数据库，不处理“能不能创建”的业务判断。
func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据主键查询用户。
// 查不到时返回 nil, nil，让 service 层决定返回什么业务错误。
func (r *userRepository) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsernameOrEmail 用于创建用户前检查唯一性。
func (r *userRepository) FindByUsernameOrEmail(ctx context.Context, username, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("username = ? OR email = ?", username, email).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByAccount 根据账号查询用户。
// 这里允许用户用 username 或 email 登录。
func (r *userRepository) FindByAccount(ctx context.Context, account string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("username = ? OR email = ?", account, account).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List 返回分页数据和总数。
// repository 接收 offset/limit，不关心 page/page_size 这种 HTTP 查询参数。
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	db := r.db.WithContext(ctx).Model(&model.User{})
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := db.Order("id DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}
