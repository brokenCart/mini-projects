package main

import "encoding/json"

func parseJSON(content []byte) (map[string]StoryArc, error) {
	stories := make(map[string]StoryArc)
	err := json.Unmarshal(content, &stories)
	if err != nil {
		return nil, err
	}
	return stories, nil
}

type StoryArc struct {
	Title           string   `json:"title"`
	StoryParagraphs []string `json:"story"`
	Options         []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}
