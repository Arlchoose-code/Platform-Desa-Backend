// Platform Desa — Admin Wisata Controller
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

func GetKategoriWisataList(c *gin.Context) {
	showTrash := c.Query("trash") == "true"
	search := c.Query("search")

	query := database.DB.Model(&models.KategoriWisata{}).Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}

	var kategori []models.KategoriWisata
	query.Order("nama ASC").Find(&kategori)

	helpers.OK(c, "Berhasil", kategori)
}

func GetKategoriWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriWisata
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", kategori)
}

func CreateKategoriWisata(c *gin.Context) {
	var req structs.KategoriWisataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.KategoriWisata
	if err := database.DB.Where("nama = ? AND is_deleted = false", req.Nama).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kategori sudah ada", nil)
		return
	}

	kategori := models.KategoriWisata{
		Nama: req.Nama,
		Slug: helpers.UniqueSlug(helpers.GenerateSlug(req.Nama)),
		Icon: req.Icon,
	}

	if err := database.DB.Create(&kategori).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah kategori")
		return
	}

	helpers.Log(c, "create", "wisata", "Menambah kategori wisata: "+req.Nama)

	helpers.Created(c, "Kategori berhasil ditambahkan", kategori)
}

func UpdateKategoriWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriWisata
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	var req structs.KategoriWisataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.KategoriWisata
	if err := database.DB.Where("nama = ? AND id != ? AND is_deleted = false", req.Nama, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kategori sudah ada", nil)
		return
	}

	if req.Nama != kategori.Nama {
		kategori.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Nama))
	}
	kategori.Nama = req.Nama
	kategori.Icon = req.Icon

	if err := database.DB.Save(&kategori).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui kategori")
		return
	}

	helpers.Log(c, "update", "wisata", "Memperbarui kategori wisata: "+req.Nama)

	helpers.OK(c, "Kategori berhasil diperbarui", kategori)
}

func DeleteKategoriWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriWisata
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	var count int64
	database.DB.Model(&models.Wisata{}).Where("kategori_id = ? AND is_deleted = false", id).Count(&count)
	if count > 0 {
		helpers.BadRequest(c, "Kategori masih digunakan oleh wisata aktif", nil)
		return
	}

	database.DB.Model(&kategori).Update("is_deleted", true)

	helpers.Log(c, "delete", "wisata", "Menghapus kategori wisata: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dihapus", nil)
}

func ForceDeleteKategoriWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriWisata
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan di trash")
		return
	}

	database.DB.Model(&models.Wisata{}).Where("kategori_id = ?", id).Updates(map[string]interface{}{
		"kategori_id":   nil,
		"kategori_nama": nil,
	})

	database.DB.Delete(&kategori)

	helpers.Log(c, "force_delete", "wisata", "Menghapus permanen kategori wisata: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dihapus permanen", nil)
}

func RestoreKategoriWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriWisata
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan di trash")
		return
	}

	database.DB.Model(&kategori).Update("is_deleted", false)

	helpers.Log(c, "restore", "wisata", "Memulihkan kategori wisata: "+kategori.Nama)

	helpers.OK(c, "Kategori berhasil dipulihkan", nil)
}

func GetWisataList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	kategoriID := c.Query("kategori_id")
	isPublished := c.Query("is_published")
	showTrash := c.Query("trash") == "true"

	query := database.DB.Model(&models.Wisata{}).
		Preload("Kategori").
		Preload("Thumbnail").
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}
	if kategoriID != "" {
		query = query.Where("kategori_id = ?", kategoriID)
	}
	if isPublished != "" {
		query = query.Where("is_published = ?", isPublished == "true")
	}

	var total int64
	query.Count(&total)

	var wisata []models.Wisata
	query.Scopes(helpers.Paginate(pg)).Order("created_at DESC").Find(&wisata)

	helpers.OKPaginated(c, "Berhasil", wisata, pg.Meta(total))
}

func GetWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var wisata models.Wisata
	if err := database.DB.Preload("Kategori").Preload("Thumbnail").Preload("Galeri.Media").
		Where("id = ? AND is_deleted = false", id).
		First(&wisata).Error; err != nil {
		helpers.NotFound(c, "Wisata tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", wisata)
}

func CreateWisata(c *gin.Context) {
	var req structs.WisataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	wisata := models.Wisata{
		Nama:        req.Nama,
		Slug:        helpers.UniqueSlug(helpers.GenerateSlug(req.Nama)),
		Deskripsi:   req.Deskripsi,
		Alamat:      req.Alamat,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		JamBuka:     req.JamBuka,
		JamTutup:    req.JamTutup,
		HargaTiket:  req.HargaTiket,
		Kontak:      req.Kontak,
		ThumbnailID: req.ThumbnailID,
		Fasilitas:   req.Fasilitas,
		IsPublished: req.IsPublished,
	}

	if req.KategoriID != nil {
		var kategori models.KategoriWisata
		if err := database.DB.Where("id = ? AND is_deleted = false", req.KategoriID).First(&kategori).Error; err == nil {
			wisata.KategoriID = req.KategoriID
			wisata.KategoriNama = &kategori.Nama
		}
	}

	if err := database.DB.Create(&wisata).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah wisata")
		return
	}

	helpers.Log(c, "create", "wisata", "Menambah wisata: "+req.Nama)

	if req.IsPublished {
		go helpers.RevalidatePaths("/wisata", "/wisata/"+wisata.Slug)
	}

	database.DB.Preload("Kategori").Preload("Thumbnail").First(&wisata, wisata.ID)
	helpers.Created(c, "Wisata berhasil ditambahkan", wisata)
}

func UpdateWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var wisata models.Wisata
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&wisata).Error; err != nil {
		helpers.NotFound(c, "Wisata tidak ditemukan")
		return
	}

	var req structs.WisataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	if req.Nama != wisata.Nama {
		wisata.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Nama))
	}

	wisata.Nama = req.Nama
	wisata.Deskripsi = req.Deskripsi
	wisata.Alamat = req.Alamat
	wisata.Latitude = req.Latitude
	wisata.Longitude = req.Longitude
	wisata.JamBuka = req.JamBuka
	wisata.JamTutup = req.JamTutup
	wisata.HargaTiket = req.HargaTiket
	wisata.Kontak = req.Kontak
	wisata.ThumbnailID = req.ThumbnailID
	wisata.Fasilitas = req.Fasilitas
	wisata.IsPublished = req.IsPublished

	if req.KategoriID != nil {
		var kategori models.KategoriWisata
		if err := database.DB.Where("id = ? AND is_deleted = false", req.KategoriID).First(&kategori).Error; err == nil {
			wisata.KategoriID = req.KategoriID
			wisata.KategoriNama = &kategori.Nama
		}
	} else {
		wisata.KategoriID = nil
		wisata.KategoriNama = nil
	}

	if err := database.DB.Save(&wisata).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui wisata")
		return
	}

	helpers.Log(c, "update", "wisata", "Memperbarui wisata: "+req.Nama)

	go helpers.RevalidatePaths("/wisata", "/wisata/"+wisata.Slug)

	database.DB.Preload("Kategori").Preload("Thumbnail").First(&wisata, wisata.ID)
	helpers.OK(c, "Wisata berhasil diperbarui", wisata)
}

func DeleteWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var wisata models.Wisata
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&wisata).Error; err != nil {
		helpers.NotFound(c, "Wisata tidak ditemukan")
		return
	}

	database.DB.Model(&wisata).Update("is_deleted", true)

	helpers.Log(c, "delete", "wisata", "Menghapus wisata: "+wisata.Nama)

	go helpers.RevalidatePath("/wisata")

	helpers.OK(c, "Wisata berhasil dihapus", nil)
}

func ForceDeleteWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var wisata models.Wisata
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&wisata).Error; err != nil {
		helpers.NotFound(c, "Wisata tidak ditemukan di trash")
		return
	}

	var galeri []models.WisataGaleri
	database.DB.Where("wisata_id = ?", wisata.ID).Preload("Media").Find(&galeri)
	for _, g := range galeri {
		if g.Media != nil {
			go helpers.DeleteFile(g.Media.Path)
			database.DB.Delete(&g.Media)
		}
	}
	database.DB.Where("wisata_id = ?", wisata.ID).Delete(&models.WisataGaleri{})

	if wisata.ThumbnailID != nil {
		var media models.Media
		if err := database.DB.First(&media, wisata.ThumbnailID).Error; err == nil {
			go helpers.DeleteFile(media.Path)
			database.DB.Delete(&media)
		}
	}

	database.DB.Delete(&wisata)

	helpers.Log(c, "force_delete", "wisata", "Menghapus permanen wisata: "+wisata.Nama)

	helpers.OK(c, "Wisata berhasil dihapus permanen", nil)
}

func RestoreWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var wisata models.Wisata
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&wisata).Error; err != nil {
		helpers.NotFound(c, "Wisata tidak ditemukan di trash")
		return
	}

	database.DB.Model(&wisata).Update("is_deleted", false)

	helpers.Log(c, "restore", "wisata", "Memulihkan wisata: "+wisata.Nama)

	helpers.OK(c, "Wisata berhasil dipulihkan", nil)
}

func PublishWisata(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var wisata models.Wisata
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&wisata).Error; err != nil {
		helpers.NotFound(c, "Wisata tidak ditemukan")
		return
	}

	newStatus := !wisata.IsPublished
	database.DB.Model(&wisata).Update("is_published", newStatus)

	go helpers.RevalidatePaths("/wisata", "/wisata/"+wisata.Slug)

	msg := "Wisata berhasil dipublikasikan"
	if !newStatus {
		msg = "Wisata berhasil disembunyikan"
	}

	helpers.Log(c, "publish", "wisata", msg+": "+wisata.Nama)

	helpers.OK(c, msg, gin.H{"is_published": newStatus})
}

func AddWisataGaleri(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var wisata models.Wisata
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&wisata).Error; err != nil {
		helpers.NotFound(c, "Wisata tidak ditemukan")
		return
	}

	var req structs.WisataGaleriRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	galeri := models.WisataGaleri{
		WisataID: wisata.ID,
		MediaID:  req.MediaID,
		Caption:  req.Caption,
		Urutan:   req.Urutan,
	}

	if err := database.DB.Create(&galeri).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah galeri")
		return
	}

	go helpers.RevalidatePaths("/wisata", "/wisata/"+wisata.Slug)

	database.DB.Preload("Media").First(&galeri, galeri.ID)
	helpers.Created(c, "Galeri berhasil ditambahkan", galeri)
}

func DeleteWisataGaleri(c *gin.Context) {
	wisataID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	galeriID, err := strconv.Atoi(c.Param("galeri_id"))
	if err != nil {
		helpers.BadRequest(c, "ID galeri tidak valid", nil)
		return
	}

	var galeri models.WisataGaleri
	if err := database.DB.Preload("Media").Where("id = ? AND wisata_id = ?", galeriID, wisataID).First(&galeri).Error; err != nil {
		helpers.NotFound(c, "Galeri tidak ditemukan")
		return
	}

	if galeri.Media != nil {
		go helpers.DeleteFile(galeri.Media.Path)
		database.DB.Delete(&galeri.Media)
	}

	database.DB.Delete(&galeri)

	var wisata models.Wisata
	if err := database.DB.First(&wisata, wisataID).Error; err == nil {
		go helpers.RevalidatePaths("/wisata", "/wisata/"+wisata.Slug)
	}

	helpers.OK(c, "Galeri berhasil dihapus", nil)
}

func UpdateUrutanWisataGaleri(c *gin.Context) {
	wisataID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var req structs.UrutanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	for _, item := range req.Items {
		database.DB.Model(&models.WisataGaleri{}).
			Where("id = ? AND wisata_id = ?", item.ID, wisataID).
			Update("urutan", item.Urutan)
	}

	helpers.OK(c, "Urutan berhasil diperbarui", nil)
}
