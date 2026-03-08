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

		protected.GET("/profil", admin.GetProfilDesa)
		protected.PUT("/profil", admin.UpdateProfilDesa)

		potensi := protected.Group("/potensi")
		{
			potensi.GET("", admin.GetPotensiList)
			potensi.POST("", admin.CreatePotensi)
			potensi.GET("/trash", admin.GetPotensiList)
			potensi.GET("/:id", admin.GetPotensi)
			potensi.PUT("/:id", admin.UpdatePotensi)
			potensi.DELETE("/:id", admin.DeletePotensi)
			potensi.DELETE("/:id/force", admin.ForceDeletePotensi)
			potensi.PUT("/:id/restore", admin.RestorePotensi)
			potensi.PUT("/:id/publish", admin.PublishPotensi)
		}

		jabatan := protected.Group("/jabatan")
		{
			jabatan.GET("", admin.GetJabatanList)
			jabatan.POST("", admin.CreateJabatan)
			jabatan.GET("/trash", admin.GetJabatanList)
			jabatan.GET("/:id", admin.GetJabatan)
			jabatan.PUT("/:id", admin.UpdateJabatan)
			jabatan.DELETE("/:id", admin.DeleteJabatan)
			jabatan.DELETE("/:id/force", admin.ForceDeleteJabatan)
			jabatan.PUT("/:id/restore", admin.RestoreJabatan)
		}

		pejabat := protected.Group("/pejabat")
		{
			pejabat.GET("", admin.GetPejabatList)
			pejabat.POST("", admin.CreatePejabat)
			pejabat.GET("/trash", admin.GetPejabatList)
			pejabat.GET("/:id", admin.GetPejabat)
			pejabat.PUT("/:id", admin.UpdatePejabat)
			pejabat.DELETE("/:id", admin.DeletePejabat)
			pejabat.DELETE("/:id/force", admin.ForceDeletePejabat)
			pejabat.PUT("/:id/restore", admin.RestorePejabat)
		}

		lembaga := protected.Group("/lembaga")
		{
			lembaga.GET("", admin.GetLembagaList)
			lembaga.POST("", admin.CreateLembaga)
			lembaga.GET("/trash", admin.GetLembagaList)
			lembaga.GET("/:id", admin.GetLembaga)
			lembaga.PUT("/:id", admin.UpdateLembaga)
			lembaga.DELETE("/:id", admin.DeleteLembaga)
			lembaga.DELETE("/:id/force", admin.ForceDeleteLembaga)
			lembaga.PUT("/:id/restore", admin.RestoreLembaga)
		}

		kategoriBerita := protected.Group("/berita/kategori")
		{
			kategoriBerita.GET("", admin.GetKategoriBeritaList)
			kategoriBerita.POST("", admin.CreateKategoriBerita)
			kategoriBerita.GET("/trash", admin.GetKategoriBeritaList)
			kategoriBerita.GET("/:id", admin.GetKategoriBerita)
			kategoriBerita.PUT("/:id", admin.UpdateKategoriBerita)
			kategoriBerita.DELETE("/:id", admin.DeleteKategoriBerita)
			kategoriBerita.DELETE("/:id/force", admin.ForceDeleteKategoriBerita)
			kategoriBerita.PUT("/:id/restore", admin.RestoreKategoriBerita)
		}

		berita := protected.Group("/berita")
		{
			berita.GET("", admin.GetBeritaList)
			berita.POST("", admin.CreateBerita)
			berita.GET("/trash", admin.GetBeritaList)
			berita.GET("/:id", admin.GetBerita)
			berita.PUT("/:id", admin.UpdateBerita)
			berita.DELETE("/:id", admin.DeleteBerita)
			berita.DELETE("/:id/force", admin.ForceDeleteBerita)
			berita.PUT("/:id/restore", admin.RestoreBerita)
			berita.PUT("/:id/publish", admin.PublishBerita)
		}

		pengumuman := protected.Group("/pengumuman")
		{
			pengumuman.GET("", admin.GetPengumumanList)
			pengumuman.POST("", admin.CreatePengumuman)
			pengumuman.GET("/trash", admin.GetPengumumanList)
			pengumuman.GET("/:id", admin.GetPengumuman)
			pengumuman.PUT("/:id", admin.UpdatePengumuman)
			pengumuman.DELETE("/:id", admin.DeletePengumuman)
			pengumuman.DELETE("/:id/force", admin.ForceDeletePengumuman)
			pengumuman.PUT("/:id/restore", admin.RestorePengumuman)
			pengumuman.PUT("/:id/publish", admin.PublishPengumuman)
		}

		agenda := protected.Group("/agenda")
		{
			agenda.GET("", admin.GetAgendaList)
			agenda.POST("", admin.CreateAgenda)
			agenda.GET("/trash", admin.GetAgendaList)
			agenda.GET("/:id", admin.GetAgenda)
			agenda.PUT("/:id", admin.UpdateAgenda)
			agenda.DELETE("/:id", admin.DeleteAgenda)
			agenda.DELETE("/:id/force", admin.ForceDeleteAgenda)
			agenda.PUT("/:id/restore", admin.RestoreAgenda)
			agenda.PUT("/:id/publish", admin.PublishAgenda)
		}
	}
}
