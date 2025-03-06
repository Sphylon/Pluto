package main

import "fmt"

func main() {
	c, _ := LoadConfig(".", true)
	c.Auth()
	fmt.Println(c.Save2File("."))
}
