package service

import (
	"context"
	"errors"

	"go_redbook/config"
	"go_redbook/internal/model"
	"go_redbook/internal/pkg/jwtutil"
	"go_redbook/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists   = errors.New("用户已存在")
	ErrUserNotFound = errors.New("用户不存在")
	ErrInvalidLogin = errors.New("账号或密码错误")
)

// CreateUserInput 是 service 层的创建用户入参。
// handler 的 request 不直接传进来，避免 HTTP 结构污染业务层。
type CreateUserInput struct {
	Username string
	Email    string
	Password string
	Nickname string
}

// LoginInput 是登录入参。
// Account 可以是用户名，也可以是邮箱。
type LoginInput struct {
	Account  string
	Password string
}

// AuthResult 是注册和登录成功后的返回值。
// Token 给前端保存，User 给前端展示当前用户信息。
type AuthResult struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

// UserList 是用户分页列表的业务返回值。
type UserList struct {
	Items []model.User
	Total int64
}

// UserService 定义用户业务能力。
// 业务规则都放在 service：比如查重、默认状态、分页限制。
type UserService interface {
	Register(ctx context.Context, input CreateUserInput) (*AuthResult, error)
	Login(ctx context.Context, input LoginInput) (*AuthResult, error)
	GetByID(ctx context.Context, id uint) (*model.User, error)
	List(ctx context.Context, page, pageSize int) (*UserList, error)
}

type userService struct {
	userRepo repository.UserRepository
	jwtCfg   config.JwtConfig
}

// NewUserService 创建用户服务实例，注入 repository。
func NewUserService(userRepo repository.UserRepository, jwtCfg config.JwtConfig) UserService {
	return &userService{
		userRepo: userRepo,
		jwtCfg:   jwtCfg,
	}
}

// Create 创建用户。
// 这里先检查用户名或邮箱是否已存在，再设置默认状态并落库。
func (s *userService) create(ctx context.Context, input CreateUserInput) (*model.User, error) {
	existing, err := s.userRepo.FindByUsernameOrEmail(ctx, input.Username, input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUserExists
	}

	passwordHash, err := hashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: passwordHash,
		Nickname: input.Nickname,
		Status:   1,
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// Register 注册用户并返回 JWT。
func (s *userService) Register(ctx context.Context, input CreateUserInput) (*AuthResult, error) {
	user, err := s.create(ctx, input)
	if err != nil {
		return nil, err
	}
	return s.buildAuthResult(user)
}

// Login 校验账号密码并返回 JWT。
func (s *userService) Login(ctx context.Context, input LoginInput) (*AuthResult, error) {
	user, err := s.userRepo.FindByAccount(ctx, input.Account)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidLogin
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, ErrInvalidLogin
	}

	return s.buildAuthResult(user)
}

// GetByID 获取用户详情，并把“查不到”转换成业务错误。
func (s *userService) GetByID(ctx context.Context, id uint) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// List 查询用户分页列表。
// page/pageSize 的兜底和上限在 service 里处理，handler 保持薄一点。
func (s *userService) List(ctx context.Context, page, pageSize int) (*UserList, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	users, total, err := s.userRepo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, err
	}
	return &UserList{Items: users, Total: total}, nil
}

func (s *userService) buildAuthResult(user *model.User) (*AuthResult, error) {
	token, err := jwtutil.GenerateToken(s.jwtCfg, user.ID, user.Username)
	if err != nil {
		return nil, err
	}
	return &AuthResult{
		Token: token,
		User:  *user,
	}, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
