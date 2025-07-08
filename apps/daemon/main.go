package main

import (
	"connectrpc.com/connect"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"panelium/daemon/internal/config"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/docker"
	"panelium/daemon/internal/handler"
	"panelium/daemon/internal/security"
	"panelium/proto_gen_go/backend"
	"panelium/proto_gen_go/backend/backendconnect"
	"time"
)

func main() {
	log.SetOutput(os.Stdout)

	err := config.Init()
	if err != nil {
		log.Printf("Failed to initialize configuration: %v", err)
		return
	}

	if config.SecretsInstance.GetBackendToken() == "" {
		log.Println("Backend token is not set. Please set it in the configuration file.")
		return
	}
	if config.SecretsInstance.GetNodeJTI() == "" {
		err = tryRegisterWithBackend()
		if err != nil {
			log.Printf("Failed to register with backend: %v", err)
			return
		}
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

func tryRegisterWithBackend() error {
	client := backendconnect.NewDaemonServiceClient(http.DefaultClient, config.ConfigInstance.GetBackendHost())

	nodeToken, nodeJTI, _, err := security.CreateNodeToken(time.Now())
	if err != nil {
		return err
	}

	req := connect.NewRequest(&backend.RegisterDaemonRequest{
		NodeToken: nodeToken,
	})
	req.Header().Add("Authorization", config.SecretsInstance.BackendToken)

	res, err := client.RegisterDaemon(context.Background(), req)
	if err != nil {
		return err
	}

	if res == nil || res.Msg.Success != true {
		return fmt.Errorf("something went wrong")
	}

	config.SecretsInstance.SetNodeJTI(nodeJTI)
	err = config.SecretsInstance.Save()
	if err != nil {
		return err
	}

	log.Printf("Successfully registered with backend. Node JTI: %s", nodeJTI)

	return nil
}
