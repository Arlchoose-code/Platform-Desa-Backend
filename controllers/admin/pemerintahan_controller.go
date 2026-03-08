// Platform Desa — Admin Pemerintahan Controller
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

func GetJabatanList(c *gin.Context) {
	showTrash := c.Query("trash") == "true"
	search := c.Query("search")

	query := database.DB.Model(&models.Jabatan{}).
		Preload("Parent").
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}

	var jabatan []models.Jabatan
	query.Order("urutan ASC, level ASC").Find(&jabatan)

	helpers.OK(c, "Berhasil", jabatan)
}

func GetJabatan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var jabatan models.Jabatan
	if err := database.DB.Preload("Parent").Preload("Children").
		Where("id = ? AND is_deleted = false", id).
		First(&jabatan).Error; err != nil {
		helpers.NotFound(c, "Jabatan tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", jabatan)
}

func CreateJabatan(c *gin.Context) {
	var req structs.JabatanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.Jabatan
	if err := database.DB.Where("nama = ? AND is_deleted = false", req.Nama).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Jabatan sudah ada", nil)
		return
	}

	jabatan := models.Jabatan{
		Nama:     req.Nama,
		ParentID: req.ParentID,
		Level:    req.Level,
		Urutan:   req.Urutan,
	}

	if err := database.DB.Create(&jabatan).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah jabatan")
		return
	}

	database.DB.Preload("Parent").First(&jabatan, jabatan.ID)
	helpers.Created(c, "Jabatan berhasil ditambahkan", jabatan)
}

func UpdateJabatan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var jabatan models.Jabatan
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&jabatan).Error; err != nil {
		helpers.NotFound(c, "Jabatan tidak ditemukan")
		return
	}

	var req structs.JabatanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var existing models.Jabatan
	if err := database.DB.Where("nama = ? AND id != ? AND is_deleted = false", req.Nama, id).First(&existing).Error; err == nil {
		helpers.BadRequest(c, "Jabatan sudah ada", nil)
		return
	}

	jabatan.Nama = req.Nama
	jabatan.ParentID = req.ParentID
	jabatan.Level = req.Level
	jabatan.Urutan = req.Urutan

	if err := database.DB.Save(&jabatan).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui jabatan")
		return
	}

	database.DB.Preload("Parent").First(&jabatan, jabatan.ID)
	helpers.OK(c, "Jabatan berhasil diperbarui", jabatan)
}

func DeleteJabatan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var jabatan models.Jabatan
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&jabatan).Error; err != nil {
		helpers.NotFound(c, "Jabatan tidak ditemukan")
		return
	}

	var count int64
	database.DB.Model(&models.Pejabat{}).Where("jabatan_id = ? AND is_deleted = false", id).Count(&count)
	if count > 0 {
		helpers.BadRequest(c, "Jabatan masih digunakan oleh pejabat aktif", nil)
		return
	}

	database.DB.Model(&jabatan).Update("is_deleted", true)
	helpers.OK(c, "Jabatan berhasil dihapus", nil)
}

func ForceDeleteJabatan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var jabatan models.Jabatan
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&jabatan).Error; err != nil {
		helpers.NotFound(c, "Jabatan tidak ditemukan di trash")
		return
	}

	database.DB.Delete(&jabatan)
	helpers.OK(c, "Jabatan berhasil dihapus permanen", nil)
}

func RestoreJabatan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var jabatan models.Jabatan
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&jabatan).Error; err != nil {
		helpers.NotFound(c, "Jabatan tidak ditemukan di trash")
		return
	}

	database.DB.Model(&jabatan).Update("is_deleted", false)
	helpers.OK(c, "Jabatan berhasil dipulihkan", nil)
}

func GetPejabatList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	jabatanID := c.Query("jabatan_id")
	isActive := c.Query("is_active")
	showTrash := c.Query("trash") == "true"

	query := database.DB.Model(&models.Pejabat{}).
		Preload("Jabatan").
		Preload("Foto").
		Preload("Pendidikan").
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("nama LIKE ?", "%"+search+"%")
	}
	if jabatanID != "" {
		query = query.Where("jabatan_id = ?", jabatanID)
	}
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	var total int64
	query.Count(&total)

	var pejabat []models.Pejabat
	query.Scopes(helpers.Paginate(pg)).Order("urutan ASC, created_at ASC").Find(&pejabat)

	helpers.OKPaginated(c, "Berhasil", pejabat, pg.Meta(total))
}

func GetPejabat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pejabat models.Pejabat
	if err := database.DB.Preload("Jabatan").Preload("Foto").Preload("Pendidikan").
		Where("id = ? AND is_deleted = false", id).
		First(&pejabat).Error; err != nil {
		helpers.NotFound(c, "Pejabat tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", pejabat)
}

