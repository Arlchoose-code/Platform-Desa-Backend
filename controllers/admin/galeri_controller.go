// Platform Desa — Admin Galeri Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package admin

import (
	"strconv"

	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/Arlchoose-code/platform-desa-backend/structs"
	"github.com/gin-gonic/gin"
)

func GetKategoriGaleriList(c *gin.Context) {
	showTrash := c.Query("trash") == "true"
	search := c.Query("search")
	tipe := c.Query("tipe")

	query := database.DB.Model(&models.KategoriGaleri{}).Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}
	if tipe != "" {
		query = query.Where("tipe = ?", tipe)
	}

	var kategori []models.KategoriGaleri
	query.Order("nama ASC").Find(&kategori)

	helpers.OK(c, "Berhasil", kategori)
}

func GetKategoriGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriGaleri
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", kategori)
}

func CreateKategoriGaleri(c *gin.Context) {
	var req structs.KategoriGaleriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.KategoriGaleri
	if err := database.DB.Where("nama = ? AND is_deleted = false", req.Nama).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kategori sudah ada", nil)
		return
	}

	kategori := models.KategoriGaleri{
		Nama: req.Nama,
		Slug: helpers.UniqueSlug(helpers.GenerateSlug(req.Nama)),
		Tipe: req.Tipe,
	}

	if err := database.DB.Create(&kategori).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah kategori")
		return
	}

	helpers.Log(c, "create", "galeri", "Menambah kategori galeri: "+req.Nama)

	helpers.Created(c, "Kategori berhasil ditambahkan", kategori)
}

func UpdateKategoriGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriGaleri
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	var req structs.KategoriGaleriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.KategoriGaleri
	if err := database.DB.Where("nama = ? AND id != ? AND is_deleted = false", req.Nama, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kategori sudah ada", nil)
		return
	}

	if req.Nama != kategori.Nama {
		kategori.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Nama))
	}
	kategori.Nama = req.Nama
	kategori.Tipe = req.Tipe

	if err := database.DB.Save(&kategori).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui kategori")
		return
	}

	helpers.Log(c, "update", "galeri", "Memperbarui kategori galeri: "+req.Nama)

	helpers.OK(c, "Kategori berhasil diperbarui", kategori)
}

func DeleteKategoriGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriGaleri
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	var count int64
	database.DB.Model(&models.Galeri{}).Where("kategori_id = ? AND is_deleted = false", id).Count(&count)
	if count > 0 {
		helpers.BadRequest(c, "Kategori masih digunakan oleh galeri aktif", nil)
		return
	}

	database.DB.Model(&kategori).Update("is_deleted", true)

	helpers.Log(c, "delete", "galeri", "Menghapus kategori galeri: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dihapus", nil)
}

func ForceDeleteKategoriGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriGaleri
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan di trash")
		return
	}

	database.DB.Model(&models.Galeri{}).Where("kategori_id = ?", id).Updates(map[string]interface{}{
		"kategori_id":   nil,
		"kategori_nama": nil,
	})

	database.DB.Delete(&kategori)

	helpers.Log(c, "force_delete", "galeri", "Menghapus permanen kategori galeri: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dihapus permanen", nil)
}

func RestoreKategoriGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriGaleri
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan di trash")
		return
	}

	database.DB.Model(&kategori).Update("is_deleted", false)

	helpers.Log(c, "restore", "galeri", "Memulihkan kategori galeri: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dipulihkan", nil)
}

func GetGaleriList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	kategoriID := c.Query("kategori_id")
	isPublished := c.Query("is_published")
	showTrash := c.Query("trash") == "true"

	query := database.DB.Model(&models.Galeri{}).
		Preload("Kategori").
		Preload("Media").
		Preload("Thumbnail").
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("judul LIKE ?", "%"+search+"%")
	}
	if kategoriID != "" {
		query = query.Where("kategori_id = ?", kategoriID)
	}
	if isPublished != "" {
		query = query.Where("is_published = ?", isPublished == "true")
	}

	var total int64
	query.Count(&total)

	var galeri []models.Galeri
	query.Scopes(helpers.Paginate(pg)).Order("urutan ASC, created_at DESC").Find(&galeri)

	helpers.OKPaginated(c, "Berhasil", galeri, pg.Meta(total))
}

func GetGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var galeri models.Galeri
	if err := database.DB.Preload("Kategori").Preload("Media").Preload("Thumbnail").
		Where("id = ? AND is_deleted = false", id).
		First(&galeri).Error; err != nil {
		helpers.NotFound(c, "Galeri tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", galeri)
}

