// Platform Desa — Admin UMKM Controller
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

func GetKategoriUMKMList(c *gin.Context) {
	showTrash := c.Query("trash") == "true"
	search := c.Query("search")

	query := database.DB.Model(&models.KategoriUMKM{}).Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}

	var kategori []models.KategoriUMKM
	query.Order("nama ASC").Find(&kategori)

	helpers.OK(c, "Berhasil", kategori)
}

func GetKategoriUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriUMKM
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", kategori)
}

func CreateKategoriUMKM(c *gin.Context) {
	var req structs.KategoriUMKMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.KategoriUMKM
	if err := database.DB.Where("nama = ? AND is_deleted = false", req.Nama).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kategori sudah ada", nil)
		return
	}

	kategori := models.KategoriUMKM{
		Nama: req.Nama,
		Slug: helpers.UniqueSlug(helpers.GenerateSlug(req.Nama)),
	}

	if err := database.DB.Create(&kategori).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah kategori")
		return
	}

	helpers.Created(c, "Kategori berhasil ditambahkan", kategori)
}

func UpdateKategoriUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriUMKM
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	var req structs.KategoriUMKMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.KategoriUMKM
	if err := database.DB.Where("nama = ? AND id != ? AND is_deleted = false", req.Nama, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Kategori sudah ada", nil)
		return
	}

	if req.Nama != kategori.Nama {
		kategori.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Nama))
	}
	kategori.Nama = req.Nama

	if err := database.DB.Save(&kategori).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui kategori")
		return
	}

	helpers.OK(c, "Kategori berhasil diperbarui", kategori)
}

func DeleteKategoriUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriUMKM
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan")
		return
	}

	var count int64
	database.DB.Model(&models.UMKM{}).Where("kategori_id = ? AND is_deleted = false", id).Count(&count)
	if count > 0 {
		helpers.BadRequest(c, "Kategori masih digunakan oleh UMKM aktif", nil)
		return
	}

	database.DB.Model(&kategori).Update("is_deleted", true)
	helpers.OK(c, "Kategori berhasil dihapus", nil)
}

func ForceDeleteKategoriUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriUMKM
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan di trash")
		return
	}

	database.DB.Model(&models.UMKM{}).Where("kategori_id = ?", id).Updates(map[string]interface{}{
		"kategori_id":   nil,
		"kategori_nama": nil,
	})

	database.DB.Delete(&kategori)
	helpers.OK(c, "Kategori berhasil dihapus permanen", nil)
}

func RestoreKategoriUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var kategori models.KategoriUMKM
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&kategori).Error; err != nil {
		helpers.NotFound(c, "Kategori tidak ditemukan di trash")
		return
	}

	database.DB.Model(&kategori).Update("is_deleted", false)
	helpers.OK(c, "Kategori berhasil dipulihkan", nil)
}

func GetUMKMList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	kategoriID := c.Query("kategori_id")
	isPublished := c.Query("is_published")
	showTrash := c.Query("trash") == "true"

	query := database.DB.Model(&models.UMKM{}).
		Preload("Kategori").
		Preload("Foto").
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("nama_usaha LIKE ? OR nama_pemilik LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if kategoriID != "" {
		query = query.Where("kategori_id = ?", kategoriID)
	}
	if isPublished != "" {
		query = query.Where("is_published = ?", isPublished == "true")
	}

	var total int64
	query.Count(&total)

	var umkm []models.UMKM
	query.Scopes(helpers.Paginate(pg)).Order("created_at DESC").Find(&umkm)

	helpers.OKPaginated(c, "Berhasil", umkm, pg.Meta(total))
}

func GetUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var umkm models.UMKM
	if err := database.DB.Preload("Kategori").Preload("Foto").Preload("Produk.Foto").
		Where("id = ? AND is_deleted = false", id).
		First(&umkm).Error; err != nil {
		helpers.NotFound(c, "UMKM tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", umkm)
}

