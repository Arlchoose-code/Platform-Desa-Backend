// Platform Desa — Structs APBDes
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type APBDesRequest struct {
	Tahun       int   `json:"tahun" binding:"required"`
	DokumenID   *uint `json:"dokumen_id"`
	IsPublished bool  `json:"is_published"`
}

type APBDesDetailRequest struct {
	Jenis     string `json:"jenis" binding:"required,oneof=pendapatan belanja pembiayaan"`
	Kategori  string `json:"kategori" binding:"required"`
	Uraian    string `json:"uraian" binding:"required"`
	Anggaran  int64  `json:"anggaran"`
	Realisasi int64  `json:"realisasi"`
	Urutan    uint   `json:"urutan"`
}
