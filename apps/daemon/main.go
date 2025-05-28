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

type DaemonServer struct{}

func (s *DaemonServer) TestMethod(
	ctx context.Context,
	req *connect.Request[proto_gen_go.TestMessage],
) (*connect.Response[proto_gen_go.TestMessage], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(req.Msg)
	return res, nil
}

func main() {
	server := &DaemonServer{}
	mux := http.NewServeMux()
	path, handler := proto_gen_goconnect.NewTestServiceHandler(server)
	corsHandler := withCORS(handler)
	mux.Handle(path, corsHandler)
	http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
