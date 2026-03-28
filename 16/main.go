package main

import "wget/downloader"

func main() {
	if err := downloader.Donwload("https://example.com/"); err != nil {
		panic(err)
	}
}
