package main

import (
	"fmt"
	"panelium/daemon/internal/config"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/docker"
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
}
