package handler

import (
	"go_redbook/internal/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CurrentUser struct {
	UserID   uint
	Username string
}

func currentUser(c *gin.Context) (*CurrentUser, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 40104, "未登录")
		return nil, false
	}

	usernameValue, exists := c.Get("username")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 40104, "未登录")
		return nil, false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40105, "用户信息错误")
		return nil, false
	}

	username, ok := usernameValue.(string)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40105, "用户信息错误")
		return nil, false
	}

	return &CurrentUser{
		UserID:   userID,
		Username: username,
	}, true
}

func parseUintParam(c *gin.Context, key string) (uint64, error) {
	return strconv.ParseUint(c.Param(key), 10, 64)
}

func queryInt(c *gin.Context, key string, defaultValue int) int {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return n
}