func CreatePejabat(c *gin.Context) {
	var req structs.PejabatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var jabatan models.Jabatan
	if err := database.DB.Where("id = ? AND is_deleted = false", req.JabatanID).First(&jabatan).Error; err != nil {
		helpers.BadRequest(c, "Jabatan tidak ditemukan", nil)
		return
	}

	pejabat := models.Pejabat{
		JabatanID:      req.JabatanID,
		JabatanNama:    jabatan.Nama,
		Nama:           req.Nama,
		Slug:           helpers.UniqueSlug(helpers.GenerateSlug(req.Nama)),
		NIP:            req.NIP,
		FotoID:         req.FotoID,
		PeriodeMulai:   req.PeriodeMulai,
		PeriodeSelesai: req.PeriodeSelesai,
		Biodata:        req.Biodata,
		IsActive:       req.IsActive,
		Urutan:         req.Urutan,
	}

	if err := database.DB.Create(&pejabat).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah pejabat")
		return
	}

	if len(req.Pendidikan) > 0 {
		for _, p := range req.Pendidikan {
			database.DB.Create(&models.PejabatPendidikan{
				PejabatID: pejabat.ID,
				Jenjang:   p.Jenjang,
				Jurusan:   p.Jurusan,
				Institusi: p.Institusi,
				Tahun:     p.Tahun,
			})
		}
	}

	go helpers.RevalidatePath("/pemerintahan")

	database.DB.Preload("Jabatan").Preload("Foto").Preload("Pendidikan").First(&pejabat, pejabat.ID)
	helpers.Created(c, "Pejabat berhasil ditambahkan", pejabat)
}

func UpdatePejabat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pejabat models.Pejabat
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&pejabat).Error; err != nil {
		helpers.NotFound(c, "Pejabat tidak ditemukan")
		return
	}

	var req structs.PejabatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var jabatan models.Jabatan
	if err := database.DB.Where("id = ? AND is_deleted = false", req.JabatanID).First(&jabatan).Error; err != nil {
		helpers.BadRequest(c, "Jabatan tidak ditemukan", nil)
		return
	}

	if req.Nama != pejabat.Nama {
		pejabat.Slug = helpers.UniqueSlug(helpers.GenerateSlug(req.Nama))
	}

	pejabat.JabatanID = req.JabatanID
	pejabat.JabatanNama = jabatan.Nama
	pejabat.Nama = req.Nama
	pejabat.NIP = req.NIP
	pejabat.FotoID = req.FotoID
	pejabat.PeriodeMulai = req.PeriodeMulai
	pejabat.PeriodeSelesai = req.PeriodeSelesai
	pejabat.Biodata = req.Biodata
	pejabat.IsActive = req.IsActive
	pejabat.Urutan = req.Urutan

	if err := database.DB.Save(&pejabat).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui pejabat")
		return
	}

	if req.Pendidikan != nil {
		database.DB.Where("pejabat_id = ?", pejabat.ID).Delete(&models.PejabatPendidikan{})
		for _, p := range req.Pendidikan {
			database.DB.Create(&models.PejabatPendidikan{
				PejabatID: pejabat.ID,
				Jenjang:   p.Jenjang,
				Jurusan:   p.Jurusan,
				Institusi: p.Institusi,
				Tahun:     p.Tahun,
			})
		}
	}

	go helpers.RevalidatePath("/pemerintahan")

	database.DB.Preload("Jabatan").Preload("Foto").Preload("Pendidikan").First(&pejabat, pejabat.ID)
	helpers.OK(c, "Pejabat berhasil diperbarui", pejabat)
}

func DeletePejabat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pejabat models.Pejabat
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&pejabat).Error; err != nil {
		helpers.NotFound(c, "Pejabat tidak ditemukan")
		return
	}

	database.DB.Model(&pejabat).Update("is_deleted", true)

	go helpers.RevalidatePath("/pemerintahan")

	helpers.OK(c, "Pejabat berhasil dihapus", nil)
}

func ForceDeletePejabat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pejabat models.Pejabat
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&pejabat).Error; err != nil {
		helpers.NotFound(c, "Pejabat tidak ditemukan di trash")
		return
	}

	database.DB.Where("pejabat_id = ?", pejabat.ID).Delete(&models.PejabatPendidikan{})

	if pejabat.FotoID != nil {
		var media models.Media
		if err := database.DB.First(&media, pejabat.FotoID).Error; err == nil {
			go helpers.DeleteFile(media.Path)
			database.DB.Delete(&media)
		}
	}

	database.DB.Delete(&pejabat)
	helpers.OK(c, "Pejabat berhasil dihapus permanen", nil)
}

