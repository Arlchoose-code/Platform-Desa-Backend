// Platform Desa — Structs Persuratan
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type JenisSuratRequest struct {
	Nama           string  `json:"nama" binding:"required,min=3"`
	Kode           string  `json:"kode" binding:"required"`
	Deskripsi      *string `json:"deskripsi"`
	Syarat         *string `json:"syarat"`
	TemplateFileID *uint   `json:"template_file_id"`
	EstimasiHari   uint    `json:"estimasi_hari"`
	IsActive       bool    `json:"is_active"`
}

type AjukanSuratRequest struct {
	JenisSuratID uint    `json:"jenis_surat_id" binding:"required"`
	NamaPemohon  string  `json:"nama_pemohon" binding:"required,min=3"`
	NIK          string  `json:"nik" binding:"required,len=16"`
	Keperluan    *string `json:"keperluan"`
	DataTambahan *string `json:"data_tambahan"`
}

type ProsesSuratRequest struct {
	CatatanAdmin *string `json:"catatan_admin"`
}

type SelesaikanSuratRequest struct {
	FileHasilID  uint    `json:"file_hasil_id" binding:"required"`
	CatatanAdmin *string `json:"catatan_admin"`
}

type TolakSuratRequest struct {
	Alasan string `json:"alasan" binding:"required"`
}
