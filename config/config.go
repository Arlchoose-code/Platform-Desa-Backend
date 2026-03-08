// Platform Desa — Config
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Name    string
	Env     string
	Port    string
	URL     string
	Version string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type JWTConfig struct {
	Secret             string
	ExpireHours        int
	RefreshExpireHours int
}

type UploadConfig struct {
	Path            string
	MaxImageSize    int64 // bytes
	MaxVideoSize    int64 // bytes, 0 = bebas
	MaxDocumentSize int64 // bytes, 0 = bebas
	ImageMaxWidth   int   // px, resize kalau melebihi
	ImageMaxHeight  int   // px, resize kalau melebihi
	ImageQuality    int   // WebP quality 1-100
}

type CORSConfig struct {
	AllowedOrigins string
}

type NextJSConfig struct {
	URL              string
	RevalidateSecret string
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
	CORS     CORSConfig
	NextJS   NextJSConfig
}

var App *Config

func Load() error {
	_ = godotenv.Load()

	jwtExpire, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	jwtRefreshExpire, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRE_HOURS", "168"))

	maxImageMB, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_IMAGE_MB", "10"), 10, 64)
	maxVideoMB, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_VIDEO_MB", "500"), 10, 64)
	maxDocumentMB, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_DOCUMENT_MB", "0"), 10, 64)

	imageMaxWidth, _ := strconv.Atoi(getEnv("IMAGE_MAX_WIDTH", "1920"))
	imageMaxHeight, _ := strconv.Atoi(getEnv("IMAGE_MAX_HEIGHT", "1920"))
	imageQuality, _ := strconv.Atoi(getEnv("IMAGE_QUALITY", "80"))

	App = &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "Platform Desa"),
			Env:     getEnv("APP_ENV", "production"),
			Port:    getEnv("APP_PORT", "8080"),
			URL:     getEnv("APP_URL", "http://localhost:8080"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			Name:     getEnv("DB_NAME", "platform_desa"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret:             getEnv("JWT_SECRET", ""),
			ExpireHours:        jwtExpire,
			RefreshExpireHours: jwtRefreshExpire,
		},
		Upload: UploadConfig{
			Path:            getEnv("UPLOAD_PATH", "./uploads"),
			MaxImageSize:    maxImageMB * 1024 * 1024,
			MaxVideoSize:    maxVideoMB * 1024 * 1024,
			MaxDocumentSize: maxDocumentMB * 1024 * 1024,
			ImageMaxWidth:   imageMaxWidth,
			ImageMaxHeight:  imageMaxHeight,
			ImageQuality:    imageQuality,
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
		},
		NextJS: NextJSConfig{
			URL:              getEnv("NEXTJS_URL", "http://localhost:3000"),
			RevalidateSecret: getEnv("NEXTJS_REVALIDATE_SECRET", ""),
		},
	}

	if App.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET wajib diisi di file .env")
	}

	return nil
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
