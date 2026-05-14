package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 是统一的 JSON 响应格式。
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// Success 返回成功响应。
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

// Created 返回创建成功响应。
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

// Error 返回失败响应。
func Error(c *gin.Context, httpStatus, code int, msg string) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
	})
}
