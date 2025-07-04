package main

import (
	"fmt"
	"os"
	"panelium/daemon/internal/config"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/docker"
	"panelium/daemon/internal/handler"
)

func main() {
	err := config.Init()
	if err != nil {
		fmt.Printf("Failed to initialize configuration: %v", err)
		return
	}

	err = db.Init()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v", err)
		return
	}

	err = docker.Init()
	if err != nil {
		fmt.Printf("Failed to initialize Docker client: %v", err)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	go func() {
		err = handler.Handle("0.0.0.0:" + port)
		if err != nil {
			fmt.Printf("Failed to start handler: %v", err)
			return
		}
	}()

	fmt.Printf("Panelium Daemon started on port %s", port)

	select {}
}
