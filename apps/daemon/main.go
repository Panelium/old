package main

import (
	"log"
	"os"
	"panelium/daemon/internal/config"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/docker"
	"panelium/daemon/internal/handler"
)

func main() {
	log.SetOutput(os.Stdout)

	err := config.Init()
	if err != nil {
		log.Printf("Failed to initialize configuration: %v", err)
		return
	}

	err = db.Init()
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return
	}

	err = docker.Init()
	if err != nil {
		log.Printf("Failed to initialize Docker client: %v", err)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	go func() {
		err = handler.Handle("0.0.0.0:" + port)
		if err != nil {
			log.Printf("Failed to start handler: %v", err)
			return
		}
	}()

	log.Printf("Panelium Daemon started on port %s", port)

	select {}
}
