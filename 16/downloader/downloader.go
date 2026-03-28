package downloader

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func Donwload(rawURL string) error {
	c := http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Proxy: nil,
		},
	}

	resp, err := c.Get(rawURL)
	if err != nil {
		return fmt.Errorf("ошибка сети, url: %v: error: %w", rawURL, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %s", resp.Status)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("failed URL parse %v", rawURL)
	}
	hostName := parsedURL.Hostname()

	if err := os.MkdirAll(hostName, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(hostName, "index.html")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	print("ok")
	return nil

}
