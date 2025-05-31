package main

import (
	"connectrpc.com/connect"
	connectcors "connectrpc.com/cors"
	"context"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	"panelium/daemon/src/global"
	proto_gen_go "panelium/proto-gen-go"
	"panelium/proto-gen-go/proto_gen_goconnect"
)

func withCORS(h http.Handler) http.Handler {
	middleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Change to specific origins in production
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return middleware.Handler(h)
}

type ServerServiceHandler struct {
	proto_gen_goconnect.UnimplementedServerServiceHandler
}

func (s *ServerServiceHandler) RunCommand(
	ctx context.Context,
	req *connect.Request[proto_gen_go.SimpleMessage],
) (*connect.Response[proto_gen_go.Empty], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&proto_gen_go.Empty{})
	return res, nil
}

func (s *ServerServiceHandler) Console(
	ctx context.Context,
	stream *connect.BidiStream[proto_gen_go.SimpleMessage, proto_gen_go.SimpleMessage],
) error {
	log.Println("Request headers: ", stream.RequestHeader())
	for {
		msg, err := stream.Receive()
		if err != nil {
			return err
		}

		log.Printf("received: %s\n", *msg.Text)

		response := &proto_gen_go.SimpleMessage{}
		responseText := "echo: " + *msg.Text
		response.Text = &responseText
		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func main() {
	err := global.Init()
	if err != nil {
		return
	}

	serverServiceHandler := &ServerServiceHandler{}
	mux := http.NewServeMux()
	path, handler := proto_gen_goconnect.NewServerServiceHandler(serverServiceHandler)
	corsHandler := withCORS(handler)
	mux.Handle(path, corsHandler)
	http.ListenAndServe(
		"localhost:9090",
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
