package main

import (
	"flag"
	"log"

	entity "main/entity"

	"net/url"
	"os"

	parser "main/parseHTML"
)

func parseFlags() *entity.Flags {
	m := flag.String("m", "", "URL of the download site (for example: https://example.com)")
	l := flag.Int("l", 0, "Recursion depth for downloading pages (0 - only the specified page)")
	p := flag.Bool("p", false, "Download all page resources (images, CSS, JS, etc.)")

	flag.Parse()

	return &entity.Flags{
		Url:   m,
		Level: l,
		Page:  p,
	}
}

func createDir(baseUrl string) string {
	os.Mkdir("downloads", 0755)
	u, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatalf("invalid start url: %v", err)
	}
	dirName := "downloads" + "/" + u.Host
	os.Mkdir(dirName, 0755)

	return dirName
}

func main() {
	f := parseFlags()
	dir := createDir(*f.Url)
	pagesVisited := make(map[string]bool)
	resVisited := make(map[string]bool)
	parser.ParseHTML(*f, dir, pagesVisited, resVisited)
}
