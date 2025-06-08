package main

import (
	"fmt"
	"gophercises/3/handler"
	"gophercises/3/story"
	"html/template"
	"net/http"
)

func main() {
	webStory, err := story.LoadStory("3/gopher.json")
	if err != nil {
		fmt.Println(err)
	}
	tpl, err := template.ParseFiles("3/templates/story.html")
	if err != nil {
		fmt.Println(err)
	}
	h := handler.Handler{Tmpl: tpl, Story: webStory}
	fmt.Println("Starting server on port :8080")
	err = http.ListenAndServe(":8080", h)
	if err != nil {
		fmt.Println(err)
	}
}
