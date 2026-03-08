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

		kependudukan := protected.Group("/kependudukan")
		{
			kependudukan.GET("/penduduk", admin.GetStatistikPendudukList)
			kependudukan.POST("/penduduk", admin.CreateStatistikPenduduk)
			kependudukan.GET("/penduduk/:id", admin.GetStatistikPenduduk)
			kependudukan.PUT("/penduduk/:id", admin.UpdateStatistikPenduduk)
			kependudukan.DELETE("/penduduk/:id", admin.DeleteStatistikPenduduk)

			kependudukan.GET("/pendidikan", admin.GetStatistikPendidikanList)
			kependudukan.POST("/pendidikan", admin.CreateStatistikPendidikan)
			kependudukan.GET("/pendidikan/:id", admin.GetStatistikPendidikan)
			kependudukan.PUT("/pendidikan/:id", admin.UpdateStatistikPendidikan)
			kependudukan.DELETE("/pendidikan/:id", admin.DeleteStatistikPendidikan)

			kependudukan.GET("/pekerjaan", admin.GetStatistikPekerjaanList)
			kependudukan.POST("/pekerjaan", admin.CreateStatistikPekerjaan)
			kependudukan.GET("/pekerjaan/:id", admin.GetStatistikPekerjaan)
			kependudukan.PUT("/pekerjaan/:id", admin.UpdateStatistikPekerjaan)
			kependudukan.DELETE("/pekerjaan/:id", admin.DeleteStatistikPekerjaan)
		}

		kategoriWisata := protected.Group("/wisata/kategori")
		{
			kategoriWisata.GET("", admin.GetKategoriWisataList)
			kategoriWisata.POST("", admin.CreateKategoriWisata)
			kategoriWisata.GET("/trash", admin.GetKategoriWisataList)
			kategoriWisata.GET("/:id", admin.GetKategoriWisata)
			kategoriWisata.PUT("/:id", admin.UpdateKategoriWisata)
			kategoriWisata.DELETE("/:id", admin.DeleteKategoriWisata)
			kategoriWisata.DELETE("/:id/force", admin.ForceDeleteKategoriWisata)
			kategoriWisata.PUT("/:id/restore", admin.RestoreKategoriWisata)
		}

		wisata := protected.Group("/wisata")
		{
			wisata.GET("", admin.GetWisataList)
			wisata.POST("", admin.CreateWisata)
			wisata.GET("/trash", admin.GetWisataList)
			wisata.GET("/:id", admin.GetWisata)
			wisata.PUT("/:id", admin.UpdateWisata)
			wisata.DELETE("/:id", admin.DeleteWisata)
			wisata.DELETE("/:id/force", admin.ForceDeleteWisata)
			wisata.PUT("/:id/restore", admin.RestoreWisata)
			wisata.PUT("/:id/publish", admin.PublishWisata)
			wisata.POST("/:id/galeri", admin.AddWisataGaleri)
			wisata.DELETE("/:id/galeri/:galeri_id", admin.DeleteWisataGaleri)
			wisata.PUT("/:id/galeri/urutan", admin.UpdateUrutanWisataGaleri)
		}

		kategoriGaleri := protected.Group("/galeri/kategori")
		{
			kategoriGaleri.GET("", admin.GetKategoriGaleriList)
			kategoriGaleri.POST("", admin.CreateKategoriGaleri)
			kategoriGaleri.GET("/trash", admin.GetKategoriGaleriList)
			kategoriGaleri.GET("/:id", admin.GetKategoriGaleri)
			kategoriGaleri.PUT("/:id", admin.UpdateKategoriGaleri)
			kategoriGaleri.DELETE("/:id", admin.DeleteKategoriGaleri)
			kategoriGaleri.DELETE("/:id/force", admin.ForceDeleteKategoriGaleri)
			kategoriGaleri.PUT("/:id/restore", admin.RestoreKategoriGaleri)
		}

		galeri := protected.Group("/galeri")
		{
			galeri.GET("", admin.GetGaleriList)
			galeri.POST("", admin.CreateGaleri)
			galeri.GET("/trash", admin.GetGaleriList)
			galeri.GET("/:id", admin.GetGaleri)
			galeri.PUT("/:id", admin.UpdateGaleri)
			galeri.DELETE("/:id", admin.DeleteGaleri)
			galeri.DELETE("/:id/force", admin.ForceDeleteGaleri)
			galeri.PUT("/:id/restore", admin.RestoreGaleri)
			galeri.PUT("/:id/publish", admin.PublishGaleri)
		}

		kategoriUMKM := protected.Group("/umkm/kategori")
		{
			kategoriUMKM.GET("", admin.GetKategoriUMKMList)
			kategoriUMKM.POST("", admin.CreateKategoriUMKM)
			kategoriUMKM.GET("/trash", admin.GetKategoriUMKMList)
			kategoriUMKM.GET("/:id", admin.GetKategoriUMKM)
			kategoriUMKM.PUT("/:id", admin.UpdateKategoriUMKM)
			kategoriUMKM.DELETE("/:id", admin.DeleteKategoriUMKM)
			kategoriUMKM.DELETE("/:id/force", admin.ForceDeleteKategoriUMKM)
			kategoriUMKM.PUT("/:id/restore", admin.RestoreKategoriUMKM)
		}

		umkm := protected.Group("/umkm")
		{
			umkm.GET("", admin.GetUMKMList)
			umkm.POST("", admin.CreateUMKM)
			umkm.GET("/trash", admin.GetUMKMList)
			umkm.GET("/:id", admin.GetUMKM)
			umkm.PUT("/:id", admin.UpdateUMKM)
			umkm.DELETE("/:id", admin.DeleteUMKM)
			umkm.DELETE("/:id/force", admin.ForceDeleteUMKM)
			umkm.PUT("/:id/restore", admin.RestoreUMKM)
			umkm.PUT("/:id/publish", admin.PublishUMKM)
			umkm.GET("/:id/produk", admin.GetProdukUMKMList)
			umkm.POST("/:id/produk", admin.CreateProdukUMKM)
			umkm.PUT("/:id/produk/:produk_id", admin.UpdateProdukUMKM)
			umkm.DELETE("/:id/produk/:produk_id", admin.DeleteProdukUMKM)
		}
	}
}