func CreateUMKM(c *gin.Context) {
	var req structs.UMKMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	umkm := models.UMKM{
		NamaUsaha:   req.NamaUsaha,
		Slug:        helpers.UniqueSlug(helpers.GenerateSlug(req.NamaUsaha)),
		NamaPemilik: req.NamaPemilik,
		Deskripsi:   req.Deskripsi,
		Alamat:      req.Alamat,
		Telepon:     req.Telepon,
		WhatsApp:    req.WhatsApp,
		Email:       req.Email,
		Instagram:   req.Instagram,
		FotoID:      req.FotoID,
		IsPublished: req.IsPublished,
	}

	if req.KategoriID != nil {
		var kategori models.KategoriUMKM
		if err := database.DB.Where("id = ? AND is_deleted = false", req.KategoriID).First(&kategori).Error; err == nil {
			umkm.KategoriID = req.KategoriID
			umkm.KategoriNama = &kategori.Nama
		}
	}

	if err := database.DB.Create(&umkm).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah UMKM")
		return
	}

	if req.IsPublished {
		go helpers.RevalidatePaths("/umkm", "/umkm/"+umkm.Slug)
	}

	database.DB.Preload("Kategori").Preload("Foto").First(&umkm, umkm.ID)
	helpers.Created(c, "UMKM berhasil ditambahkan", umkm)
}

func UpdateUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var umkm models.UMKM
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&umkm).Error; err != nil {
		helpers.NotFound(c, "UMKM tidak ditemukan")
		return
	}

	var req structs.UMKMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	if req.NamaUsaha != umkm.NamaUsaha {
		umkm.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.NamaUsaha))
	}

	umkm.NamaUsaha = req.NamaUsaha
	umkm.NamaPemilik = req.NamaPemilik
	umkm.Deskripsi = req.Deskripsi
	umkm.Alamat = req.Alamat
	umkm.Telepon = req.Telepon
	umkm.WhatsApp = req.WhatsApp
	umkm.Email = req.Email
	umkm.Instagram = req.Instagram
	umkm.FotoID = req.FotoID
	umkm.IsPublished = req.IsPublished

	if req.KategoriID != nil {
		var kategori models.KategoriUMKM
		if err := database.DB.Where("id = ? AND is_deleted = false", req.KategoriID).First(&kategori).Error; err == nil {
			umkm.KategoriID = req.KategoriID
			umkm.KategoriNama = &kategori.Nama
		}
	} else {
		umkm.KategoriID = nil
		umkm.KategoriNama = nil
	}

	if err := database.DB.Save(&umkm).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui UMKM")
		return
	}

	go helpers.RevalidatePaths("/umkm", "/umkm/"+umkm.Slug)

	database.DB.Preload("Kategori").Preload("Foto").First(&umkm, umkm.ID)
	helpers.OK(c, "UMKM berhasil diperbarui", umkm)
}

func DeleteUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var umkm models.UMKM
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&umkm).Error; err != nil {
		helpers.NotFound(c, "UMKM tidak ditemukan")
		return
	}

	database.DB.Model(&umkm).Update("is_deleted", true)

	go helpers.RevalidatePath("/umkm")

	helpers.OK(c, "UMKM berhasil dihapus", nil)
}

func ForceDeleteUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var umkm models.UMKM
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&umkm).Error; err != nil {
		helpers.NotFound(c, "UMKM tidak ditemukan di trash")
		return
	}

	var produk []models.ProdukUMKM
	database.DB.Where("umkm_id = ?", umkm.ID).Preload("Foto").Find(&produk)
	for _, p := range produk {
		if p.Foto != nil {
			go helpers.DeleteFile(p.Foto.Path)
			database.DB.Delete(&p.Foto)
		}
	}
	database.DB.Where("umkm_id = ?", umkm.ID).Delete(&models.ProdukUMKM{})

	if umkm.FotoID != nil {
		var foto models.Media
		if err := database.DB.First(&foto, umkm.FotoID).Error; err == nil {
			go helpers.DeleteFile(foto.Path)
			database.DB.Delete(&foto)
		}
	}

	database.DB.Delete(&umkm)
	helpers.OK(c, "UMKM berhasil dihapus permanen", nil)
}

func RestoreUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var umkm models.UMKM
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&umkm).Error; err != nil {
		helpers.NotFound(c, "UMKM tidak ditemukan di trash")
		return
	}

	database.DB.Model(&umkm).Update("is_deleted", false)
	helpers.OK(c, "UMKM berhasil dipulihkan", nil)
}

func PublishUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var umkm models.UMKM
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&umkm).Error; err != nil {
		helpers.NotFound(c, "UMKM tidak ditemukan")
		return
	}

	newStatus := !umkm.IsPublished
	database.DB.Model(&umkm).Update("is_published", newStatus)

	go helpers.RevalidatePaths("/umkm", "/umkm/"+umkm.Slug)

	msg := "UMKM berhasil dipublikasikan"
	if !newStatus {
		msg = "UMKM berhasil disembunyikan"
	}
	helpers.OK(c, msg, gin.H{"is_published": newStatus})
}

func GetProdukUMKMList(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var umkm models.UMKM
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&umkm).Error; err != nil {
		helpers.NotFound(c, "UMKM tidak ditemukan")
		return
	}

	var produk []models.ProdukUMKM
	database.DB.Preload("Foto").Where("umkm_id = ?", id).Order("created_at ASC").Find(&produk)

	helpers.OK(c, "Berhasil", produk)
}

func CreateProdukUMKM(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var umkm models.UMKM
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&umkm).Error; err != nil {
		helpers.NotFound(c, "UMKM tidak ditemukan")
		return
	}

	var req structs.ProdukUMKMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	produk := models.ProdukUMKM{
		UMKMID:      umkm.ID,
		Nama:        req.Nama,
		Deskripsi:   req.Deskripsi,
		Harga:       req.Harga,
		FotoID:      req.FotoID,
		IsAvailable: req.IsAvailable,
	}

	if err := database.DB.Create(&produk).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah produk")
		return
	}

	go helpers.RevalidatePaths("/umkm", "/umkm/"+umkm.Slug)

	database.DB.Preload("Foto").First(&produk, produk.ID)
	helpers.Created(c, "Produk berhasil ditambahkan", produk)
}

func UpdateProdukUMKM(c *gin.Context) {
	umkmID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	produkID, err := strconv.Atoi(c.Param("produk_id"))
	if err != nil {
		helpers.BadRequest(c, "ID produk tidak valid", nil)
		return
	}

	var produk models.ProdukUMKM
	if err := database.DB.Where("id = ? AND umkm_id = ?", produkID, umkmID).First(&produk).Error; err != nil {
		helpers.NotFound(c, "Produk tidak ditemukan")
		return
	}

	var req structs.ProdukUMKMRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	produk.Nama = req.Nama
	produk.Deskripsi = req.Deskripsi
	produk.Harga = req.Harga
	produk.FotoID = req.FotoID
	produk.IsAvailable = req.IsAvailable

	if err := database.DB.Save(&produk).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui produk")
		return
	}

	var umkm models.UMKM
	if err := database.DB.First(&umkm, umkmID).Error; err == nil {
		go helpers.RevalidatePaths("/umkm", "/umkm/"+umkm.Slug)
	}

	database.DB.Preload("Foto").First(&produk, produk.ID)
	helpers.OK(c, "Produk berhasil diperbarui", produk)
}

func DeleteProdukUMKM(c *gin.Context) {
	umkmID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	produkID, err := strconv.Atoi(c.Param("produk_id"))
	if err != nil {
		helpers.BadRequest(c, "ID produk tidak valid", nil)
		return
	}

	var produk models.ProdukUMKM
	if err := database.DB.Preload("Foto").Where("id = ? AND umkm_id = ?", produkID, umkmID).First(&produk).Error; err != nil {
		helpers.NotFound(c, "Produk tidak ditemukan")
		return
	}

	if produk.Foto != nil {
		go helpers.DeleteFile(produk.Foto.Path)
		database.DB.Delete(&produk.Foto)
	}

	database.DB.Delete(&produk)

	var umkm models.UMKM
	if err := database.DB.First(&umkm, umkmID).Error; err == nil {
		go helpers.RevalidatePaths("/umkm", "/umkm/"+umkm.Slug)
	}

	helpers.OK(c, "Produk berhasil dihapus", nil)
}
