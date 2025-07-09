package server

import (
	"bufio"
	"connectrpc.com/connect"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"io"
	"panelium/daemon/internal/docker"
	"panelium/proto_gen_go"
	"strings"
	"sync"
	"time"
)

func Console(
	sid string,
	stm *connect.ServerStream[proto_gen_go.SimpleMessage],
) error {
	c, err := console(sid)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to attach to container console"))
	}

	defer func(c *types.HijackedResponse) {
		if c == nil {
			return
		}

		c.Close()
	}(c)

	logs, err := ConsoleLogs(sid)
	if err == nil {
		for _, log := range logs {
			response := &proto_gen_go.SimpleMessage{Text: log}
			if err := stm.Send(response); err != nil {
				return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to send console log"))
			}
		}
	}

	err = attach(c, stm)
	if err != nil {
		return err
	}

	return nil
}

func Terminal(
	sid string,
	stm *connect.ServerStream[proto_gen_go.SimpleMessage],
) error {
	t, err := terminal(sid)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create container terminal"))
	}

	defer func(t *types.HijackedResponse) {
		if t == nil {
			return
		}

		_, _ = t.Conn.Write([]byte("\x03"))
		t.Close()
		terminals.Delete(sid)
	}(t)

	err = attach(t, stm)
	if err != nil {
		return err
	}

	return nil
}

func attach(
	attachStream *types.HijackedResponse,
	stm *connect.ServerStream[proto_gen_go.SimpleMessage],
) error {
	attachStmCh := make(chan string)
	errCh := make(chan error, 2)

	go func() {
		scanner := bufio.NewScanner(attachStream.Reader)
		for scanner.Scan() {
			line := scanner.Text()
			lineWithTimestamp := fmt.Sprintf("[%s] %s", time.Now().Format(time.TimeOnly), line) // TODO: this might be mildly inaccurate
			attachStmCh <- lineWithTimestamp
		}
		if err := scanner.Err(); err != nil {
			errCh <- connect.NewError(connect.CodeInternal, err)
			return
		}
		close(attachStmCh)
	}()

	for {
		select {
		case data, ok := <-attachStmCh:
			if !ok {
				return nil // Channel closed, exit the loop
			}
			response := &proto_gen_go.SimpleMessage{Text: data}
			if err := stm.Send(response); err != nil {
				return err
			}
		case err := <-errCh:
			return err
		}
	}
}

func ConsoleLogs(sid string) ([]string, error) {
	rc, err := docker.Instance().ContainerLogs(context.Background(), fmt.Sprint("server_", sid), container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Tail:       "100",
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get container logs"))
	}
	defer func(rc io.ReadCloser) {
		_ = rc.Close()
	}(rc)

	scanner := bufio.NewScanner(rc)
	var logs []string
	for scanner.Scan() {
		lineRaw := scanner.Text()

		spaceIndex := strings.Index(lineRaw, " ")

		timestamp := lineRaw[:spaceIndex] // everything before the first space
		timestampTime, err := time.Parse(time.RFC3339Nano, timestamp)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to parse timestamp"))
		}

		logText := lineRaw[spaceIndex+1:] // everything after the first space

		formattedTimestamp := timestampTime.Format(time.TimeOnly)
		formattedLine := fmt.Sprintf("[%s] %s", formattedTimestamp, logText)
		logs = append(logs, formattedLine)
	}
	if err := scanner.Err(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to read container logs"))
	}

	return logs, nil
}

func ConsoleCommand(sid string, command string) error {
	if command == "" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("command cannot be empty"))
	}

	c, err := console(sid)
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to attach to container console"))
	}

	_, err = c.Conn.Write([]byte(command + "\n"))
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to write command to container console"))
	}

	return nil
}

func console(sid string) (*types.HijackedResponse, error) {
	c, err := docker.Instance().ContainerAttach(context.Background(), fmt.Sprint("server_", sid), container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		Logs:   false,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to attach to container console"))
	}
	return &c, nil
}

func TerminalCommand(sid string, command string) error {
	if command == "" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("command cannot be empty"))
	}

	var t *types.HijackedResponse

	tm, ok := terminals.Load(sid)
	if !ok {
		var err error
		t, err = terminal(sid)
		if err != nil {
			return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create container terminal"))
		}
	} else {
		var ok bool
		t, ok = tm.(*types.HijackedResponse)
		if !ok {
			return connect.NewError(connect.CodeInternal, errors.New("invalid terminal instance"))
		}
	}

	if t == nil {
		return connect.NewError(connect.CodeInternal, errors.New("terminal not found"))
	}

	_, err := t.Conn.Write([]byte(command + "\n"))
	if err != nil {
		return connect.NewError(connect.CodeInternal, fmt.Errorf("failed to write command to container terminal"))
	}

	return nil
}

var terminals sync.Map

func terminal(sid string) (*types.HijackedResponse, error) {
	if val, ok := terminals.Load(sid); ok {
		if t, ok := val.(*types.HijackedResponse); ok {
			return t, nil
		}
		return nil, connect.NewError(connect.CodeInternal, errors.New("invalid terminal instance"))
	}

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

	err = docker.Instance().ContainerExecStart(context.Background(), eid.ID, container.ExecStartOptions{
		Detach: true,
		Tty:    true,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to start exec instance"))
	}

	t, err := docker.Instance().ContainerExecAttach(context.Background(), eid.ID, container.ExecAttachOptions{
		Tty: true,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to attach to container terminal"))
	}

	terminals.Store(sid, &t)

	return &t, nil
}
