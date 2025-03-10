package main

import "encoding/json"

type Body map[string]interface{}
type Result Body

func (j Body) Marshal() ([]byte, error) {
	return json.Marshal(j)
}

type Param map[string]string

type Header Param

func (p Param) Marshal() ([]byte, error) {
	return json.Marshal(p)
}
