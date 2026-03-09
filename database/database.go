// Platform Desa — Database
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package database

import (
	"fmt"
	"log"

	"github.com/Arlchoose-code/platform-desa-backend/config"
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.App.Database.User,
		config.App.Database.Password,
		config.App.Database.Host,
		config.App.Database.Port,
		config.App.Database.Name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}
	fmt.Println("Database connected successfully!")

	err = DB.AutoMigrate(
		// System
		&models.Setting{},
		&models.User{},
		&models.ActivityLog{},
		&models.Media{},

		// Profil
		&models.ProfilDesa{},
		&models.PotensiDesa{},

		// Pemerintahan
		&models.Jabatan{},
		&models.Pejabat{},
		&models.LembagaDesa{},

		// Informasi Publik
		&models.KategoriBerita{},
		&models.Berita{},
		&models.Pengumuman{},
		&models.Agenda{},

		// Kependudukan
		&models.StatistikPenduduk{},
		&models.StatistikPendidikan{},
		&models.StatistikPekerjaan{},

		// Wisata
		&models.KategoriWisata{},
		&models.Wisata{},
		&models.WisataGaleri{},

		// Galeri
		&models.KategoriGaleri{},
		&models.Galeri{},

		// UMKM
		&models.KategoriUMKM{},
		&models.UMKM{},
		&models.ProdukUMKM{},

		// Persuratan
		&models.JenisSurat{},
		&models.PengajuanSurat{},
		&models.RiwayatSurat{},

		// APBDes
		&models.APBDes{},
		&models.APBDesDetail{},

		// Pengaduan
		&models.Pengaduan{},

		// Regulasi
		&models.KategoriRegulasi{},
		&models.Regulasi{},

		// Sistem
		&models.PetaFasilitas{},
		&models.FAQ{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}
	fmt.Println("Database migrated successfully!")

	SeedSettings(DB)
	fmt.Println("Settings seeded successfully!")

	SeedUsers(DB)
	fmt.Println("Users seeded successfully!")
}
