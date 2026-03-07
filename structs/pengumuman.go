// Platform Desa — Structs Pengumuman & Agenda
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

import "time"

type PengumumanRequest struct {
	Judul          string     `json:"judul" binding:"required,min=3"`
	Isi            string     `json:"isi" binding:"required"`
	FileID         *uint      `json:"file_id"`
	IsPenting      bool       `json:"is_penting"`
	Status         string     `json:"status" binding:"required,oneof=draft published"`
	TanggalMulai   *time.Time `json:"tanggal_mulai"`
	TanggalSelesai *time.Time `json:"tanggal_selesai"`
}

type AgendaRequest struct {
	Judul          string     `json:"judul" binding:"required,min=3"`
	Deskripsi      *string    `json:"deskripsi"`
	Lokasi         *string    `json:"lokasi"`
	TanggalMulai   time.Time  `json:"tanggal_mulai" binding:"required"`
	TanggalSelesai *time.Time `json:"tanggal_selesai"`
	Penyelenggara  *string    `json:"penyelenggara"`
	Status         string     `json:"status" binding:"required,oneof=upcoming ongoing done cancelled"`
	IsPublished    bool       `json:"is_published"`
}
