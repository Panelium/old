package server_files

import (
	"connectrpc.com/connect"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"panelium/daemon/internal/security"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go/daemon"
	"path/filepath"
	"strings"
)

func (s *ServerFilesServiceHandler) ListDirectory(ctx context.Context, req *connect.Request[daemon.ListDirectoryRequest]) (*connect.Response[daemon.ListDirectoryResponse], error) {
	err := security.CheckServerAccess(ctx, req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, err)
	}

	root, err := server.GetRoot(req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer func(root *os.Root) {
		_ = root.Close()
	}(root)

	dir, err := root.Open(req.Msg.Path)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	files, err := dir.ReadDir(-1)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var resFiles []*daemon.FileEntry

	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		entry := &daemon.FileEntry{
			Path:         req.Msg.Path + "/" + file.Name(),
			IsDirectory:  file.IsDir(),
			Size:         fileInfo.Size(),
			LastModified: timestamppb.New(fileInfo.ModTime()),
		}
		resFiles = append(resFiles, entry)
	}

	res := &daemon.ListDirectoryResponse{
		Files: resFiles,
	}

	return connect.NewResponse(res), nil
}

func (s *ServerFilesServiceHandler) CreateDirectory(ctx context.Context, req *connect.Request[daemon.CreateDirectoryRequest]) (*connect.Response[daemon.CreateDirectoryResponse], error) {
	err := security.CheckServerAccess(ctx, req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, err)
	}

	root, err := server.GetRoot(req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer func(root *os.Root) {
		_ = root.Close()
	}(root)

	err = rootMkdirAll(root, req.Msg.Path, 0755)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &daemon.CreateDirectoryResponse{
		Success: true,
	}

	return connect.NewResponse(res), nil
}
func (s *ServerFilesServiceHandler) GetDirectorySize(ctx context.Context, req *connect.Request[daemon.GetDirectorySizeRequest]) (*connect.Response[daemon.GetDirectorySizeResponse], error) {
	err := security.CheckServerAccess(ctx, req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, err)
	}

	root, err := server.GetRoot(req.Msg.ServerId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	defer func(root *os.Root) {
		_ = root.Close()
	}(root)

	size, err := rootDirSize(root, req.Msg.Path)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &daemon.GetDirectorySizeResponse{
		Size: size,
	}

	return connect.NewResponse(res), nil
}

func rootDirSize(root *os.Root, path string) (int64, error) {
	dir, err := root.Open(path)
	if err != nil {
		return 0, err
	}

	files, err := dir.ReadDir(-1)
	if err != nil {
		return 0, err
	}

	var totalSize int64

	for _, file := range files {
		if file.IsDir() {
			size, err := rootDirSize(root, path+"/"+file.Name())
			if err != nil {
				return 0, err
			}
			totalSize += size
		} else {
			fileInfo, err := file.Info()
			if err != nil {
				return 0, err
			}
			totalSize += fileInfo.Size()
		}
	}

	return totalSize, nil
}

func rootMkdirAll(root *os.Root, path string, perm os.FileMode) error {
	path = filepath.Clean(path)
	if path == "." {
		return nil
	}

	parts := strings.Split(path, "/")

	var curr string
	for _, part := range parts {
		if part == "" {
			continue
		}

		curr = filepath.Join(curr, part)

		err := root.Mkdir(curr, perm)
		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	return nil
}
