package main

import (
	"connectrpc.com/connect"
	"context"
	"crypto/tls"
	"golang.org/x/net/http2"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"panelium/common/id"
	proto_gen_go "panelium/proto-gen-go"
	"panelium/proto-gen-go/proto_gen_goconnect"
	"strconv"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "idGen" {
		idGen()
		return
	}

	message := proto_gen_go.SimpleMessage{}
	messageText := "somecommand"
	message.Text = &messageText

	client := proto_gen_goconnect.NewServerServiceClient(
		&http.Client{
			Transport: &http2.Transport{
				AllowHTTP: true,
				DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
					return net.Dial(network, addr)
				},
			},
		},
		"http://localhost:9090",
		//connect.WithGRPC(),
		//connect.WithGRPCWeb(),
		//connect.WithProtoJSON(),
	)
	_, err := client.RunCommand(
		context.Background(),
		connect.NewRequest(&message),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("ok")

	stream := client.Console(
		context.Background(),
	)

	for i := 0; i < 10; i++ {
		message := proto_gen_go.SimpleMessage{}
		messageText := "message " + strconv.Itoa(rand.Intn(100))
		message.Text = &messageText

		if err := stream.Send(&message); err != nil {
			log.Println("Error sending message:", err)
			return
		}

		response, err := stream.Receive()
		if err != nil {
			log.Println("Error receiving response:", err)
			return
		}
		log.Printf("Received response: %s\n", *response.Text)
	}
}

func idGen() {
	s, err := id.New()
	if err != nil {
		log.Fatalf("Failed to generate ID: %v", err)
	}
	log.Printf("Generated ID: %s\n", s)
}
