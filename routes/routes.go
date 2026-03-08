// Platform Desa — Routes
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package routes

import (
	"github.com/Arlchoose-code/platform-desa-backend/controllers"
	"github.com/Arlchoose-code/platform-desa-backend/controllers/admin"
	"github.com/Arlchoose-code/platform-desa-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api/v1")

	auth := api.Group("/admin/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/refresh", controllers.RefreshToken)
	}

	protected := api.Group("/admin")
	protected.Use(middlewares.Auth())
	{
		protected.POST("/auth/logout", controllers.Logout)
		protected.GET("/auth/me", controllers.GetMe)
		protected.PUT("/auth/me", controllers.UpdateMe)
		protected.PUT("/auth/me/password", controllers.ChangePassword)

		users := protected.Group("/users")
		users.Use(middlewares.Role("superadmin"))
		{
			users.GET("", admin.GetUsers)
			users.POST("", admin.CreateUser)
			users.GET("/activity-logs", admin.GetActivityLogs)
			users.GET("/:id", admin.GetUser)
			users.PUT("/:id", admin.UpdateUser)
			users.PUT("/:id/avatar", admin.UpdateAvatar)
			users.DELETE("/:id/avatar", admin.DeleteAvatar)
			users.PUT("/:id/reset-password", admin.ResetPassword)
			users.DELETE("/:id", admin.DeleteUser)
		}
		media := protected.Group("/media")
		{
			media.GET("", admin.GetMediaList)
			media.GET("/:id", admin.GetMedia)
			media.POST("", admin.UploadMedia)
			media.POST("/batch", admin.UploadMultipleMedia)
			media.POST("/youtube", admin.AddYoutubeMedia)
			media.POST("/drive", admin.AddDriveMedia)
			media.DELETE("/batch", admin.DeleteMultipleMedia)
			media.DELETE("/:id", admin.DeleteMedia)
		}
	}
}
