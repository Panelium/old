package main

import (
	"log"
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
		log.Fatalf("Failed to initialize configuration: %v", err)
		return
	}

	err = db.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
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

func passwordHashTest() {
	pass := "test1234"
	pepper := security.GenerateRandomString()
	hashed, salt := security.HashPassword(pass, pepper)
	log.Printf("Hashed password: %s, Salt: %s\n", hashed, salt)
	verified := security.VerifyPassword(pass, salt, pepper, hashed)
	if verified {
		log.Println("Password verification successful")
	} else {
		log.Println("Password verification failed")
	}
}
