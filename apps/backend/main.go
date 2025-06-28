package main

import (
	"fmt"
	"os"
	"panelium/backend/internal/config"
	"panelium/backend/internal/db"
	"panelium/backend/internal/global"
	"panelium/backend/internal/handler"
	"panelium/backend/internal/security"
	"panelium/common/id"
)

func main() {
	err := global.InitValidator()
	if err != nil {
		fmt.Printf("Failed to initialize validator: %v", err)
		return
	}

	err = config.Init()
	if err != nil {
		fmt.Printf("Failed to initialize configuration: %v", err)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "idGen" {
		idGen()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "passwordHashTest" {
		passwordHashTest()
		return
	}

	err = db.Init()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v", err)
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	go func() {
		err = handler.Handle("0.0.0.0:" + port)
		if err != nil {
			fmt.Printf("Failed to start handler: %v", err)
			return
		}
	}()

	fmt.Printf("Panelium Backend started on port %s", port)

	select {}
}

func idGen() {
	s, err := id.New()
	if err != nil {
		fmt.Printf("Failed to generate ID: %v", err)
	}
	fmt.Printf("Generated ID: %s\n", s)
}

func passwordHashTest() {
	pass := "test1234"
	hashed, salt, err := security.HashPassword(pass)
	if err != nil {
		fmt.Printf("Failed to hash password: %v", err)
		return
	}
	fmt.Printf("Hashed password: %s, Salt: %s\n", hashed, salt)
	verified := security.VerifyPassword(pass, salt, hashed)
	if verified {
		fmt.Printf("Password verification successful")
	} else {
		fmt.Printf("Password verification failed")
	}
}
