package main

import (
	"fmt"
	"log"
	//"time"
)

func main() {

	forth, err := NewSerialForth()

	if err != nil {

		log.Fatalf("Error opening arduino: %v", err)
	}

	defer forth.Close()

	fmt.Println(forth.Query("6 IN 1 = ."))
}
