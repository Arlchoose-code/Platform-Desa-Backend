// Platform Desa — Helper Activity Log
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package helpers

import (
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"github.com/gin-gonic/gin"
)

func LogWithUser(userID uint, c *gin.Context, action, module, description string) {
	ip := c.ClientIP()
	log := models.ActivityLog{
		UserID:      &userID,
		Action:      action,
		Module:      module,
		Description: &description,
		IPAddress:   &ip,
	}
	go database.DB.Create(&log)
}

func Log(c *gin.Context, action, module, description string) {
	userID, exists := c.Get("user_id")
	ip := c.ClientIP()

	log := models.ActivityLog{
		Action:      action,
		Module:      module,
		Description: &description,
		IPAddress:   &ip,
	}

	if exists {
		id := userID.(uint)
		log.UserID = &id
	}

	go database.DB.Create(&log)
}
