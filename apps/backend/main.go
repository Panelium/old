package main

import (
	"connectrpc.com/connect"
	"context"
	"log"
	"net/http"
	proto_gen_go "panelium/proto-gen-go"
	"panelium/proto-gen-go/proto_gen_goconnect"
)

func main() {
	test := proto_gen_go.TestMessage{}

	text := "Hello, World!"
	test.Text = &text

	number := int32(42)
	test.Number = &number

	array := []string{"item1", "item2", "item3"}
	test.Array = array

	boolean := true
	test.Boolean = &boolean

	log.Println("req: " + test.String())

	client := proto_gen_goconnect.NewTestServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		//connect.WithGRPC(),
		//connect.WithGRPCWeb(),
		//connect.WithProtoJSON(),
	)
	res, err := client.TestMethod(
		context.Background(),
		connect.NewRequest(&test),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("res: " + res.Msg.String())
}
