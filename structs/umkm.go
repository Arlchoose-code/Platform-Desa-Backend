// Platform Desa — Structs UMKM
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type UMKMRequest struct {
	KategoriID  *uint   `json:"kategori_id"`
	NamaUsaha   string  `json:"nama_usaha" binding:"required,min=3"`
	NamaPemilik string  `json:"nama_pemilik" binding:"required,min=2"`
	Deskripsi   *string `json:"deskripsi"`
	Alamat      *string `json:"alamat"`
	Telepon     *string `json:"telepon"`
	WhatsApp    *string `json:"whatsapp"`
	Email       *string `json:"email"`
	Instagram   *string `json:"instagram"`
	FotoID      *uint   `json:"foto_id"`
	IsPublished bool    `json:"is_published"`
}

type ProdukUMKMRequest struct {
	Nama        string  `json:"nama" binding:"required,min=2"`
	Deskripsi   *string `json:"deskripsi"`
	Harga       *int64  `json:"harga"`
	FotoID      *uint   `json:"foto_id"`
	IsAvailable bool    `json:"is_available"`
}

type KategoriUMKMRequest struct {
	Nama string `json:"nama" binding:"required,min=2"`
}
