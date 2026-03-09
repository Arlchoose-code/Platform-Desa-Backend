// Platform Desa — Public Settings Controller
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package public

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/helpers"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

// GetSettingsPublik mengembalikan settings yang boleh diakses publik:
// general, seo, sosmed, hero, footer
// Group sensitif (notifikasi, surat, pengaduan, maintenance) tidak dikembalikan
var allowedGroups = map[string]bool{
	"general": true,
	"seo":     true,
	"sosmed":  true,
	"hero":    true,
	"footer":  true,
}

func GetSettingsPublik(c *gin.Context) {
	group := c.Query("group")

	var settings []models.Setting

	if group != "" {
		if !allowedGroups[group] {
			helpers.Forbidden(c, "Group settings tidak tersedia")
			return
		}
		database.DB.Where("`group` = ?", group).Order("`key` ASC").Find(&settings)
	} else {
		keys := make([]string, 0, len(allowedGroups))
		for k := range allowedGroups {
			keys = append(keys, k)
		}
		database.DB.Where("`group` IN ?", keys).Order("`group` ASC, `key` ASC").Find(&settings)
	}

	result := make(map[string]map[string]interface{})
	for _, s := range settings {
		if _, ok := result[s.Group]; !ok {
			result[s.Group] = make(map[string]interface{})
		}
		result[s.Group][s.Key] = s.Value
	}

	helpers.OK(c, "Berhasil", result)
}
