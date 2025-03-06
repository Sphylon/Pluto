package main

import "encoding/json"

type Body map[string]interface{}

func (b Body) Add(key, content string) {
	b[key] = content
}

func (b Body) Marshal() ([]byte, error) {
	return json.Marshal(b)
}
