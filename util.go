package main

import (
	"encoding/json"
)

const TimeFormat = "2006-01-02 15:04:05" // Golang의 특별한 시간 형식

type Body map[string]interface{}
type Result Body

func (b Body) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

type Param map[string]string

type Header Param

func (p Param) Marshal() ([]byte, error) {
	return json.Marshal(p)
}
