package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// response 是统一的 JSON 响应格式。
type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// ok 返回成功响应。
func ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, response{
		Code: 0,
		Msg:  "ok",
		Data: data,
	})
}

// fail 返回失败响应。
// httpStatus 是 HTTP 状态码，code 是业务错误码。
func fail(c *gin.Context, httpStatus, code int, msg string) {
	c.JSON(httpStatus, response{
		Code: code,
		Msg:  msg,
	})
}
