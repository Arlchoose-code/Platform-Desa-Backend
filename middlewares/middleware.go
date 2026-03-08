// Platform Desa — Middlewares
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package middlewares

import (
	"net/http"
	"strings"

	"github.com/Arlchoose-code/platform-desa-backend/config"
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.Unauthorized(c, "Token tidak ditemukan")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			helpers.Unauthorized(c, "Format token tidak valid")
			c.Abort()
			return
		}

		claims, err := helpers.ParseToken(parts[1])
		if err != nil {
			helpers.Unauthorized(c, "Token tidak valid atau sudah kadaluarsa")
			c.Abort()
			return
		}

		if claims.Type != "access" {
			helpers.Unauthorized(c, "Token tidak valid")
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.Select("id, role, is_active").First(&user, claims.UserID).Error; err != nil {
			helpers.Unauthorized(c, "User tidak ditemukan")
			c.Abort()
			return
		}

		if !user.IsActive {
			helpers.Unauthorized(c, "Akun tidak aktif")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", user.Role)
		c.Next()
	}
}

func Role(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("user_role")
		role := userRole.(string)

		if role == "superadmin" {
			c.Next()
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		helpers.Forbidden(c, "Anda tidak memiliki akses ke halaman ini")
		c.Abort()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		for _, allowed := range strings.Split(config.App.CORS.AllowedOrigins, ",") {
			if strings.TrimSpace(allowed) == origin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
