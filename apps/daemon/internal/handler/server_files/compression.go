package server_files

import (
	"connectrpc.com/connect"
	"context"
	"os"
	"panelium/daemon/internal/security"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go/daemon"
)

func (s *ServerFilesServiceHandler) CompressFile(ctx context.Context, req *connect.Request[daemon.CompressFileRequest]) (*connect.Response[daemon.CompressFileResponse], error) {
	err := security.CheckServerAccess(ctx, req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	root, err := server.GetRoot(req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer func(root *os.Root) {
		_ = root.Close()
	}(root)

	// TODO

	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}
func (s *ServerFilesServiceHandler) DecompressFile(ctx context.Context, req *connect.Request[daemon.DecompressFileRequest]) (*connect.Response[daemon.DecompressFileResponse], error) {
	err := security.CheckServerAccess(ctx, req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodePermissionDenied, err)
	}

	root, err := server.GetRoot(req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer func(root *os.Root) {
		_ = root.Close()
	}(root)

	// TODO

	return nil, connect.NewError(connect.CodeUnimplemented, nil)
}
