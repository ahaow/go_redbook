package middleware

import (
	"go_redbook/config"
	"go_redbook/internal/pkg/jwtutil"
	"go_redbook/internal/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuth(jwtCfg config.JwtConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, 40101, "未登录")
			c.Abort()
			return
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			response.Error(c, http.StatusUnauthorized, 40102, "token格式错误")
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, prefix)

		claims, err := jwtutil.ParseToken(jwtCfg, tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 40103, "token无效或已过期")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
