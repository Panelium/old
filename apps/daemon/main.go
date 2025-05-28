package main

import (
	"connectrpc.com/connect"
	"context"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
	"net/http"
	proto_gen_go "panelium/proto-gen-go"
	"panelium/proto-gen-go/proto_gen_goconnect"
)

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
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
