package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseUintParam 解析路径参数里的无符号整数。
// 放在 handler 包里，所有 HTTP handler 都可以复用。
func parseUintParam(c *gin.Context, key string) (uint64, error) {
	return strconv.ParseUint(c.Param(key), 10, 64)
}

// queryInt 解析查询参数，解析失败时使用默认值。
// 常用于分页、状态筛选这类 query 参数。
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
