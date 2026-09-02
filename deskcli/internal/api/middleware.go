package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pikachuim/deskcli/internal/config"
)

type jwtClaims struct {
	UserID   int64  `json:"uid"`
	Username string `json:"sub"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func jwtSecret(cfg *config.Config) []byte {
	if cfg.JWTSecret != "" {
		return []byte(cfg.JWTSecret)
	}
	return []byte("deskcli-default-secret-change-me")
}

func signJWT(cfg *config.Config, userID int64, username, role string) (string, error) {
	claims := jwtClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(jwtSecret(cfg))
}

func parseJWT(cfg *config.Config, tokenStr string) (*jwtClaims, error) {
	t, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret(cfg), nil
	})
	if err != nil || !t.Valid {
		return nil, err
	}
	claims, ok := t.Claims.(*jwtClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}

// authMiddleware supports both JWT (Bearer) and legacy X-Token for backward compatibility.
func jwtAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. JWT Bearer token
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token := authHeader[7:]
			claims, err := parseJWT(cfg, token)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				c.Next()
				return
			}
		}

		// 2. Legacy X-Token (backward compat)
		provided := c.GetHeader("X-Token")
		if cfg.Token != "" && provided == cfg.Token {
			c.Set("user_id", int64(0))
			c.Set("username", "admin")
			c.Set("role", "admin")
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
	}
}

// requireAdmin aborts if the authenticated user is not an admin.
func requireAdmin(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": "admin required"})
		return
	}
	c.Next()
}

// wsAuthMiddleware reads token from query param for WebSocket upgrades.
func wsAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token != "" {
			claims, err := parseJWT(cfg, token)
			if err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				c.Next()
				return
			}
		}
		// Legacy token in query
		if cfg.Token != "" && c.Query("x_token") == cfg.Token {
			c.Set("user_id", int64(0))
			c.Set("username", "admin")
			c.Set("role", "admin")
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthorized"})
	}
}
