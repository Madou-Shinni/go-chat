package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-chat/internal/entity"
	"go-chat/internal/pkg/auth"
)

type SessionInterface interface {
	IsExistBlackList(ctx context.Context, token string) bool
}

// JwtAuth 授权中间件
func JwtAuth(secret string, guard string, session SessionInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := auth.GetJwtToken(c)

		claims, err := check(guard, secret, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		// 这里还需要验证 token 黑名单
		if session.IsExistBlackList(context.Background(), token) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "请登录再试！"})
			c.Abort()
			return
		}

		// 设置登录用户ID
		uid, _ := strconv.Atoi(claims.Id)

		c.Set(entity.LoginUserID, uid)

		// 记录 jwt 相关信息
		c.Set("jwt", map[string]string{
			"token":      token,
			"expires_at": strconv.Itoa(int(claims.ExpiresAt)),
		})

		c.Next()
	}
}

func check(guard string, secret string, token string) (*auth.JwtAuthClaims, error) {
	if token == "" {
		return nil, errors.New("请登录后操作! ")
	}

	claims, err := auth.VerifyJwtToken(token, secret)
	if err != nil {
		return nil, err
	}

	// 判断权限认证守卫是否一致
	if claims.Valid() != nil || claims.Guard != guard {
		return nil, errors.New("请登录后操作! ")
	}

	return claims, nil
}
