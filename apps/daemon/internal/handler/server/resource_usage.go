package server

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types/container"
	"io"
	"log"
	"panelium/common/fs"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/docker"
	"panelium/daemon/internal/model"
	"panelium/proto_gen_go"
	"time"
)

func (s *ServerServiceHandler) ResourceUsage(
	ctx context.Context,
	req *connect.Request[proto_gen_go.Empty],
	stm *connect.ServerStream[proto_gen_go.ResourceUsageMessage],
) error {
	serverId := ctx.Value("server_id").(string)
	if serverId == "" {
		return connect.NewError(connect.CodeInvalidArgument, errors.New("server ID is required"))
	}

	var srv *model.Server
	tx := db.Instance().First(&srv, "sid = ?", serverId)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return connect.NewError(connect.CodeNotFound, errors.New("server not found"))
	}

	if !srv.ContainerExists {
		return connect.NewError(connect.CodeFailedPrecondition, errors.New("server does not have a container"))
	}

	vol, err := docker.Instance().VolumeInspect(context.Background(), serverId)
	if err != nil {
		return connect.NewError(connect.CodeInternal, errors.New("failed to inspect volume"))
	}

	csr, err := docker.Instance().ContainerStats(context.Background(), serverId, true)
	if err != nil {
		return connect.NewError(connect.CodeInternal, errors.New("failed to get container stats"))
	}

	rc := csr.Body
	defer func(csrc io.ReadCloser) {
		_ = csrc.Close()
	}(rc)

	jsonDecoder := json.NewDecoder(rc)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	var lastStorageMB float32

	go func() {
		for range ticker.C {
			dirSizeBytes, err := fs.DirSize(vol.Mountpoint)
			if err == nil {
				lastStorageMB = float32(dirSizeBytes) / (1024 * 1024)
			}
		}
	}()

	for {
		var stat container.StatsResponse

		if err := jsonDecoder.Decode(&stat); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err) // malformed JSON or other error
		}

		timeNow := stat.Read
		cpuNow := stat.CPUStats
		if timeNow.IsZero() || cpuNow.SystemUsage == 0 {
			continue
		}

		timeBef := stat.PreRead
		cpuBef := stat.PreCPUStats
		if timeBef.IsZero() || cpuBef.SystemUsage == 0 {
			continue
		}

		cpuDelta := cpuNow.SystemUsage - cpuBef.SystemUsage
		if cpuDelta <= 0 {
			continue
		}
		cpuUsage := float32(cpuNow.CPUUsage.TotalUsage-cpuBef.CPUUsage.TotalUsage) / float32(cpuDelta) * 100.0
		if cpuUsage < 0 {
			cpuUsage = 0
		}
		if cpuUsage > 100 {
			cpuUsage = 100
		}

		memNowMB := float32(stat.MemoryStats.Usage) / (1024 * 1024)

		msg := &proto_gen_go.ResourceUsageMessage{
			Usage: &proto_gen_go.ResourceUsage{
				Cpu:     cpuUsage,
				Ram:     memNowMB,
				Storage: lastStorageMB,
			},
		}

		if err := stm.Send(msg); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil // client disconnected
			}
			return connect.NewError(connect.CodeInternal, err)
		}
	}

	return nil
}
