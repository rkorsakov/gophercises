package story

import (
	"encoding/json"
	"os"
)

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

type Story map[string]Chapter

func LoadStory(filename string) (Story, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	var story Story
	err = dec.Decode(&story)
	if err != nil {
		return nil, err
	}
	return story, nil
}
