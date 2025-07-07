package main

import (
	"log"
	"os"
	"panelium/backend/internal/config"
	"panelium/backend/internal/db"
	"panelium/backend/internal/global"
	"panelium/backend/internal/handler"
	"panelium/backend/internal/security"
	"panelium/common/id"
	"panelium/common/jwt"
	"time"

	"panelium/backend/internal/model"
)

func main() {
	log.SetOutput(os.Stdout)

	err := global.InitValidator()
	if err != nil {
		log.Printf("Failed to initialize validator: %v", err)
		return
	}

	err = config.Init()
	if err != nil {
		log.Printf("Failed to initialize configuration: %v", err)
		return
	}

	err = global.InitEncryption()
	if err != nil {
		log.Printf("Failed to initialize encryption: %v", err)
		return
	}

	err = db.Init()
	if err != nil {
		log.Printf("Failed to initialize database: %v", err)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "--test-idgen" {
		idGen()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "--test-password-hash" {
		passwordHashTest()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "--test-jwt" {
		jwtTest()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "--make-admin" {
		if len(os.Args) < 3 {
			log.Println("Usage: backend --make-admin <username or email>")
			return
		}
		makeAdmin(os.Args[2])
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	go func() {
		err = handler.Handle("0.0.0.0:" + port)
		if err != nil {
			log.Printf("Failed to start handler: %v", err)
			return
		}
	}()

	// cron job to delete expired sessions
	go func() {
		ticker := time.NewTicker(time.Hour)
		defer ticker.Stop()
		for {
			err := db.DeleteExpiredSessions()
			if err != nil {
				log.Printf("Failed to delete expired sessions: %v", err)
			}
			<-ticker.C
		}
	}()

	log.Printf("Panelium Backend started on port %s", port)

	select {}
}

func idGen() {
	s, err := id.New()
	if err != nil {
		log.Printf("Failed to generate ID: %v", err)
	}
	log.Printf("Generated ID: %s\n", s)
}

func passwordHashTest() {
	pass := "test1234"
	hashed, salt, err := security.HashPassword(pass)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return
	}
	log.Printf("Hashed password: %s, Salt: %s\n", hashed, salt)
	verified := security.VerifyPassword(pass, salt, hashed)
	if verified {
		log.Printf("Password verification successful")
	} else {
		log.Printf("Password verification failed")
	}
}

func jwtTest() {
	sessionId, err := id.New()
	if err != nil {
		log.Printf("Failed to generate session ID: %v", err)
		return
	}

	jti, err := id.New()
	if err != nil {
		log.Printf("Failed to generate JTI: %v", err)
		return
	}

	claims := &jwt.Claims{
		IssuedAt:   time.Now().Unix(),
		Expiration: time.Now().Add(time.Hour).Unix(),
		Audience:   &sessionId,
		Issuer:     jwt.BackendIssuer,
		TokenType:  jwt.MFATokenType,
		JTI:        jti,
	}
	token, err := jwt.CreateJWT(*claims, config.JWTPrivateKeyInstance)
	if err != nil {
		log.Printf("Failed to create JWT: %v", err)
		return
	}
	log.Printf("Generated JWT: %s\n", token)
}

func makeAdmin(arg string) {
	var user model.User
	tx := db.Instance().Where("username = ?", arg).Or("email = ?", arg).First(&user)
	if tx.Error != nil || tx.RowsAffected == 0 {
		log.Printf("User '%s' not found.", arg)
		return
	}
	if user.Admin {
		log.Printf("User '%s' is already an admin.", arg)
		return
	}
	user.Admin = true
	if err := db.Instance().Save(&user).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
		return
	}
	log.Printf("User '%s' promoted to admin.", arg)
}
