// Platform Desa — Settings Seeder
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package database

import (
	"github.com/Arlchoose-code/platform-desa-backend/models"
	"gorm.io/gorm"
)

func SeedSettings(db *gorm.DB) {
	settings := []models.Setting{
		// =====================
		// Group: general
		// =====================
		{Key: "site_name", Value: strPtr("Desa Maju Bersama"), Group: "general"},
		{Key: "site_tagline", Value: strPtr("Membangun Desa, Membangun Bangsa"), Group: "general"},
		{Key: "site_description", Value: strPtr("Website resmi Desa Maju Bersama"), Group: "general"},
		{Key: "site_logo", Value: nil, Group: "general"},
		{Key: "site_favicon", Value: nil, Group: "general"},
		{Key: "site_email", Value: strPtr("desa@example.com"), Group: "general"},
		{Key: "site_phone", Value: strPtr("0812-3456-7890"), Group: "general"},
		{Key: "site_address", Value: strPtr("Jl. Desa No. 1, Kecamatan Maju, Kabupaten Bersama"), Group: "general"},
		{Key: "site_maps_url", Value: strPtr(""), Group: "general"},
		{Key: "site_latitude", Value: strPtr("-6.200000"), Group: "general"},
		{Key: "site_longitude", Value: strPtr("106.816666"), Group: "general"},
		{Key: "site_kode_desa", Value: strPtr(""), Group: "general"},
		{Key: "site_kecamatan", Value: strPtr(""), Group: "general"},
		{Key: "site_kabupaten", Value: strPtr(""), Group: "general"},
		{Key: "site_provinsi", Value: strPtr(""), Group: "general"},
		{Key: "site_kode_pos", Value: strPtr(""), Group: "general"},

		// =====================
		// Group: seo
		// =====================
		{Key: "seo_meta_title", Value: strPtr(""), Group: "seo"},
		{Key: "seo_meta_description", Value: strPtr(""), Group: "seo"},
		{Key: "seo_meta_keywords", Value: strPtr(""), Group: "seo"},
		{Key: "seo_og_image", Value: nil, Group: "seo"},
		{Key: "seo_og_type", Value: strPtr("website"), Group: "seo"},
		{Key: "seo_twitter_card", Value: strPtr("summary_large_image"), Group: "seo"},
		{Key: "seo_google_analytics", Value: strPtr(""), Group: "seo"},
		{Key: "seo_google_site_verification", Value: strPtr(""), Group: "seo"},
		{Key: "seo_robots", Value: strPtr("index, follow"), Group: "seo"},
		{Key: "seo_canonical_url", Value: strPtr(""), Group: "seo"},
		{Key: "seo_sitemap_aktif", Value: strPtr("true"), Group: "seo"},
		{Key: "seo_schema_org", Value: strPtr("true"), Group: "seo"},

		// =====================
		// Group: sosmed
		// =====================
		{Key: "sosmed_facebook", Value: strPtr(""), Group: "sosmed"},
		{Key: "sosmed_instagram", Value: strPtr(""), Group: "sosmed"},
		{Key: "sosmed_youtube", Value: strPtr(""), Group: "sosmed"},
		{Key: "sosmed_twitter", Value: strPtr(""), Group: "sosmed"},
		{Key: "sosmed_tiktok", Value: strPtr(""), Group: "sosmed"},
		{Key: "sosmed_whatsapp", Value: strPtr(""), Group: "sosmed"},

		// =====================
		// Group: hero
		// =====================
		{Key: "hero_judul", Value: strPtr("Selamat Datang di Desa Maju Bersama"), Group: "hero"},
		{Key: "hero_subjudul", Value: strPtr("Membangun desa yang mandiri, sejahtera, dan berbudaya"), Group: "hero"},
		{Key: "hero_foto", Value: nil, Group: "hero"},
		{Key: "hero_video", Value: nil, Group: "hero"},
		{Key: "hero_cta_label", Value: strPtr("Jelajahi Desa"), Group: "hero"},
		{Key: "hero_cta_url", Value: strPtr("/profil"), Group: "hero"},

		// =====================
		// Group: surat
		// =====================
		{Key: "surat_jam_buka", Value: strPtr("08:00"), Group: "surat"},
		{Key: "surat_jam_tutup", Value: strPtr("15:00"), Group: "surat"},
		{Key: "surat_hari_kerja", Value: strPtr("Senin - Jumat"), Group: "surat"},
		{Key: "surat_estimasi_default", Value: strPtr("3"), Group: "surat"},
		{Key: "surat_catatan", Value: strPtr("Harap membawa dokumen persyaratan yang lengkap"), Group: "surat"},
		{Key: "surat_aktif", Value: strPtr("true"), Group: "surat"},

		// =====================
		// Group: pengaduan
		// =====================
		{Key: "pengaduan_aktif", Value: strPtr("true"), Group: "pengaduan"},
		{Key: "pengaduan_anonim_aktif", Value: strPtr("true"), Group: "pengaduan"},
		{Key: "pengaduan_catatan", Value: strPtr("Pengaduan akan diproses dalam 3-7 hari kerja"), Group: "pengaduan"},
		{Key: "pengaduan_max_foto", Value: strPtr("3"), Group: "pengaduan"},

		// =====================
		// Group: footer
		// =====================
		{Key: "footer_tentang", Value: strPtr(""), Group: "footer"},
		{Key: "footer_copyright", Value: strPtr("© 2026 Desa Maju Bersama. All rights reserved."), Group: "footer"},
		{Key: "footer_logo", Value: nil, Group: "footer"},
		{Key: "footer_show_sosmed", Value: strPtr("true"), Group: "footer"},
		{Key: "footer_show_maps", Value: strPtr("true"), Group: "footer"},

		// =====================
		// Group: notifikasi
		// =====================
		{Key: "notif_email_aktif", Value: strPtr("false"), Group: "notifikasi"},
		{Key: "notif_email_host", Value: strPtr(""), Group: "notifikasi"},
		{Key: "notif_email_port", Value: strPtr("587"), Group: "notifikasi"},
		{Key: "notif_email_user", Value: strPtr(""), Group: "notifikasi"},
		{Key: "notif_email_password", Value: strPtr(""), Group: "notifikasi"},
		{Key: "notif_email_from", Value: strPtr(""), Group: "notifikasi"},
		{Key: "notif_email_from_name", Value: strPtr(""), Group: "notifikasi"},
		{Key: "notif_whatsapp_aktif", Value: strPtr("false"), Group: "notifikasi"},
		{Key: "notif_whatsapp_api_url", Value: strPtr(""), Group: "notifikasi"},
		{Key: "notif_whatsapp_token", Value: strPtr(""), Group: "notifikasi"},

		// =====================
		// Group: maintenance
		// =====================
		{Key: "maintenance_aktif", Value: strPtr("false"), Group: "maintenance"},
		{Key: "maintenance_pesan", Value: strPtr("Website sedang dalam pemeliharaan. Silakan coba beberapa saat lagi."), Group: "maintenance"},
		{Key: "maintenance_estimasi", Value: strPtr(""), Group: "maintenance"},
	}

	for _, s := range settings {
		var existing models.Setting
		if err := db.Where("`key` = ?", s.Key).First(&existing).Error; err != nil {
			db.Create(&s)
		}
	}
}

func strPtr(s string) *string {
	return &s
}
