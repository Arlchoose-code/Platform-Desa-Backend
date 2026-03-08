// Platform Desa — Upload Helper
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package helpers

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/Arlchoose-code/platform-desa-backend/config"
	"github.com/disintegration/imaging"
	"github.com/gen2brain/webp"
	"github.com/google/uuid"
)

var allowedImages = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
	"image/gif":  true,
}

var allowedVideos = map[string]bool{
	"video/mp4":       true,
	"video/webm":      true,
	"video/quicktime": true,
}

var allowedDocuments = map[string]bool{
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.ms-excel": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
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
	mimeType := file.Header.Get("Content-Type")

	var fileType string
	switch {
	case allowedImages[mimeType]:
		fileType = "image"
		if config.App.Upload.MaxImageSize > 0 && file.Size > config.App.Upload.MaxImageSize {
			return nil, fmt.Errorf("ukuran foto maksimal %dMB",
				config.App.Upload.MaxImageSize/1024/1024)
		}
	case allowedVideos[mimeType]:
		fileType = "video"
		if config.App.Upload.MaxVideoSize > 0 && file.Size > config.App.Upload.MaxVideoSize {
			return nil, fmt.Errorf("ukuran video maksimal %dMB",
				config.App.Upload.MaxVideoSize/1024/1024)
		}
	case allowedDocuments[mimeType]:
		fileType = "document"
		if config.App.Upload.MaxDocumentSize > 0 && file.Size > config.App.Upload.MaxDocumentSize {
			return nil, fmt.Errorf("ukuran dokumen maksimal %dMB",
				config.App.Upload.MaxDocumentSize/1024/1024)
		}
	default:
		return nil, fmt.Errorf("tipe file tidak diizinkan: %s", mimeType)
	}

	dir := filepath.Join(config.App.Upload.Path, folder)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("gagal buat folder: %w", err)
	}

	if fileType == "image" {
		return uploadImage(file, folder, dir)
	}

	return uploadRaw(file, folder, dir, mimeType, fileType)
}

func uploadImage(file *multipart.FileHeader, folder, dir string) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("gagal buka file: %w", err)
	}
	defer src.Close()

	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("gagal decode gambar: %w", err)
	}

	maxW := config.App.Upload.ImageMaxWidth
	maxH := config.App.Upload.ImageMaxHeight
	bounds := img.Bounds()
	if bounds.Dx() > maxW || bounds.Dy() > maxH {
		img = imaging.Fit(img, maxW, maxH, imaging.Lanczos)
	}

	var buf bytes.Buffer
	options := webp.Options{
		Lossless: false,
		Quality:  config.App.Upload.ImageQuality,
	}
	if err := webp.Encode(&buf, img, options); err != nil {
		return nil, fmt.Errorf("gagal encode WebP: %w", err)
	}

	filename := uuid.New().String() + ".webp"
	fullPath := filepath.Join(dir, filename)
	if err := os.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
		return nil, fmt.Errorf("gagal simpan file: %w", err)
	}

	relativePath := strings.ReplaceAll(
		filepath.Join("uploads", folder, filename), "\\", "/",
	)

	return &UploadResult{
		Filename:     filename,
		OriginalName: file.Filename,
		MimeType:     "image/webp",
		Size:         int64(buf.Len()),
		Path:         relativePath,
		Type:         "image",
	}, nil
}

func uploadRaw(file *multipart.FileHeader, folder, dir, mimeType, fileType string) (*UploadResult, error) {
	ext := mimeToExt(mimeType)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("gagal buka file: %w", err)
	}
	defer src.Close()

	filename := uuid.New().String() + ext
	fullPath := filepath.Join(dir, filename)

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, fmt.Errorf("gagal baca file: %w", err)
	}

	if err := os.WriteFile(fullPath, buf.Bytes(), 0644); err != nil {
		return nil, fmt.Errorf("gagal simpan file: %w", err)
	}

	relativePath := strings.ReplaceAll(
		filepath.Join("uploads", folder, filename), "\\", "/",
	)

	return &UploadResult{
		Filename:     filename,
		OriginalName: file.Filename,
		MimeType:     mimeType,
		Size:         int64(buf.Len()),
		Path:         relativePath,
		Type:         fileType,
	}, nil
}

func DeleteFile(path string) error {
	if path == "" {
		return nil
	}
	full := filepath.Join(".", path)
	if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func mimeToExt(mimeType string) string {
	m := map[string]string{
		"video/mp4":          ".mp4",
		"video/webm":         ".webm",
		"video/quicktime":    ".mov",
		"application/pdf":    ".pdf",
		"application/msword": ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"application/vnd.ms-excel": ".xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": ".xlsx",
	}
	if ext, ok := m[mimeType]; ok {
		return ext
	}
	return ".bin"
}
