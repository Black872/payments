package main

import (
	"log"
	"payments/cmd"
)

func main() {
	if err := cmd.RunApp(); err != nil {
		log.Printf("exit reason: %v\n", err)
	}
	log.Println("app closed")
}
