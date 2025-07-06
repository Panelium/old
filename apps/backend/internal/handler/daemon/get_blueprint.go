package daemon

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"panelium/backend/internal/db"
	"panelium/backend/internal/middleware"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (s *DaemonServiceHandler) GetBlueprint(
	ctx context.Context,
	req *connect.Request[proto_gen_go.SimpleIDMessage],
) (*connect.Response[backend.Blueprint], error) {
	daemonInfoData := ctx.Value("panelium_daemon_info") // this isn't really necessary for this method rn
	daemonInfo, ok := daemonInfoData.(*middleware.DaemonInfo)
	if !ok || daemonInfo == nil || daemonInfo.NID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
	}

	var blueprint *model.Blueprint
	tx := db.Instance().First(&blueprint, "bid = ?", req.Msg.Id)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, connect.NewError(connect.CodeInternal, tx.Error)
	}

	var flags []string
	err := blueprint.Flags.Scan(&flags)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var blockedFiles []struct {
		File     string `json:"file"`
		Visible  bool   `json:"visible"`
		Readable bool   `json:"readable"`
	}
	err = blueprint.BlockedFiles.Scan(&blockedFiles)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var blockedFilesProto []*backend.BlockedFile
	for _, file := range blockedFiles {
		blockedFilesProto = append(blockedFilesProto, &backend.BlockedFile{
			File:     file.File,
			Visible:  file.Visible,
			Readable: file.Readable,
		})
	}

	var dockerImages []struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}
	err = blueprint.DockerImages.Scan(&dockerImages)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	var dockerImagesProto []string
	for _, image := range dockerImages {
		dockerImagesProto = append(dockerImagesProto, image.Image)
	}

	blueprintProto := &backend.Blueprint{
		Bid:                    blueprint.BID,
		Version:                uint32(blueprint.Version),
		Flags:                  flags,
		BlockedFiles:           blockedFilesProto,
		DockerImages:           dockerImagesProto,
		ServerBinary:           blueprint.ServerBinary,
		StartCommand:           blueprint.StartCommand,
		StopCommand:            blueprint.StopCommand,
		SetupScriptBase64:      blueprint.SetupScriptBase64,
		SetupDockerImage:       blueprint.SetupDockerImage,
		SetupScriptInterpreter: blueprint.SetupScriptInterpreter,
	}

	return connect.NewResponse(blueprintProto), nil
}
