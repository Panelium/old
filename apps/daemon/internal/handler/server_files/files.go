package server_files

import (
	"connectrpc.com/connect"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"panelium/daemon/internal/security"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go/daemon"
)

func (s *ServerFilesServiceHandler) ReadFile(ctx context.Context, req *connect.Request[daemon.ReadFileRequest]) (*connect.Response[daemon.ReadFileResponse], error) {
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
	if stat.IsDir() {
		return nil, connect.NewError(connect.CodeInvalidArgument, os.ErrInvalid)
	}

	content := make([]byte, stat.Size())
	_, err = file.Read(content)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	fileInfo := &daemon.FileEntry{
		Path:         req.Msg.Path,
		IsDirectory:  stat.IsDir(),
		Size:         stat.Size(),
		LastModified: timestamppb.New(stat.ModTime()),
	}

	res := &daemon.ReadFileResponse{
		Content:  content,
		FileInfo: fileInfo,
	}

	return connect.NewResponse(res), nil
}
func (s *ServerFilesServiceHandler) WriteFile(ctx context.Context, req *connect.Request[daemon.WriteFileRequest]) (*connect.Response[daemon.WriteFileResponse], error) {
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

	file, err := root.OpenFile(req.Msg.Path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, err = file.Write(req.Msg.Content)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &daemon.WriteFileResponse{
		Success: true,
	}

	return connect.NewResponse(res), nil
}
func (s *ServerFilesServiceHandler) DeleteFile(ctx context.Context, req *connect.Request[daemon.DeleteFileRequest]) (*connect.Response[daemon.DeleteFileResponse], error) {
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

	err = root.Remove(req.Msg.Path)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &daemon.DeleteFileResponse{
		Success: true,
	}

	return connect.NewResponse(res), nil
}
