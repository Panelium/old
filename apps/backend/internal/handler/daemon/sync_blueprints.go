package daemon

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"panelium/backend/internal/db"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
)

func (s *DaemonServiceHandler) SyncBlueprints(
	ctx context.Context,
	req *connect.Request[proto_gen_go.Empty],
	stm *connect.ServerStream[backend.Blueprint],
) error {
	// TODO: bruh
	//daemonInfoData := ctx.Value("panelium_daemon_info")
	//daemonInfo, ok := daemonInfoData.(*middleware.DaemonInfo)
	//if !ok || daemonInfo == nil || daemonInfo.NID == "" {
	//	log.Printf("invalid daemon info in context: %v - %v", daemonInfoData, daemonInfo)
	//	return connect.NewError(connect.CodeUnauthenticated, errors.New("invalid node token"))
	//}

	var blueprints []*model.Blueprint
	tx := db.Instance().Find(&blueprints)
	if tx.Error != nil {
		return connect.NewError(connect.CodeInternal, tx.Error)
	}

	for _, blueprint := range blueprints {
		var flags []string
		err := json.Unmarshal(blueprint.Flags, &flags)
		if err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}

		var blockedFiles []struct {
			File     string `json:"file"`
			Visible  bool   `json:"visible"`
			Readable bool   `json:"readable"`
		}
		err = json.Unmarshal(blueprint.BlockedFiles, &blockedFiles)
		if err != nil {
			return connect.NewError(connect.CodeInternal, err)
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
		err = json.Unmarshal(blueprint.DockerImages, &dockerImages)
		if err != nil {
			return connect.NewError(connect.CodeInternal, err)
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

		if err := stm.Send(blueprintProto); err != nil {
			return connect.NewError(connect.CodeInternal, err)
		}
	}

	return nil
}
