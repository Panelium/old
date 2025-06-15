package main

import (
	"log"
	"os"
	"panelium/backend/global"
	"panelium/backend/handler"
	"panelium/common/id"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "idGen" {
		idGen()
		return
	}

	err := global.Init()
	if err != nil {
		log.Fatalf("Failed to initialize global settings: %v", err)
		return
	}

	err = handler.Handle("localhost:9090")
	if err != nil {
		log.Fatalf("Failed to start handler: %v", err)
		return
	}
}

func idGen() {
	s, err := id.New()
	if err != nil {
		log.Fatalf("Failed to generate ID: %v", err)
	}
	log.Printf("Generated ID: %s\n", s)
}
