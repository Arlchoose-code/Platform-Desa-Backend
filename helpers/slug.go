// Platform Desa — Slug Helper
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package helpers

import (
	"fmt"
	"regexp"
	"strings"
	"time"
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
	return fmt.Sprintf("%s-%d", base, time.Now().Unix())
}
