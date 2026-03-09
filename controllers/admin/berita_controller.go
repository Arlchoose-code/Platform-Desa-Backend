// Platform Desa — Admin Berita Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package admin

import (
	"strconv"
	"time"

	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/Arlchoose-code/platform-desa-backend/structs"
	"github.com/gin-gonic/gin"
)

func GetKategoriBeritaList(c *gin.Context) {
	showTrash := c.Query("trash") == "true"
	search := c.Query("search")

	query := database.DB.Model(&models.KategoriBerita{}).
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}

	var kategori []models.KategoriBerita
	query.Order("nama ASC").Find(&kategori)

	helpers.OK(c, "Berhasil", kategori)
}

func GetKategoriBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriBerita
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", kategori)
}

func CreateKategoriBerita(c *gin.Context) {
	var req structs.KategoriBeritaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.KategoriBerita
	if err := database.DB.Where("nama = ? AND is_deleted = false", req.Nama).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kategori sudah ada", nil)
		return
	}

	kategori := models.KategoriBerita{
		Nama:  req.Nama,
		Slug:  helpers.UniqueSlug(helpers.GenerateSlug(req.Nama)),
		Warna: req.Warna,
	}

	if err := database.DB.Create(&kategori).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah kategori")
		return
	}

	helpers.Log(c, "create", "berita", "Menambah kategori berita: "+req.Nama)

	helpers.Created(c, "Kategori berhasil ditambahkan", kategori)
}

func UpdateKategoriBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriBerita
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	var req structs.KategoriBeritaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.KategoriBerita
	if err := database.DB.Where("nama = ? AND id != ? AND is_deleted = false", req.Nama, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kategori sudah ada", nil)
		return
	}

	if req.Nama != kategori.Nama {
		kategori.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Nama))
	}
	kategori.Nama = req.Nama
	kategori.Warna = req.Warna

	if err := database.DB.Save(&kategori).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui kategori")
		return
	}

	helpers.Log(c, "update", "berita", "Memperbarui kategori berita: "+req.Nama)

	helpers.OK(c, "Kategori berhasil diperbarui", kategori)
}

func DeleteKategoriBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriBerita
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	var count int64
	database.DB.Model(&models.Berita{}).Where("kategori_id = ? AND is_deleted = false", id).Count(&count)
	if count > 0 {
		helpers.BadRequest(c, "Kategori masih digunakan oleh berita aktif", nil)
		return
	}

	database.DB.Model(&kategori).Update("is_deleted", true)

	helpers.Log(c, "delete", "berita", "Menghapus kategori berita: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dihapus", nil)
}

func ForceDeleteKategoriBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriBerita
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan di trash")
		return
	}

	database.DB.Model(&models.Berita{}).Where("kategori_id = ?", id).Updates(map[string]interface{}{
		"kategori_id":   nil,
		"kategori_nama": nil,
	})

	database.DB.Delete(&kategori)

	helpers.Log(c, "force_delete", "berita", "Menghapus permanen kategori berita: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dihapus permanen", nil)
}

func RestoreKategoriBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriBerita
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan di trash")
		return
	}

	database.DB.Model(&kategori).Update("is_deleted", false)

	helpers.Log(c, "restore", "berita", "Memulihkan kategori berita: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dipulihkan", nil)
}

func GetBeritaList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	kategoriID := c.Query("kategori_id")
	status := c.Query("status")
	showTrash := c.Query("trash") == "true"

	query := database.DB.Model(&models.Berita{}).
		Preload("Kategori").
		Preload("Thumbnail").
		Preload("Penulis").
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("judul LIKE ? OR ringkasan LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if kategoriID != "" {
		query = query.Where("kategori_id = ?", kategoriID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var berita []models.Berita
	query.Scopes(helpers.Paginate(pg)).Order("created_at DESC").Find(&berita)

	helpers.OKPaginated(c, "Berhasil", berita, pg.Meta(total))
}

func GetBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var berita models.Berita
	if err := database.DB.Preload("Kategori").Preload("Thumbnail").Preload("Penulis").
		Where("id = ? AND is_deleted = false", id).
		First(&berita).Error; err != nil {
		helpers.NotFound(c, "Berita tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", berita)
}

func CreateBerita(c *gin.Context) {
	var req structs.BeritaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	userID := c.GetUint("user_id")

	berita := models.Berita{
		Judul:       req.Judul,
		Slug:        helpers.UniqueSlug(helpers.GenerateSlug(req.Judul)),
		Ringkasan:   req.Ringkasan,
		Isi:         req.Isi,
		ThumbnailID: req.ThumbnailID,
		PenulisID:   &userID,
		Status:      req.Status,
	}

	if req.KategoriID != nil {
		var kategori models.KategoriBerita
		if err := database.DB.Where("id = ? AND is_deleted = false", req.KategoriID).First(&kategori).Error; err == nil {
			berita.KategoriID = req.KategoriID
			berita.KategoriNama = &kategori.Nama
		}
	}

	if req.Status == "published" {
		now := time.Now()
		berita.PublishedAt = &now
	}

	if err := database.DB.Create(&berita).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah berita")
		return
	}

	helpers.Log(c, "create", "berita", "Menambah berita: "+req.Judul)

	if req.Status == "published" {
		go helpers.RevalidatePaths("/berita", "/berita/"+berita.Slug)
	}

	database.DB.Preload("Kategori").Preload("Thumbnail").Preload("Penulis").First(&berita, berita.ID)
	helpers.Created(c, "Berita berhasil ditambahkan", berita)
}

func UpdateBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var berita models.Berita
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&berita).Error; err != nil {
		helpers.NotFound(c, "Berita tidak ditemukan")
		return
	}

	var req structs.BeritaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	if req.Judul != berita.Judul {
		berita.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Judul))
	}

	berita.Judul = req.Judul
	berita.Ringkasan = req.Ringkasan
	berita.Isi = req.Isi
	berita.ThumbnailID = req.ThumbnailID
	berita.Status = req.Status

	if req.KategoriID != nil {
		var kategori models.KategoriBerita
		if err := database.DB.Where("id = ? AND is_deleted = false", req.KategoriID).First(&kategori).Error; err == nil {
			berita.KategoriID = req.KategoriID
			berita.KategoriNama = &kategori.Nama
		}
	} else {
		berita.KategoriID = nil
		berita.KategoriNama = nil
	}

	if req.Status == "published" && berita.PublishedAt == nil {
		now := time.Now()
		berita.PublishedAt = &now
	}

	if err := database.DB.Save(&berita).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui berita")
		return
	}

	helpers.Log(c, "update", "berita", "Memperbarui berita: "+req.Judul)

	go helpers.RevalidatePaths("/berita", "/berita/"+berita.Slug)

	database.DB.Preload("Kategori").Preload("Thumbnail").Preload("Penulis").First(&berita, berita.ID)
	helpers.OK(c, "Berita berhasil diperbarui", berita)
}

func DeleteBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var berita models.Berita
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&berita).Error; err != nil {
		helpers.NotFound(c, "Berita tidak ditemukan")
		return
	}

	database.DB.Model(&berita).Update("is_deleted", true)

	helpers.Log(c, "delete", "berita", "Menghapus berita: "+berita.Judul)

	go helpers.RevalidatePath("/berita")

	helpers.OK(c, "Berita berhasil dihapus", nil)
}

func ForceDeleteBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var berita models.Berita
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&berita).Error; err != nil {
		helpers.NotFound(c, "Berita tidak ditemukan di trash")
		return
	}

	if berita.ThumbnailID != nil {
		var media models.Media
		if err := database.DB.First(&media, berita.ThumbnailID).Error; err == nil {
			go helpers.DeleteFile(media.Path)
			database.DB.Delete(&media)
		}
	}

	database.DB.Delete(&berita)

	helpers.Log(c, "force_delete", "berita", "Menghapus permanen berita: "+berita.Judul)

	helpers.OK(c, "Berita berhasil dihapus permanen", nil)
}

func RestoreBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var berita models.Berita
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&berita).Error; err != nil {
		helpers.NotFound(c, "Berita tidak ditemukan di trash")
		return
	}

	database.DB.Model(&berita).Update("is_deleted", false)

	helpers.Log(c, "restore", "berita", "Memulihkan berita: "+berita.Judul)

	helpers.OK(c, "Berita berhasil dipulihkan", nil)
}

func PublishBerita(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var berita models.Berita
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&berita).Error; err != nil {
		helpers.NotFound(c, "Berita tidak ditemukan")
		return
	}

	if berita.Status == "published" {
		database.DB.Model(&berita).Update("status", "draft")
		helpers.Log(c, "unpublish", "berita", "Menyembunyikan berita: "+berita.Judul)
		go helpers.RevalidatePaths("/berita", "/berita/"+berita.Slug)
		helpers.OK(c, "Berita berhasil di-unpublish", gin.H{"status": "draft"})
		return
	}

	now := time.Now()
	database.DB.Model(&berita).Updates(map[string]interface{}{
		"status":       "published",
		"published_at": now,
	})

	helpers.Log(c, "publish", "berita", "Mempublikasikan berita: "+berita.Judul)

	go helpers.RevalidatePaths("/berita", "/berita/"+berita.Slug)
	helpers.OK(c, "Berita berhasil dipublikasikan", gin.H{"status": "published"})
}
