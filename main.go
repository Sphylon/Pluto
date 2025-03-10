package main

import "fmt"

func main() {
	kisd, err := NewKISD(false)
	if err != nil {
		fmt.Println(err)
	}

	kisd.Auth()
	kisd.GetPrice("012450")
	kisd.Order("012450", true)
}
