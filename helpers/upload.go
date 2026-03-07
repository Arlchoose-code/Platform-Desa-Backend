// Platform Desa — Upload Helper
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/Arlchoose-code/platform-desa-backend/config"
	"github.com/google/uuid"
)

var allowedImages = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

var allowedVideos = map[string]string{
	"video/mp4":       ".mp4",
	"video/webm":      ".webm",
	"video/quicktime": ".mov",
}

var allowedDocuments = map[string]string{
	"application/pdf":    ".pdf",
	"application/msword": ".doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	"application/vnd.ms-excel": ".xls",
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
}

type UploadResult struct {
	Filename     string
	OriginalName string
	MimeType     string
	Size         int64
	Path         string
	Type         string
}

func UploadFile(file *multipart.FileHeader, folder string) (*UploadResult, error) {
	if file.Size > config.App.Upload.MaxUploadSize {
		return nil, fmt.Errorf("ukuran file maksimal %dMB",
			config.App.Upload.MaxUploadSize/1024/1024)
	}

	mimeType := file.Header.Get("Content-Type")
	fileType := ""
	ext := ""

	if e, ok := allowedImages[mimeType]; ok {
		fileType, ext = "image", e
	} else if e, ok := allowedVideos[mimeType]; ok {
		fileType, ext = "video", e
	} else if e, ok := allowedDocuments[mimeType]; ok {
		fileType, ext = "document", e
	} else {
		return nil, fmt.Errorf("tipe file tidak diizinkan: %s", mimeType)
	}

	dir := filepath.Join(config.App.Upload.Path, folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("gagal buat folder: %w", err)
	}

	filename := uuid.New().String() + ext
	fullPath := filepath.Join(dir, filename)

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	relativePath := strings.ReplaceAll(
		filepath.Join("uploads", folder, filename), "\\", "/",
	)

	return &UploadResult{
		Filename:     filename,
		OriginalName: file.Filename,
		MimeType:     mimeType,
		Size:         file.Size,
		Path:         relativePath,
		Type:         fileType,
	}, nil
}

func DeleteFile(path string) error {
	full := filepath.Join(".", path)
	if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
