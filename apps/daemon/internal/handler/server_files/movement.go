package server_files

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"os"
	"panelium/daemon/internal/security"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go/daemon"
)

func (s *ServerFilesServiceHandler) MoveFile(ctx context.Context, req *connect.Request[daemon.MoveFileRequest]) (*connect.Response[daemon.MoveFileResponse], error) {
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

	stat, err := root.Stat(req.Msg.SourcePath)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	if stat.IsDir() {
		err = rootCopyDirectory(root, req.Msg.SourcePath, req.Msg.DestinationPath, true)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	} else {
		err = rootCopyFile(root, req.Msg.SourcePath, req.Msg.DestinationPath, true)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	res := &daemon.MoveFileResponse{
		Success: true,
	}

	return connect.NewResponse(res), nil
}
func (s *ServerFilesServiceHandler) CopyFile(ctx context.Context, req *connect.Request[daemon.CopyFileRequest]) (*connect.Response[daemon.CopyFileResponse], error) {
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

	stat, err := root.Stat(req.Msg.SourcePath)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	if stat.IsDir() {
		err = rootCopyDirectory(root, req.Msg.SourcePath, req.Msg.DestinationPath, false)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	} else {
		err = rootCopyFile(root, req.Msg.SourcePath, req.Msg.DestinationPath, false)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	res := &daemon.CopyFileResponse{
		Success: true,
	}

	return connect.NewResponse(res), nil
}

func rootCopyFile(root *os.Root, sourcePath string, destinationPath string, move bool) error {
	if sourcePath == destinationPath {
		return nil
	}

	sourceFile, err := root.OpenFile(sourcePath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(sourceFile)

	stat, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	destinationFile, err := root.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, stat.Mode())
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(destinationFile)

	_, err = sourceFile.WriteTo(destinationFile)
	if err != nil {
		return err
	}

	if move {
		err = root.Remove(sourcePath)
		if err != nil {
			return err
		}
	}

	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

func rootCopyDirectory(root *os.Root, sourcePath string, destinationPath string, move bool) error {
	return errors.New("unimplemented")

	// TODO
}
