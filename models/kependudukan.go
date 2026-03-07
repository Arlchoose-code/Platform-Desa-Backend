// Platform Desa — Model Kependudukan
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package models

import "time"

type StatistikPenduduk struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	Tahun         int       `gorm:"not null" json:"tahun"`
	Bulan         int       `gorm:"not null" json:"bulan"`
	TotalPenduduk uint      `gorm:"default:0" json:"total_penduduk"`
	LakiLaki      uint      `gorm:"default:0" json:"laki_laki"`
	Perempuan     uint      `gorm:"default:0" json:"perempuan"`
	TotalKK       uint      `gorm:"default:0" json:"total_kk"`
	Kelahiran     uint      `gorm:"default:0" json:"kelahiran"`
	Kematian      uint      `gorm:"default:0" json:"kematian"`
	PindahMasuk   uint      `gorm:"default:0" json:"pindah_masuk"`
	PindahKeluar  uint      `gorm:"default:0" json:"pindah_keluar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type StatistikPendidikan struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Tahun        int       `gorm:"not null;uniqueIndex" json:"tahun"`
	TidakSekolah uint      `gorm:"default:0" json:"tidak_sekolah"`
	SD           uint      `gorm:"default:0" json:"sd"`
	SMP          uint      `gorm:"default:0" json:"smp"`
	SMA          uint      `gorm:"default:0" json:"sma"`
	Diploma      uint      `gorm:"default:0" json:"diploma"`
	Sarjana      uint      `gorm:"default:0" json:"sarjana"`
	Pascasarjana uint      `gorm:"default:0" json:"pascasarjana"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StatistikPekerjaan struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	Tahun      int       `gorm:"not null;uniqueIndex" json:"tahun"`
	Petani     uint      `gorm:"default:0" json:"petani"`
	Pedagang   uint      `gorm:"default:0" json:"pedagang"`
	PNS        uint      `gorm:"default:0" json:"pns"`
	Swasta     uint      `gorm:"default:0" json:"swasta"`
	Wiraswasta uint      `gorm:"default:0" json:"wiraswasta"`
	Buruh      uint      `gorm:"default:0" json:"buruh"`
	Lainnya    uint      `gorm:"default:0" json:"lainnya"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
