package server_files

import (
	"connectrpc.com/connect"
	"context"
	"os"
	"panelium/daemon/internal/security"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go/daemon"
)

func (s *ServerFilesServiceHandler) ChangeFilePermissions(ctx context.Context, req *connect.Request[daemon.ChangeFilePermissionsRequest]) (*connect.Response[daemon.ChangeFilePermissionsResponse], error) {
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

	file, err := root.OpenFile(req.Msg.Path, os.O_RDWR, 0)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	err = file.Chmod(os.FileMode(req.Msg.Permissions))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &daemon.ChangeFilePermissionsResponse{
		Success: true,
	}

	return connect.NewResponse(res), nil
}
func (s *ServerFilesServiceHandler) GetFilePermissions(ctx context.Context, req *connect.Request[daemon.GetFilePermissionsRequest]) (*connect.Response[daemon.GetFilePermissionsResponse], error) {
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

	file, err := root.OpenFile(req.Msg.Path, os.O_RDONLY, 0)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	stat, err := file.Stat()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	permissions := stat.Mode().Perm()

	res := &daemon.GetFilePermissionsResponse{
		Permissions: uint32(permissions),
	}

	return connect.NewResponse(res), nil
}
