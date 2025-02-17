package middleware

import (
	"context"
	"danmu-http/internal/app"
	"danmu-http/internal/model"
	"danmu-http/logger"
	"danmu-http/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	RoleAdmin = "admin"
	RoleGuest = "guest"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			logger.Error().Msg("token is empty")
			app.NewGin(c).Response(http.StatusUnauthorized, app.Unauthorized, nil)
			c.Abort()
			return
		}

		// Remove 'Bearer ' prefix if exists
		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := utils.ParseToken(token)
		if err != nil {
			logger.Error().Err(err).Msg("parse token failed")
			app.NewGin(c).Response(http.StatusUnauthorized, app.Unauthorized, nil)
			c.Abort()
			return
		}

		// Set auth info to context
		auth := &model.Auth{
			ID:    claims.ID,
			Email: claims.Email,
			Role:  claims.Role,
		}
		c.Set("auth", auth)
		ctx := context.WithValue(c.Request.Context(), "auth", auth)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, exists := c.Get("auth")
		if !exists {
			logger.Error().Msg("auth not found in context")
			app.NewGin(c).Response(http.StatusUnauthorized, app.Unauthorized, nil)
			c.Abort()
			return
		}

		if auth.(*model.Auth).Role != RoleAdmin {
			logger.Error().Str("role", auth.(*model.Auth).Role).Msg("permission denied")
			app.NewGin(c).Response(http.StatusForbidden, app.ERROR, "permission denied")
			c.Abort()
			return
		}
		c.Next()
	}
}

func GetAuthFromContext(ctx context.Context) (*model.Auth, error) {
	auth, exists := ctx.Value("auth").(*model.Auth)
	if !exists {
		return nil, fmt.Errorf("auth not found in context")
	}

	if auth == nil {
		return nil, fmt.Errorf("auth is nil")
	}

	return auth, nil
}
