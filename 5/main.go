package main

import (
	"flag"
	"gophercises/4/htmlparser"
	"log"
	"net/http"
)

func main() {
	link := flag.String("l", "https://calhoun.io", "Link to build Sitemap")
	flag.Parse()
	res, err := http.Get(*link)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	htmlparser.ParseHTML(res.Body)
}
