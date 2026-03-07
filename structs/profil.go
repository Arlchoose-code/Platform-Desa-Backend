// Platform Desa — Structs Profil Desa
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type ProfilDesaRequest struct {
	NamaDesa      string   `json:"nama_desa" binding:"required,min=3"`
	NamaKecamatan *string  `json:"nama_kecamatan"`
	NamaKabupaten *string  `json:"nama_kabupaten"`
	NamaProvinsi  *string  `json:"nama_provinsi"`
	KodePos       *string  `json:"kode_pos"`
	Alamat        *string  `json:"alamat"`
	Telepon       *string  `json:"telepon"`
	Email         *string  `json:"email"`
	Website       *string  `json:"website"`
	LuasWilayah   *float64 `json:"luas_wilayah"`
	JumlahDusun   uint     `json:"jumlah_dusun"`
	JumlahRW      uint     `json:"jumlah_rw"`
	JumlahRT      uint     `json:"jumlah_rt"`
	BatasUtara    *string  `json:"batas_utara"`
	BatasSelatan  *string  `json:"batas_selatan"`
	BatasTimur    *string  `json:"batas_timur"`
	BatasBarat    *string  `json:"batas_barat"`
	Latitude      *float64 `json:"latitude"`
	Longitude     *float64 `json:"longitude"`
	LogoID        *uint    `json:"logo_id"`
	FotoDesaID    *uint    `json:"foto_desa_id"`
	Sejarah       *string  `json:"sejarah"`
	Visi          *string  `json:"visi"`
	Misi          *string  `json:"misi"`
}

type PotensiDesaRequest struct {
	Kategori    string  `json:"kategori" binding:"required,oneof=sda sdm infrastruktur budaya lainnya"`
	Judul       string  `json:"judul" binding:"required,min=3"`
	Deskripsi   *string `json:"deskripsi"`
	FotoID      *uint   `json:"foto_id"`
	Urutan      uint    `json:"urutan"`
	IsPublished bool    `json:"is_published"`
}
