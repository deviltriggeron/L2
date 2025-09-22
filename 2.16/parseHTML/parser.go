package parsehtml

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"main/entity"
	"main/loader"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func extractLinks(baseUrl *url.URL, r io.Reader) []string {
	var links []string

	doc, err := html.Parse(r)
	if err != nil {
		log.Println("parse html:", err)
		return links
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if u, err := baseUrl.Parse(a.Val); err == nil {
						links = append(links, u.String())
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return links
}

func extractResources(baseUrl *url.URL, r io.Reader) []string {
	var res []string

	doc, err := html.Parse(r)
	if err != nil {
		log.Println("parse html:", err)
		return res
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "img", "script":
				for _, a := range n.Attr {
					if a.Key == "src" {
						if u, err := baseUrl.Parse(a.Val); err == nil {
							res = append(res, u.String())
						}
					}
				}
			case "link":
				var href, rel string
				for _, a := range n.Attr {
					if a.Key == "href" {
						href = a.Val
					}
					if a.Key == "rel" {
						rel = a.Val
					}
				}
				if rel == "stylesheet" && href != "" {
					if u, err := baseUrl.Parse(href); err == nil {
						res = append(res, u.String())
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return res
}

func ParseHTML(flag entity.Flags, dir string, pagesVisited map[string]bool, resVisited map[string]bool) {

	if pagesVisited[*flag.Url] {
		return
	}
	pagesVisited[*flag.Url] = true

	netClient := setConfig()
	u, err := url.Parse(*flag.Url)
	if err != nil {
		log.Println("bad url:", *flag.Url, err)
		return
	}

	response, err := netClient.Get(*flag.Url)
	if err != nil {
		fmt.Println("error get response:", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error read response:", err)
		os.Exit(1)
	}

	var file *os.File
	file, err = createHTML(file, dir, *flag.Url)

	if err != nil {
		fmt.Println("Unable to create file:", err)
	}

	file.WriteString(string(content))
	file.Close()

	links := extractLinks(u, bytes.NewReader(content))
	if *flag.Page {
		res := extractResources(u, bytes.NewReader(content))
		for _, r := range res {
			if !resVisited[r] {
				resVisited[r] = true
				_, err := loader.Downloads(r, dir)
				if err != nil {
					log.Println("resource error:", r, err)
				}
			}
		}
	}

	if *flag.Level == 0 {
		return
	}

	for _, l := range links {
		lu, err := url.Parse(l)
		if err == nil && lu.Host == u.Host {
			newL := *flag.Level - 1
			newF := &entity.Flags{Url: &l, Level: &newL, Page: flag.Page}
			ParseHTML(*newF, dir, pagesVisited, resVisited)
		}
	}
}

func setConfig() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify:          true,
			PreferServerCipherSuites:    true,
			SessionTicketsDisabled:      true,
			DynamicRecordSizingDisabled: true,
		},
	}

	netClient := &http.Client{
		Timeout:   20 * time.Second,
		Transport: tr,
	}

	return netClient
}

func createHTML(f *os.File, dir string, baseUrl string) (*os.File, error) {
	name := ""
	if strings.Contains(baseUrl, "https://") {
		name = strings.ReplaceAll(baseUrl, "https://", "")
	}
	if strings.Contains(baseUrl, "http://") {
		name = strings.ReplaceAll(baseUrl, "https://", "")
	}

	name = strings.ReplaceAll(name, "/", "_")

	f, err := os.Create(dir + "/" + name + ".html")
	if err != nil {
		return f, err
	}

	return f, nil
}
