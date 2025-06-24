package handler

import (
	"connectrpc.com/authn"
	connectcors "connectrpc.com/cors"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
	"panelium/backend/internal/handler/auth"
	middleware "panelium/backend/internal/middleware"
	"panelium/proto-gen-go/proto_gen_goconnect"
)

func withCORS(h http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Change to specific origins in production
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return corsMiddleware.Handler(h)
}

func Handle(host string) error {
	mux := http.NewServeMux()

	mux.Handle(proto_gen_goconnect.NewAuthServiceHandler(&auth.AuthServiceHandler{}))

	authMiddleware := authn.NewMiddleware(middleware.Authentication)
	authedMux := authMiddleware.Wrap(mux)

	handler := h2c.NewHandler(authedMux, &http2.Server{})
	corsHandler := withCORS(handler)
	err := http.ListenAndServe(
		host,
		corsHandler,
	)
	if err != nil {
		return err
	}
	return nil
}

//message := proto_gen_go.SimpleMessage{
//	Text: "somecommand",
//}
//
//client := proto_gen_goconnect.NewServerServiceClient(
//	&http.Client{
//		Transport: &http2.Transport{
//			AllowHTTP: true,
//			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
//				return net.Dial(network, addr)
//			},
//		},
//	},
//	"http://localhost:9090",
//	//connect.WithGRPC(),
//	//connect.WithGRPCWeb(),
//	//connect.WithProtoJSON(),
//)
//_, err := client.RunCommand(
//	context.Background(),
//	connect.NewRequest(&message),
//)
//if err != nil {
//	log.Println(err)
//	return
//}
//log.Println("ok")
//
//stream := client.Console(
//	context.Background(),
//)
//
//for i := 0; i < 10; i++ {
//	message := proto_gen_go.SimpleMessage{
//		Text: "message " + strconv.Itoa(rand.Intn(100)),
//	}
//
//	if err := stream.Send(&message); err != nil {
//		log.Println("Error sending message:", err)
//		return
//	}
//
//	response, err := stream.Receive()
//	if err != nil {
//		log.Println("Error receiving response:", err)
//		return
//	}
//	log.Printf("Received response: %s\n", response.Text)
//}

//type ServerServiceHandler struct {
//	proto_gen_goconnect.UnimplementedServerServiceHandler
//}
//
//func (s *ServerServiceHandler) RunCommand(
//	ctx context.Context,
//	req *connect.Request[proto_gen_go.SimpleMessage],
//) (*connect.Response[proto_gen_go.Empty], error) {
//	log.Println("Request headers: ", req.Header())
//	res := connect.NewResponse(&proto_gen_go.Empty{})
//	return res, nil
//}
//
//func (s *ServerServiceHandler) Console(
//	ctx context.Context,
//	stream *connect.BidiStream[proto_gen_go.SimpleMessage, proto_gen_go.SimpleMessage],
//) error {
//	log.Println("Request headers: ", stream.RequestHeader())
//	for {
//		msg, err := stream.Receive()
//		if err != nil {
//			return err
//		}
//
//		log.Printf("received: %s\n", msg.Text)
//
//		response := &proto_gen_go.SimpleMessage{
//			Text: "echo: " + msg.Text,
//		}
//		if err := stream.Send(response); err != nil {
//			return err
//		}
//	}
//}
