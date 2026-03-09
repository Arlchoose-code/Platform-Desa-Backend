// Platform Desa — Admin Media Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package admin

import (
	"strconv"

	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetMediaList(c *gin.Context) {
	pg := helpers.GetPagination(c)

	var media []models.Media
	var total int64

	q := database.DB.Model(&models.Media{})

	if folder := c.Query("folder"); folder != "" {
		q = q.Where("folder = ?", folder)
	}
	if t := c.Query("type"); t != "" {
		q = q.Where("type = ?", t)
	}
	if search := c.Query("search"); search != "" {
		q = q.Where("original_name LIKE ?", "%"+search+"%")
	}
	if isExternal := c.Query("is_external"); isExternal != "" {
		q = q.Where("is_external = ?", isExternal == "true")
	}

	q.Count(&total)
	q.Scopes(helpers.Paginate(pg)).Order("created_at DESC").Find(&media)

	helpers.OKPaginated(c, "Berhasil", media, pg.Meta(total))
}
func GetMedia(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var media models.Media
	if err := database.DB.First(&media, id).Error; err != nil {
		helpers.NotFound(c, "Media tidak ditemukan")
		return
	}

	helpers.OK(c, "Berhasil", media)
}

func UploadMedia(c *gin.Context) {
	folder := c.DefaultPostForm("folder", "umum")

	file, err := c.FormFile("file")
	if err != nil {
		helpers.BadRequest(c, "File tidak ditemukan", nil)
		return
	}

	result, err := helpers.UploadFile(file, folder)
	if err != nil {
		helpers.BadRequest(c, err.Error(), nil)
		return
	}

	userID := c.GetUint("user_id")
	media := models.Media{
		Folder:       folder,
		Filename:     result.Filename,
		OriginalName: result.OriginalName,
		MimeType:     result.MimeType,
		Size:         result.Size,
		Path:         result.Path,
		Type:         result.Type,
		IsExternal:   false,
		UploadedBy:   &userID,
	}

	if err := database.DB.Create(&media).Error; err != nil {
		go helpers.DeleteFile(result.Path)
		helpers.InternalError(c, "Gagal menyimpan media")
		return
	}

	helpers.Log(c, "create", "media", "Mengupload media: "+result.OriginalName)

	helpers.Created(c, "Media berhasil diupload", media)
}

func UploadMultipleMedia(c *gin.Context) {
	folder := c.DefaultPostForm("folder", "umum")

	form, err := c.MultipartForm()
	if err != nil {
		helpers.BadRequest(c, "Form tidak valid", nil)
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		helpers.BadRequest(c, "Tidak ada file yang diupload", nil)
		return
	}
	if len(files) > 10 {
		helpers.BadRequest(c, "Maksimal 10 file sekaligus", nil)
		return
	}

	userID := c.GetUint("user_id")
	var uploaded []models.Media
	var failed []string

	for _, file := range files {
		result, err := helpers.UploadFile(file, folder)
		if err != nil {
			failed = append(failed, file.Filename+": "+err.Error())
			continue
		}

		media := models.Media{
			Folder:       folder,
			Filename:     result.Filename,
			OriginalName: result.OriginalName,
			MimeType:     result.MimeType,
			Size:         result.Size,
			Path:         result.Path,
			Type:         result.Type,
			IsExternal:   false,
			UploadedBy:   &userID,
		}

		if err := database.DB.Create(&media).Error; err != nil {
			go helpers.DeleteFile(result.Path)
			failed = append(failed, file.Filename+": gagal simpan ke database")
			continue
		}

		uploaded = append(uploaded, media)
	}

	helpers.OK(c, "Upload selesai", gin.H{
		"uploaded": uploaded,
		"failed":   failed,
		"total":    len(files),
		"success":  len(uploaded),
	})
}

func AddYoutubeMedia(c *gin.Context) {
	var req struct {
		YoutubeURL   string `json:"youtube_url" binding:"required,url"`
		OriginalName string `json:"original_name" binding:"required"`
		Folder       string `json:"folder"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	if req.Folder == "" {
		req.Folder = "video"
	}

	userID := c.GetUint("user_id")
	media := models.Media{
		Folder:       req.Folder,
		Filename:     "",
		OriginalName: req.OriginalName,
		MimeType:     "video/youtube",
		Size:         0,
		Path:         "",
		Type:         "video",
		IsExternal:   true,
		YoutubeURL:   &req.YoutubeURL,
		UploadedBy:   &userID,
	}

	if err := database.DB.Create(&media).Error; err != nil {
		helpers.InternalError(c, "Gagal menyimpan media")
		return
	}

	helpers.Created(c, "Video YouTube berhasil ditambahkan", media)
}

func AddDriveMedia(c *gin.Context) {
	var req struct {
		DriveURL     string `json:"drive_url" binding:"required,url"`
		OriginalName string `json:"original_name" binding:"required"`
		Folder       string `json:"folder"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	if req.Folder == "" {
		req.Folder = "dokumen"
	}

	userID := c.GetUint("user_id")
	media := models.Media{
		Folder:       req.Folder,
		Filename:     "",
		OriginalName: req.OriginalName,
		MimeType:     "application/drive",
		Size:         0,
		Path:         "",
		Type:         "document",
		IsExternal:   true,
		DriveURL:     &req.DriveURL,
		UploadedBy:   &userID,
	}

	if err := database.DB.Create(&media).Error; err != nil {
		helpers.InternalError(c, "Gagal menyimpan media")
		return
	}

	helpers.Created(c, "Dokumen Google Drive berhasil ditambahkan", media)
}

func DeleteMedia(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helpers.BadRequest(c, "ID tidak valid", nil)
		return
	}

	var media models.Media
	if err := database.DB.First(&media, id).Error; err != nil {
		helpers.NotFound(c, "Media tidak ditemukan")
		return
	}

	if !media.IsExternal && media.Path != "" {
		go helpers.DeleteFile(media.Path)
	}

	if err := database.DB.Delete(&media).Error; err != nil {
		helpers.InternalError(c, "Gagal menghapus media")
		return
	}

	helpers.Log(c, "delete", "media", "Menghapus media: "+media.OriginalName)

	helpers.OK(c, "Media berhasil dihapus", nil)
}

func DeleteMultipleMedia(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	var mediaList []models.Media
	database.DB.Where("id IN ?", req.IDs).Find(&mediaList)

	if len(mediaList) == 0 {
		helpers.NotFound(c, "Media tidak ditemukan")
		return
	}

	for _, m := range mediaList {
		if !m.IsExternal && m.Path != "" {
			go helpers.DeleteFile(m.Path)
		}
	}

	database.DB.Where("id IN ?", req.IDs).Delete(&models.Media{})

	helpers.OK(c, "Media berhasil dihapus", gin.H{
		"deleted": len(mediaList),
	})
}
