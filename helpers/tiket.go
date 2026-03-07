// Platform Desa — Tiket Helper
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package helpers

import (
	"fmt"
	"time"
)

func GenerateNomorPengajuan() string {
	now := time.Now()
	return fmt.Sprintf("SURAT-%s-%06d",
		now.Format("20060102"),
		now.UnixNano()%1000000,
	)
}

func GenerateNomorTiket() string {
	now := time.Now()
	return fmt.Sprintf("ADU-%s-%06d",
		now.Format("20060102"),
		now.UnixNano()%1000000,
	)
}
