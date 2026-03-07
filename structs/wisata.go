// Platform Desa — Structs Wisata
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type WisataRequest struct {
	KategoriID  *uint    `json:"kategori_id"`
	Nama        string   `json:"nama" binding:"required,min=3"`
	Deskripsi   *string  `json:"deskripsi"`
	Alamat      *string  `json:"alamat"`
	Latitude    *float64 `json:"latitude"`
	Longitude   *float64 `json:"longitude"`
	JamBuka     *string  `json:"jam_buka"`
	JamTutup    *string  `json:"jam_tutup"`
	HargaTiket  int64    `json:"harga_tiket"`
	Kontak      *string  `json:"kontak"`
	ThumbnailID *uint    `json:"thumbnail_id"`
	Fasilitas   *string  `json:"fasilitas"`
	IsPublished bool     `json:"is_published"`
}

type KategoriWisataRequest struct {
	Nama string  `json:"nama" binding:"required,min=2"`
	Icon *string `json:"icon"`
}

type WisataGaleriRequest struct {
	MediaID uint    `json:"media_id" binding:"required"`
	Caption *string `json:"caption"`
	Urutan  uint    `json:"urutan"`
}

type UrutanRequest struct {
	Items []struct {
		ID     uint `json:"id" binding:"required"`
		Urutan uint `json:"urutan"`
	} `json:"items" binding:"required"`
}