func CreateGaleri(c *gin.Context) {
	var req structs.GaleriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	galeri := models.Galeri{
		Judul:       req.Judul,
		Deskripsi:   req.Deskripsi,
		MediaID:     req.MediaID,
		ThumbnailID: req.ThumbnailID,
		Tanggal:     req.Tanggal,
		Fotografer:  req.Fotografer,
		IsPublished: req.IsPublished,
		Urutan:      req.Urutan,
	}

	if req.KategoriID != nil {
		var kategori models.KategoriGaleri
		if err := database.DB.Where("id = ? AND is_deleted = false", req.KategoriID).First(&kategori).Error; err == nil {
			galeri.KategoriID = req.KategoriID
			galeri.KategoriNama = &kategori.Nama
		}
	}

	if err := database.DB.Create(&galeri).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah galeri")
		return
	}

	helpers.Log(c, "create", "galeri", "Menambah galeri: "+req.Judul)

	if req.IsPublished {
		go helpers.RevalidatePath("/galeri")
	}

	database.DB.Preload("Kategori").Preload("Media").Preload("Thumbnail").First(&galeri, galeri.ID)
	helpers.Created(c, "Galeri berhasil ditambahkan", galeri)
}

func UpdateGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var galeri models.Galeri
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&galeri).Error; err != nil {
		helpers.NotFound(c, "Galeri tidak ditemukan")
		return
	}

	var req structs.GaleriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	galeri.Judul = req.Judul
	galeri.Deskripsi = req.Deskripsi
	galeri.MediaID = req.MediaID
	galeri.ThumbnailID = req.ThumbnailID
	galeri.Tanggal = req.Tanggal
	galeri.Fotografer = req.Fotografer
	galeri.IsPublished = req.IsPublished
	galeri.Urutan = req.Urutan

	if req.KategoriID != nil {
		var kategori models.KategoriGaleri
		if err := database.DB.Where("id = ? AND is_deleted = false", req.KategoriID).First(&kategori).Error; err == nil {
			galeri.KategoriID = req.KategoriID
			galeri.KategoriNama = &kategori.Nama
		}
	} else {
		galeri.KategoriID = nil
		galeri.KategoriNama = nil
	}

	if err := database.DB.Save(&galeri).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui galeri")
		return
	}

	helpers.Log(c, "update", "galeri", "Memperbarui galeri: "+req.Judul)

	go helpers.RevalidatePath("/galeri")

	database.DB.Preload("Kategori").Preload("Media").Preload("Thumbnail").First(&galeri, galeri.ID)
	helpers.OK(c, "Galeri berhasil diperbarui", galeri)
}

func DeleteGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var galeri models.Galeri
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&galeri).Error; err != nil {
		helpers.NotFound(c, "Galeri tidak ditemukan")
		return
	}

	database.DB.Model(&galeri).Update("is_deleted", true)

	helpers.Log(c, "delete", "galeri", "Menghapus galeri: "+galeri.Judul)

	go helpers.RevalidatePath("/galeri")

	helpers.OK(c, "Galeri berhasil dihapus", nil)
}

func ForceDeleteGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var galeri models.Galeri
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&galeri).Error; err != nil {
		helpers.NotFound(c, "Galeri tidak ditemukan di trash")
		return
	}

	var media models.Media
	if err := database.DB.First(&media, galeri.MediaID).Error; err == nil {
		go helpers.DeleteFile(media.Path)
		database.DB.Delete(&media)
	}

	if galeri.ThumbnailID != nil && *galeri.ThumbnailID != galeri.MediaID {
		var thumb models.Media
		if err := database.DB.First(&thumb, galeri.ThumbnailID).Error; err == nil {
			go helpers.DeleteFile(thumb.Path)
			database.DB.Delete(&thumb)
		}
	}

	database.DB.Delete(&galeri)

	helpers.Log(c, "force_delete", "galeri", "Menghapus permanen galeri: "+galeri.Judul)

	helpers.OK(c, "Galeri berhasil dihapus permanen", nil)
}

func RestoreGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var galeri models.Galeri
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&galeri).Error; err != nil {
		helpers.NotFound(c, "Galeri tidak ditemukan di trash")
		return
	}

	database.DB.Model(&galeri).Update("is_deleted", false)

	helpers.Log(c, "restore", "galeri", "Memulihkan galeri: "+galeri.Judul)

	helpers.OK(c, "Galeri berhasil dipulihkan", nil)
}

func PublishGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var galeri models.Galeri
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&galeri).Error; err != nil {
		helpers.NotFound(c, "Galeri tidak ditemukan")
		return
	}

	newStatus := !galeri.IsPublished
	database.DB.Model(&galeri).Update("is_published", newStatus)

	go helpers.RevalidatePath("/galeri")

	msg := "Galeri berhasil dipublikasikan"
	if !newStatus {
		msg = "Galeri berhasil disembunyikan"
	}

	helpers.Log(c, "publish", "galeri", msg+": "+galeri.Judul)

	helpers.OK(c, msg, gin.H{"is_published": newStatus})
}
