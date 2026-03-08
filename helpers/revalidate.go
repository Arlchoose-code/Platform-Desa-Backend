// Platform Desa — Revalidate Helper
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Arlchoose-code/platform-desa-backend/config"
)

var revalidateClient = &http.Client{Timeout: 5 * time.Second}

type revalidatePayload struct {
	Secret string `json:"secret"`
	Path   string `json:"path"`
}

func RevalidatePath(path string) {
	go triggerRevalidate(path)
}

func RevalidatePaths(paths ...string) {
	for _, path := range paths {
		go triggerRevalidate(path)
	}
}

func triggerRevalidate(path string) {
	if config.App.NextJS.URL == "" || config.App.NextJS.RevalidateSecret == "" {
		return
	}

	payload, err := json.Marshal(revalidatePayload{
		Secret: config.App.NextJS.RevalidateSecret,
		Path:   path,
	})
	if err != nil {
		fmt.Printf("[revalidate] gagal marshal payload: %v\n", err)
		return
	}

	maxRetry := 5
	for attempt := 1; attempt <= maxRetry; attempt++ {
		err := doRevalidate(payload)
		if err == nil {
			if attempt > 1 {
				fmt.Printf("[revalidate] berhasil setelah %d percobaan: %s\n", attempt, path)
			}
			return
		}

		if attempt == maxRetry {
			fmt.Printf("[revalidate] menyerah setelah %d percobaan untuk path: %s — %v\n",
				maxRetry, path, err)
			return
		}

		wait := time.Duration(1<<attempt) * time.Second
		fmt.Printf("[revalidate] percobaan %d gagal untuk %s: %v — retry dalam %s\n",
			attempt, path, err, wait)
		time.Sleep(wait)
	}
}

func doRevalidate(payload []byte) error {
	url := config.App.NextJS.URL + "/api/revalidate"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("gagal buat request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := revalidateClient.Do(req)
	if err != nil {
		return fmt.Errorf("gagal kirim request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Next.js response %d", resp.StatusCode)
	}

	return nil
}
