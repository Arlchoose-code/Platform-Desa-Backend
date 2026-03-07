// Platform Desa — Structs Pemerintahan
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type JabatanRequest struct {
	Nama     string `json:"nama" binding:"required,min=2"`
	ParentID *uint  `json:"parent_id"`
	Level    uint   `json:"level"`
	Urutan   uint   `json:"urutan"`
}

type PejabatRequest struct {
	JabatanID      uint    `json:"jabatan_id" binding:"required"`
	Nama           string  `json:"nama" binding:"required,min=2"`
	NIP            *string `json:"nip"`
	FotoID         *uint   `json:"foto_id"`
	Pendidikan     *string `json:"pendidikan"`
	PeriodeMulai   *int    `json:"periode_mulai"`
	PeriodeSelesai *int    `json:"periode_selesai"`
	Biodata        *string `json:"biodata"`
	IsActive       bool    `json:"is_active"`
	Urutan         uint    `json:"urutan"`
}

type LembagaDesaRequest struct {
	Nama      string  `json:"nama" binding:"required,min=3"`
	Singkatan *string `json:"singkatan"`
	Deskripsi *string `json:"deskripsi"`
	Ketua     *string `json:"ketua"`
	LogoID    *uint   `json:"logo_id"`
	IsActive  bool    `json:"is_active"`
	Urutan    uint    `json:"urutan"`
}
