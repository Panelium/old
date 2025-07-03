package server

import (
	"bufio"
	"connectrpc.com/connect"
	"context"
	"errors"
	"github.com/docker/docker/api/types/container"
	"panelium/daemon/internal/docker"
	"panelium/proto_gen_go"
)

func (s *ServerServiceHandler) Console(
	ctx context.Context,
	stm *connect.BidiStream[proto_gen_go.SimpleMessage, proto_gen_go.SimpleMessage],
) error {
	serverId := ctx.Value("server_id").(string)
	if serverId == "" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("server ID is required"))
	}

	console, err := docker.Instance().ContainerAttach(context.Background(), serverId, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   true, // TODO: remove this and move to different method or just send at start
	})
	if err != nil {
		return connect.NewError(connect.CodeInternal, errors.New("failed to attach to container console"))
	}
	defer console.Close()

	streamCh := make(chan *proto_gen_go.SimpleMessage)
	consoleCh := make(chan []byte)
	errCh := make(chan error, 2)

	// receiving from stream
	go func() {
		for {
			msg, err := stm.Receive()
			if err != nil {
				errCh <- err
				return
			}
			streamCh <- msg
		}
	}()

	// receiving from console
	go func() {
		scanner := bufio.NewScanner(console.Reader)
		for scanner.Scan() {
			consoleCh <- scanner.Bytes()
		}
		if err := scanner.Err(); err != nil {
			errCh <- connect.NewError(connect.CodeInternal, err)
			return
		}
	}()

	for {
		select {
		case msg := <-streamCh:
			// write client command to console
			_, err := console.Conn.Write([]byte(msg.Text + "\n"))
			if err != nil {
				return err
			}
		case data := <-consoleCh:
			// send console output to client
			response := &proto_gen_go.SimpleMessage{Text: string(data)}
			if err := stm.Send(response); err != nil {
				return err
			}
		case err := <-errCh:
			return err
		}
	}
}
