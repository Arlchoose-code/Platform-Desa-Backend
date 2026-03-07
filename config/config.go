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
	Path          string
	MaxUploadSize int64 // bytes
}

type CORSConfig struct {
	AllowedOrigins string
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
	CORS     CORSConfig
}

var App *Config

func Load() error {
	_ = godotenv.Load()

	jwtExpire, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	jwtRefreshExpire, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRE_HOURS", "168"))
	maxUploadMB, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "10"), 10, 64)

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
			Path:          getEnv("UPLOAD_PATH", "./uploads"),
			MaxUploadSize: maxUploadMB * 1024 * 1024,
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
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
