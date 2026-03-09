// Platform Desa — Structs Sistem (Peta, FAQ, Settings)
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type PetaFasilitasRequest struct {
	Nama        string  `json:"nama" binding:"required,min=3"`
	Kategori    string  `json:"kategori" binding:"required,oneof=pendidikan kesehatan ibadah pemerintahan olahraga lainnya"`
	Alamat      *string `json:"alamat"`
	Latitude    float64 `json:"latitude" binding:"required"`
	Longitude   float64 `json:"longitude" binding:"required"`
	Deskripsi   *string `json:"deskripsi"`
	FotoID      *uint   `json:"foto_id"`
	IsPublished bool    `json:"is_published"`
}

type FAQRequest struct {
	Pertanyaan  string  `json:"pertanyaan" binding:"required,min=5"`
	Jawaban     string  `json:"jawaban" binding:"required,min=5"`
	Kategori    *string `json:"kategori"`
	Urutan      uint    `json:"urutan"`
	IsPublished bool    `json:"is_published"`
}

type SettingsRequest struct {
	Settings map[string]string `json:"settings" binding:"required"`
}
