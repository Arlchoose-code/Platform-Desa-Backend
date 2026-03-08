// Platform Desa — Slug Helper
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Arlchoose-code/platform-desa-backend/database"
)

var (
	reNonAlnum = regexp.MustCompile(`[^a-z0-9\s-]`)
	reSpaces   = regexp.MustCompile(`[\s-]+`)
)

func GenerateSlug(text string) string {
	s := strings.ToLower(text)
	s = strings.NewReplacer(
		"ä", "a", "ö", "o", "ü", "u",
		"à", "a", "á", "a", "â", "a",
		"è", "e", "é", "e", "ê", "e",
		"ì", "i", "í", "i", "î", "i",
		"ò", "o", "ó", "o", "ô", "o",
		"ù", "u", "ú", "u", "û", "u",
	).Replace(s)
	s = reNonAlnum.ReplaceAllString(s, "")
	s = reSpaces.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func UniqueSlug(base string) string {
	slug := base
	tables := []string{
		"beritas", "kategori_beritas",
		"wisatas", "kategori_wisatas",
		"umkms", "kategori_u_m_k_m_s",
		"galeris", "kategori_galeris",
		"regulasis", "kategori_regulasis",
		"f_a_q_s",
		"pejabats",
		"pengumumans",
		"agendas",
		"potensi_desas",
	}

	for i := 1; ; i++ {
		found := false
		for _, table := range tables {
			var count int64
			database.DB.Table(table).Where("slug = ?", slug).Count(&count)
			if count > 0 {
				found = true
				break
			}
		}
		if !found {
			return slug
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}
}
