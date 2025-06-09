package main

import (
	"bufio"
	"fmt"
	"gophercises/4/htmlparser"
	"log"
	"os"
)

func main() {
	file, err := os.Open("4/data/ex1.html")
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	parsed := htmlparser.ParseHTML(reader)
	for _, val := range parsed {
		fmt.Print(val)
	}
}
