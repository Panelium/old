package main

import (
	"fmt"
	"os"
	"panelium/backend/internal/config"
	"panelium/backend/internal/db"
	"panelium/backend/internal/handler"
	"panelium/backend/internal/security"
	"panelium/common/id"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "idGen" {
		idGen()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "passwordHashTest" {
		passwordHashTest()
		return
	}

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
	pepper := security.GenerateRandomString()
	hashed, salt := security.HashPassword(pass, pepper)
	fmt.Printf("Hashed password: %s, Salt: %s\n", hashed, salt)
	verified := security.VerifyPassword(pass, salt, pepper, hashed)
	if verified {
		fmt.Printf("Password verification successful")
	} else {
		fmt.Printf("Password verification failed")
	}
}
