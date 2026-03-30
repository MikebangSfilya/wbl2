package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"wget/downloader"
	"wget/parser"
)

var visited = make(map[string]bool)

func Crawl(targetURL string, currentDepth, maxDepth int, baseHost string, savePath string) {

	if currentDepth > maxDepth {
		return
	}
	if visited[targetURL] {
		return
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil || parsedURL.Hostname() != baseHost {
		return
	}

	visited[targetURL] = true
	fmt.Printf("=> Качаем [%d/%d]: %s\n", currentDepth, maxDepth, targetURL)

	body, err := downloader.Download(targetURL)
	if err != nil {
		slog.Error("Ошибка скачивания", "url", targetURL, "err", err)
		return
	}

	if err := downloader.Save(body, savePath); err != nil {
		slog.Error("Ошибка сохранения", "path", savePath, "err", err)
		return
	}

	links, err := parser.ExtractLinks(body, parsedURL)
	if err != nil {
		slog.Error("Ошибка парсинга", "url", targetURL, "err", err)
		return
	}

	for _, link := range links {
		refURL, err := url.Parse(link)
		if err != nil {
			continue
		}

		nextPath := filepath.Join(baseHost, refURL.Path)

		if strings.HasSuffix(nextPath, string(filepath.Separator)) || filepath.Ext(nextPath) == "" {
			nextPath = filepath.Join(nextPath, "index.html")
		}

		Crawl(link, currentDepth+1, maxDepth, baseHost, nextPath)
	}
}

func main() {
	depth := flag.Int("d", 1, "recursion depth")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: mywget [-d <depth>] <url>")
		os.Exit(1)
	}

	startURL := args[0]
	parsedURL, err := url.Parse(startURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("Старт загрузки:", startURL)

	baseHost := parsedURL.Hostname()
	mainFilePath := filepath.Join(baseHost, "index.html")
	Crawl(startURL, 0, *depth, baseHost, mainFilePath)

	fmt.Println("Готово! Все ресурсы скачаны.")
}
