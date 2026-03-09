// Platform Desa — Admin Settings Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package admin

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/Arlchoose-code/platform-desa-backend/structs"
	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	group := c.Query("group")

	query := database.DB.Model(&models.Setting{})
	if group != "" {
		query = query.Where("`group` = ?", group)
	}

	var settings []models.Setting
	query.Order("`group` ASC, `key` ASC").Find(&settings)

	result := make(map[string]map[string]interface{})
	for _, s := range settings {
		if _, ok := result[s.Group]; !ok {
			result[s.Group] = make(map[string]interface{})
		}
		result[s.Group][s.Key] = s.Value
	}

	helpers.OK(c, "Berhasil", result)
}

func UpdateSettings(c *gin.Context) {
	var req structs.SettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	for key, value := range req.Settings {
		val := value
		var existing models.Setting
		if err := database.DB.Where("`key` = ?", key).First(&existing).Error; err == nil {
			database.DB.Model(&existing).Update("value", &val)
		} else {
			database.DB.Create(&models.Setting{
				Key:   key,
				Value: &val,
				Group: "general",
			})
		}
	}

	helpers.Log(c, "update", "settings", "Memperbarui settings")

	go helpers.RevalidatePath("/")

	helpers.OK(c, "Settings berhasil diperbarui", nil)
}

func UpdateSettingsByGroup(c *gin.Context) {
	group := c.Param("group")

	var req structs.SettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.BadRequest(c, "Data tidak valid", err.Error())
		return
	}

	for key, value := range req.Settings {
		val := value
		var existing models.Setting
		if err := database.DB.Where("`key` = ?", key).First(&existing).Error; err == nil {
			database.DB.Model(&existing).Update("value", &val)
		} else {
			database.DB.Create(&models.Setting{
				Key:   key,
				Value: &val,
				Group: group,
			})
		}
	}

	helpers.Log(c, "update", "settings", "Memperbarui settings group: "+group)

	go helpers.RevalidatePath("/")

	helpers.OK(c, "Settings berhasil diperbarui", nil)
}