func RestorePejabat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var pejabat models.Pejabat
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&pejabat).Error; err != nil {
		helpers.NotFound(c, "Pejabat tidak ditemukan di trash")
		return
	}

	database.DB.Model(&pejabat).Update("is_deleted", false)
	helpers.OK(c, "Pejabat berhasil dipulihkan", nil)
}

func GetLembagaList(c *gin.Context) {
	pg := helpers.GetPagination(c)
	search := c.Query("search")
	isActive := c.Query("is_active")
	showTrash := c.Query("trash") == "true"

	query := database.DB.Model(&models.LembagaDesa{}).
		Preload("Logo").
		Where("is_deleted = ?", showTrash)

	if search != "" {
		query = query.Where("nama LIKE ? OR singkatan LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	var total int64
	query.Count(&total)

	var lembaga []models.LembagaDesa
	query.Scopes(helpers.Paginate(pg)).Order("urutan ASC, created_at ASC").Find(&lembaga)

	helpers.OKPaginated(c, "Berhasil", lembaga, pg.Meta(total))
}

func GetLembaga(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var lembaga models.LembagaDesa
	if err := database.DB.Preload("Logo").
		Where("id = ? AND is_deleted = false", id).
		First(&lembaga).Error; err != nil {
		helpers.NotFound(c, "Lembaga tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", lembaga)
}

func CreateLembaga(c *gin.Context) {
	var req structs.LembagaDesaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	lembaga := models.LembagaDesa{
		Nama:      req.Nama,
		Singkatan: req.Singkatan,
		Deskripsi: req.Deskripsi,
		Ketua:     req.Ketua,
		LogoID:    req.LogoID,
		IsActive:  req.IsActive,
		Urutan:    req.Urutan,
	}

	if err := database.DB.Create(&lembaga).Error; err != nil {
		helpers.InternalError(c, "Gagal menambah lembaga")
		return
	}

	go helpers.RevalidatePath("/pemerintahan")

	database.DB.Preload("Logo").First(&lembaga, lembaga.ID)
	helpers.Created(c, "Lembaga berhasil ditambahkan", lembaga)
}

func UpdateLembaga(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var lembaga models.LembagaDesa
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&lembaga).Error; err != nil {
		helpers.NotFound(c, "Lembaga tidak ditemukan")
		return
	}

	var req structs.LembagaDesaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	lembaga.Nama = req.Nama
	lembaga.Singkatan = req.Singkatan
	lembaga.Deskripsi = req.Deskripsi
	lembaga.Ketua = req.Ketua
	lembaga.LogoID = req.LogoID
	lembaga.IsActive = req.IsActive
	lembaga.Urutan = req.Urutan

	if err := database.DB.Save(&lembaga).Error; err != nil {
		helpers.InternalError(c, "Gagal memperbarui lembaga")
		return
	}

	go helpers.RevalidatePath("/pemerintahan")

	database.DB.Preload("Logo").First(&lembaga, lembaga.ID)
	helpers.OK(c, "Lembaga berhasil diperbarui", lembaga)
}

func DeleteLembaga(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var lembaga models.LembagaDesa
	if err := database.DB.Where("id = ? AND is_deleted = false", id).First(&lembaga).Error; err != nil {
		helpers.NotFound(c, "Lembaga tidak ditemukan")
		return
	}

	database.DB.Model(&lembaga).Update("is_deleted", true)

	go helpers.RevalidatePath("/pemerintahan")

	helpers.OK(c, "Lembaga berhasil dihapus", nil)
}

func ForceDeleteLembaga(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var lembaga models.LembagaDesa
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&lembaga).Error; err != nil {
		helpers.NotFound(c, "Lembaga tidak ditemukan di trash")
		return
	}

	if lembaga.LogoID != nil {
		var media models.Media
		if err := database.DB.First(&media, lembaga.LogoID).Error; err == nil {
			go helpers.DeleteFile(media.Path)
			database.DB.Delete(&media)
		}
	}

	database.DB.Delete(&lembaga)
	helpers.OK(c, "Lembaga berhasil dihapus permanen", nil)
}

func RestoreLembaga(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var lembaga models.LembagaDesa
	if err := database.DB.Where("id = ? AND is_deleted = true", id).First(&lembaga).Error; err != nil {
		helpers.NotFound(c, "Lembaga tidak ditemukan di trash")
		return
	}

	database.DB.Model(&lembaga).Update("is_deleted", false)
	helpers.OK(c, "Lembaga berhasil dipulihkan", nil)
}
