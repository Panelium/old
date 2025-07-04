package server

import (
	"bufio"
	"connectrpc.com/connect"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"net"
	"panelium/daemon/internal/docker"
	"panelium/proto_gen_go"
)

func Console(sid string, stm *connect.BidiStream[proto_gen_go.SimpleMessage, proto_gen_go.SimpleMessage]) error {
	c, err := console(sid)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to attach to container console"))
	}

	//send console logs

	err = attach(c, stm)
	if err != nil {
		return err
	}

	return nil
}

func Terminal(sid string, stm *connect.BidiStream[proto_gen_go.SimpleMessage, proto_gen_go.SimpleMessage]) error {
	t, err := terminal(sid)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create container terminal"))
	}

	defer func(Conn net.Conn, b []byte) {
		_, _ = Conn.Write(b)
	}(t.Conn, []byte("\x03")) // send Ctrl+C to terminate the terminal session

	err = attach(t, stm)
	if err != nil {
		return err
	}

	return nil
}

func attach(attachStream *types.HijackedResponse, stm *connect.BidiStream[proto_gen_go.SimpleMessage, proto_gen_go.SimpleMessage]) error {
	defer attachStream.Close()

	resStmCh := make(chan *proto_gen_go.SimpleMessage)
	attachStmCh := make(chan []byte)
	errCh := make(chan error, 2)

	// receiving from stream
	go func() {
		for {
			msg, err := stm.Receive()
			if err != nil {
				errCh <- err
				return
			}
			resStmCh <- msg
		}
	}()

	// receiving from console
	go func() {
		scanner := bufio.NewScanner(attachStream.Reader)
		for scanner.Scan() {
			attachStmCh <- scanner.Bytes()
		}
		if err := scanner.Err(); err != nil {
			errCh <- connect.NewError(connect.CodeInternal, err)
			return
		}
	}()

	for {
		select {
		case msg := <-resStmCh:
			// write client command to console
			_, err := attachStream.Conn.Write([]byte(msg.Text + "\n"))
			if err != nil {
				return err
			}
		case data := <-attachStmCh:
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

func ConsoleLogs(sid string) ([]byte, error) {
	return nil, nil
}

func console(sid string) (*types.HijackedResponse, error) {
	console, err := docker.Instance().ContainerAttach(context.Background(), sid, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   false,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to attach to container console"))
	}
	return &console, nil
}

func terminal(sid string) (*types.HijackedResponse, error) {
	eid, err := docker.Instance().ContainerExecCreate(context.Background(), sid, container.ExecOptions{
		User:         "root",
		Privileged:   true,
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Detach:       true,
		Cmd:          []string{"sh"}, // TODO: sh might not always be available
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create exec instance"))
	}

	terminal, err := docker.Instance().ContainerExecAttach(context.Background(), eid.ID, container.ExecAttachOptions{
		Tty: true,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to attach to container terminal"))
	}
	return &terminal, nil
}
