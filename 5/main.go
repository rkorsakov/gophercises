package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"gophercises/4/htmlparser"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

type URL struct {
	Loc string `xml:"loc"`
}

func main() {
	root := flag.String("url", "https://gophercises.com", "URL to build sitemap")
	flag.Parse()

	baseURL, err := url.Parse(*root)
	if err != nil {
		log.Fatal("Invalid URL:", err)
	}

	visited := make(map[string]bool)
	queue := []string{*root}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] {
			continue
		}
		visited[current] = true
		resp, err := http.Get(current)
		if err != nil {
			log.Printf("Error fetching %s: %v", current, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("Skipping %s: status %d", current, resp.StatusCode)
			continue
		}
		links := htmlparser.ParseHTML(resp.Body)
		resp.Body.Close()
		for _, link := range links {
			href := normalizeURL(link.Href, baseURL)
			if href != "" && !visited[href] && sameDomain(href, baseURL) {
				queue = append(queue, href)
			}
		}
	}

	var urls []URL
	for u := range visited {
		urls = append(urls, URL{Loc: u})
	}
	urlset := URLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  urls,
	}

	file, err := os.Create("sitemap.xml")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	file.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	if err := encoder.Encode(urlset); err != nil {
		log.Fatal("Error writing XML:", err)
	}

	fmt.Println("Sitemap saved to sitemap.xml")
}

func sameDomain(href string, base *url.URL) bool {
	u, err := url.Parse(href)
	if err != nil {
		return false
	}
	return u.Host == "" || u.Host == base.Host
}

func normalizeURL(href string, base *url.URL) string {
	u, err := url.Parse(href)
	if err != nil {
		return ""
	}
	if u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https" {
		return ""
	}
	if u.Host == "" {
		u = base.ResolveReference(u)
	}
	u.Fragment = ""
	u.Path = strings.TrimRight(u.Path, "/")
	return u.String()
}
