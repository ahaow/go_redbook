package handler

import (
	"errors"
	"net/http"

	"go_redbook/internal/service"

	"github.com/gin-gonic/gin"
)

// UserHandler 只处理 HTTP 层逻辑：
// 参数绑定、参数校验、调用 service、组装响应。
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户 handler，注入 service。
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// RegisterRoutes 注册用户相关路由。
// 这里约定当前 group 是 /api/v1，所以用户路由最终是 /api/v1/users。
func (h *UserHandler) RegisterRoutes(group *gin.RouterGroup) {
	auth := group.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	users := group.Group("/users")
	users.GET("", h.List)
	users.GET("/:id", h.GetByID)
}

// createUserRequest 是创建用户接口的 HTTP 请求体。
// binding 标签由 Gin 自动校验，校验失败会在 handler 里返回 400。
type createUserRequest struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=128"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Nickname string `json:"nickname" binding:"max=64"`
}

type loginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register 处理 POST /api/v1/auth/register。
func (h *UserHandler) Register(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, 40001, "参数错误: "+err.Error())
		return
	}

	result, err := h.userService.Register(c.Request.Context(), service.CreateUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	})
	if errors.Is(err, service.ErrUserExists) {
		fail(c, http.StatusConflict, 40002, err.Error())
		return
	}
	if err != nil {
		fail(c, http.StatusInternalServerError, 50001, "注册失败")
		return
	}

	c.JSON(http.StatusCreated, response{
		Code: 0,
		Msg:  "ok",
		Data: result,
	})
}

// Login 处理 POST /api/v1/auth/login。
func (h *UserHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, 40005, "参数错误: "+err.Error())
		return
	}

	result, err := h.userService.Login(c.Request.Context(), service.LoginInput{
		Account:  req.Account,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidLogin) {
		fail(c, http.StatusUnauthorized, 40006, err.Error())
		return
	}
	if err != nil {
		fail(c, http.StatusInternalServerError, 50004, "登录失败")
		return
	}

	ok(c, result)
}

// GetByID 处理 GET /api/v1/users/:id。
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		fail(c, http.StatusBadRequest, 40003, "用户ID不正确")
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), uint(id))
	if errors.Is(err, service.ErrUserNotFound) {
		fail(c, http.StatusNotFound, 40004, err.Error())
		return
	}
	if err != nil {
		fail(c, http.StatusInternalServerError, 50002, "查询用户失败")
		return
	}

	ok(c, user)
}

// List 处理 GET /api/v1/users?page=1&page_size=10。
func (h *UserHandler) List(c *gin.Context) {
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "page_size", 10)

	result, err := h.userService.List(c.Request.Context(), page, pageSize)
	if err != nil {
		fail(c, http.StatusInternalServerError, 50003, "查询用户列表失败")
		return
	}

	ok(c, gin.H{
		"items":     result.Items,
		"total":     result.Total,
		"page":      page,
		"page_size": pageSize,
	})
}
