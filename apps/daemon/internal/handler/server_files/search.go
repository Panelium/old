package server_files

import (
	"connectrpc.com/connect"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"panelium/daemon/internal/security"
	"panelium/daemon/internal/server"
	"panelium/proto_gen_go/daemon"
	"strings"
)

func (s *ServerFilesServiceHandler) SearchFiles(ctx context.Context, req *connect.Request[daemon.SearchFilesRequest]) (*connect.Response[daemon.SearchFilesResponse], error) {
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

	results, err := rootSearchFiles(root, req.Msg.Path, req.Msg.Query)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var files []*daemon.FileEntry

	for _, result := range results {
		file, err := root.OpenFile(result, os.O_RDONLY, 0)
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

		files = append(files, &daemon.FileEntry{
			Path:         result,
			IsDirectory:  stat.IsDir(),
			Size:         stat.Size(),
			LastModified: timestamppb.New(stat.ModTime()),
		})
	}

	return nil, nil
}

func rootSearchFiles(root *os.Root, path string, query string) ([]string, error) {
	dir, err := root.Open(path)
	if err != nil {
		return nil, err
	}

	defer func(dir *os.File) {
		_ = dir.Close()
	}(dir)

	entries, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var results []string
	for _, entry := range entries {
		if entry.IsDir() {
			subResults, err := rootSearchFiles(root, path+"/"+entry.Name(), query)
			if err != nil {
				return nil, err
			}
			results = append(results, subResults...)
		}

		if strings.Contains(entry.Name(), query) {
			results = append(results, path+"/"+entry.Name())
		}
	}

	return results, nil
}
