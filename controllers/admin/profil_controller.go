// Platform Desa — Admin Profil Controller
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

func GetProfilDesa(c *gin.Context) {
	var profil models.ProfilDesa
	if err := database.DB.Preload("Logo").Preload("FotoDesa").First(&profil).Error; err != nil {
		helpers.OK(c, "Berhasil", nil)
		return
	}
	helpers.OK(c, "Berhasil", profil)
}

func UpdateProfilDesa(c *gin.Context) {
	var req structs.ProfilDesaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var profil models.ProfilDesa
	database.DB.First(&profil)

	profil.NamaDesa = req.NamaDesa
	profil.NamaKecamatan = req.NamaKecamatan
	profil.NamaKabupaten = req.NamaKabupaten
	profil.NamaProvinsi = req.NamaProvinsi
	profil.KodePos = req.KodePos
	profil.Alamat = req.Alamat
	profil.Telepon = req.Telepon
	profil.Email = req.Email
	profil.Website = req.Website
	profil.Facebook = req.Facebook
	profil.Instagram = req.Instagram
	profil.Twitter = req.Twitter
	profil.Youtube = req.Youtube
	profil.Tiktok = req.Tiktok
	profil.LuasWilayah = req.LuasWilayah
	profil.JumlahDusun = req.JumlahDusun
	profil.JumlahRW = req.JumlahRW
	profil.JumlahRT = req.JumlahRT
	profil.BatasUtara = req.BatasUtara
	profil.BatasSelatan = req.BatasSelatan
	profil.BatasTimur = req.BatasTimur
	profil.BatasBarat = req.BatasBarat
	profil.Latitude = req.Latitude
	profil.Longitude = req.Longitude
	profil.LogoID = req.LogoID
	profil.FotoDesaID = req.FotoDesaID
	profil.Sejarah = req.Sejarah
	profil.Visi = req.Visi
	profil.Misi = req.Misi

	if err := database.DB.Save(&profil).Error; err != nil {
		helpers.InternalError(c, "Gagal menyimpan profil desa")
		return
	}

	helpers.Log(c, "update", "profil", "Memperbarui profil desa: "+req.NamaDesa)

	go helpers.RevalidatePaths("/", "/profil")

	database.DB.Preload("Logo").Preload("FotoDesa").First(&profil)
	helpers.OK(c, "Profil desa berhasil diperbarui", profil)
}

func GetPotensiList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	kategori := c.Query("kategori")
	search := c.Query("search")
	showTrash := c.Query("trash") == "true"

	query := database.DB.Model(&models.PotensiDesa{}).
		Preload("Foto").
		Where("is_deleted = ?", showTrash)

	if kategori != "" {
		query = query.Where("kategori = ?", kategori)
	}
	if search != "" {
		query = query.Where("judul LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var potensi []models.PotensiDesa
	query.Scopes(helpers.Paginate(pg)).Order("urutan ASC, created_at DESC").Find(&potensi)

	helpers.OKPaginated(c, "Berhasil", potensi, pg.Meta(total))
}

func GetPotensi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var potensi models.PotensiDesa
	if err := database.DB.Preload("Foto").
		Where("id = ? AND is_deleted = false", id).
		First(&potensi).Error; err != nil {
		helpers.NotFound(c, "Potensi desa tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", potensi)
}

func CreatePotensi(c *gin.Context) {
	var req structs.PotensiDesaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	potensi := models.PotensiDesa{
		Kategori:    req.Kategori,
		Judul:       req.Judul,
		Deskripsi:   req.Deskripsi,
		FotoID:      req.FotoID,
		Urutan:      req.Urutan,
		IsPublished: req.IsPublished,
	}

	if err := database.DB.Create(&potensi).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah potensi desa")
		return
	}

	helpers.Log(c, "create", "profil", "Menambah potensi desa: "+req.Judul)

	go helpers.RevalidatePath("/potensi")

	database.DB.Preload("Foto").First(&potensi, potensi.ID)
	helpers.Created(c, "Potensi desa berhasil ditambahkan", potensi)
}

func UpdatePotensi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var potensi models.PotensiDesa
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&potensi).Error; err != nil {
		helpers.NotFound(c, "Potensi desa tidak ditemukan")
		return
	}

	var req structs.PotensiDesaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	potensi.Kategori = req.Kategori
	potensi.Judul = req.Judul
	potensi.Deskripsi = req.Deskripsi
	potensi.FotoID = req.FotoID
	potensi.Urutan = req.Urutan
	potensi.IsPublished = req.IsPublished

	if err := database.DB.Save(&potensi).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui potensi desa")
		return
	}

	helpers.Log(c, "update", "profil", "Memperbarui potensi desa: "+req.Judul)

	go helpers.RevalidatePath("/potensi")

	database.DB.Preload("Foto").First(&potensi, potensi.ID)
	helpers.OK(c, "Potensi desa berhasil diperbarui", potensi)
}

func DeletePotensi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var potensi models.PotensiDesa
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&potensi).Error; err != nil {
		helpers.NotFound(c, "Potensi desa tidak ditemukan")
		return
	}

	database.DB.Model(&potensi).Update("is_deleted", true)

	helpers.Log(c, "delete", "profil", "Menghapus potensi desa: "+potensi.Judul)

	go helpers.RevalidatePath("/potensi")

	helpers.OK(c, "Potensi desa berhasil dihapus", nil)
}

func ForceDeletePotensi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var potensi models.PotensiDesa
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&potensi).Error; err != nil {
		helpers.NotFound(c, "Potensi desa tidak ditemukan di trash")
		return
	}

	if potensi.FotoID != nil {
		var media models.Media
		if err := database.DB.First(&media, potensi.FotoID).Error; err == nil {
			go helpers.DeleteFile(media.Path)
			database.DB.Delete(&media)
		}
	}

	database.DB.Delete(&potensi)

	helpers.Log(c, "force_delete", "profil", "Menghapus permanen potensi desa: "+potensi.Judul)

	helpers.OK(c, "Potensi desa berhasil dihapus permanen", nil)
}

func RestorePotensi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var potensi models.PotensiDesa
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&potensi).Error; err != nil {
		helpers.NotFound(c, "Potensi desa tidak ditemukan di trash")
		return
	}

	database.DB.Model(&potensi).Update("is_deleted", false)

	helpers.Log(c, "restore", "profil", "Memulihkan potensi desa: "+potensi.Judul)

	helpers.OK(c, "Potensi desa berhasil dipulihkan", nil)
}

func PublishPotensi(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var potensi models.PotensiDesa
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&potensi).Error; err != nil {
		helpers.NotFound(c, "Potensi desa tidak ditemukan")
		return
	}

	newStatus := !potensi.IsPublished
	database.DB.Model(&potensi).Update("is_published", newStatus)

	go helpers.RevalidatePath("/potensi")

	msg := "Potensi desa berhasil dipublikasikan"
	if !newStatus {
		msg = "Potensi desa berhasil disembunyikan"
	}

	helpers.Log(c, "publish", "profil", msg+": "+potensi.Judul)

	helpers.OK(c, msg, gin.H{"is_published": newStatus})
}
