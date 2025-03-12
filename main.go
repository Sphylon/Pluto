package main

import (
	"log"
)

func main() {
	kisd, err := NewKISD(false, ".")
	if err != nil {
		log.Panic(err)
	} else {
		log.Printf("Create new KISD object ...")
	}

	if err := kisd.AddStock("329750"); err != nil {
		log.Panic(err)
	} else {
		log.Printf("Add %s stock ..", "329750")
	}

	if err := kisd.OrderCash("329750", 1, 13200, true); err != nil {
		log.Panic(err)
	} else {
		log.Printf("Buy %s stock ..", "329750")
	}
}
