// Platform Desa — Public Kependudukan Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func GetStatistikPendudukList(c *gin.Context) {
	tahun := c.Query("tahun")

	query := database.DB.Model(&models.StatistikPenduduk{})

	if tahun != "" {
		query = query.Where("tahun = ?", tahun)
	}

	var data []models.StatistikPenduduk
	query.Order("tahun DESC, bulan DESC").Find(&data)

	helpers.OK(c, "Berhasil", data)
}

func GetStatistikPendidikanList(c *gin.Context) {
	var data []models.StatistikPendidikan
	database.DB.Order("tahun DESC").Find(&data)

	helpers.OK(c, "Berhasil", data)
}

func GetStatistikPekerjaanList(c *gin.Context) {
	var data []models.StatistikPekerjaan
	database.DB.Order("tahun DESC").Find(&data)

	helpers.OK(c, "Berhasil", data)
}
