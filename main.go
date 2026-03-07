// Platform Desa — Main Entry Point
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package main

import (
	"fmt"

	"github.com/Arlchoose-code/platform-desa-backend/config"
	"github.com/Arlchoose-code/platform-desa-backend/database"
	"github.com/Arlchoose-code/platform-desa-backend/middlewares"
	"github.com/Arlchoose-code/platform-desa-backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	config.Load()

	// Connect database + AutoMigrate
	database.InitDB()

	// Setup Gin
	if config.App.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.CORS())

	// Serve uploads
	r.Static("/uploads", config.App.Upload.Path)

	// Register routes
	routes.Register(r)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"version": config.App.App.Version,
		})
	})

	addr := fmt.Sprintf(":%s", config.App.App.Port)
	fmt.Printf("✓ Platform Desa Backend berjalan di http://localhost%s\n", addr)

	r.Run(addr)
}
