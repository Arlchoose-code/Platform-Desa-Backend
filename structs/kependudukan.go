// Platform Desa — Structs Kependudukan
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package structs

type StatistikPendudukRequest struct {
	Tahun         int  `json:"tahun" binding:"required"`
	Bulan         int  `json:"bulan" binding:"required,min=1,max=12"`
	TotalPenduduk uint `json:"total_penduduk"`
	LakiLaki      uint `json:"laki_laki"`
	Perempuan     uint `json:"perempuan"`
	TotalKK       uint `json:"total_kk"`
	Kelahiran     uint `json:"kelahiran"`
	Kematian      uint `json:"kematian"`
	PindahMasuk   uint `json:"pindah_masuk"`
	PindahKeluar  uint `json:"pindah_keluar"`
}

type StatistikPendidikanRequest struct {
	Tahun        int  `json:"tahun" binding:"required"`
	TidakSekolah uint `json:"tidak_sekolah"`
	SD           uint `json:"sd"`
	SMP          uint `json:"smp"`
	SMA          uint `json:"sma"`
	Diploma      uint `json:"diploma"`
	Sarjana      uint `json:"sarjana"`
	Pascasarjana uint `json:"pascasarjana"`
}

type StatistikPekerjaanRequest struct {
	Tahun      int  `json:"tahun" binding:"required"`
	Petani     uint `json:"petani"`
	Pedagang   uint `json:"pedagang"`
	PNS        uint `json:"pns"`
	Swasta     uint `json:"swasta"`
	Wiraswasta uint `json:"wiraswasta"`
	Buruh      uint `json:"buruh"`
	Lainnya    uint `json:"lainnya"`
}
