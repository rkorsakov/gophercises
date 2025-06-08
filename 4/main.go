package main

import (
	"fmt"
	"gophercises/4/htmlparser"
)

func main() {
	parsed := htmlparser.ParseHTML("4/data/ex1.html")
	for _, val := range parsed {
		fmt.Print(val)
	}
}
