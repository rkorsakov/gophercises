package handler

import (
	"gophercises/3/story"
	"html/template"
	"net/http"
	"strings"
)

type Handler struct {
	Tmpl  *template.Template
	Story story.Story
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	arc := strings.TrimPrefix(r.URL.Path, "/")
	if arc == "" {
		arc = "intro"
	}
	chapter, ok := h.Story[arc]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
	err := h.Tmpl.Execute(w, chapter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
