package sync

import (
	"connectrpc.com/connect"
	"context"
	"encoding/json"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
	"log"
	"net/http"
	"panelium/daemon/internal/config"
	"panelium/daemon/internal/db"
	"panelium/daemon/internal/model"
	"panelium/proto_gen_go"
	"panelium/proto_gen_go/backend"
	"panelium/proto_gen_go/backend/backendconnect"
)

func SyncBlueprints() error {
	client := backendconnect.NewDaemonServiceClient(http.DefaultClient, config.ConfigInstance.GetBackendHost())

	req := connect.NewRequest(&proto_gen_go.Empty{})
	req.Header().Add("Authorization", config.SecretsInstance.BackendToken)

	stm, err := client.SyncBlueprints(context.Background(), req)
	if err != nil {
		return err
	}

	dbInstance := db.Instance()

	defer func(s *connect.ServerStreamForClient[backend.Blueprint]) {
		_ = s.Close()
	}(stm)

	for stm.Receive() {
		blueprint := stm.Msg()
		if blueprint == nil {
			continue
		}

		flagsJson, err := json.Marshal(blueprint.Flags)
		if err != nil {
			return err
		}
		flags := datatypes.JSON(flagsJson)

		dockerImagesJson, err := json.Marshal(blueprint.DockerImages)
		if err != nil {
			return err
		}
		dockerImages := datatypes.JSON(dockerImagesJson)

		blockedFilesJson, err := json.Marshal(blueprint.BlockedFiles)
		if err != nil {
			return err
		}
		blockedFiles := datatypes.JSON(blockedFilesJson)

		dbBlueprint := &model.Blueprint{
			BID:                    blueprint.Bid,
			Version:                uint(blueprint.Version),
			Flags:                  flags,
			DockerImages:           dockerImages,
			BlockedFiles:           blockedFiles,
			ServerBinary:           blueprint.ServerBinary,
			StartCommand:           blueprint.StartCommand,
			StopCommand:            blueprint.StopCommand,
			SetupScriptBase64:      blueprint.SetupScriptBase64,
			SetupDockerImage:       blueprint.SetupDockerImage,
			SetupScriptInterpreter: blueprint.SetupScriptInterpreter,
		}

		tx := dbInstance.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "bid"}},
			DoUpdates: clause.AssignmentColumns([]string{"version", "flags", "docker_images", "blocked_files", "server_binary", "start_command", "stop_command", "setup_script_base64", "setup_docker_image", "setup_script_interpreter"}),
		}).Create(dbBlueprint)
		if tx.Error != nil || tx.RowsAffected == 0 {
			log.Printf("failed to sync blueprint %s: %v", blueprint.Bid, tx.Error)
		}
	}
	if err := stm.Err(); err != nil {
		log.Printf("error while syncing blueprints: %v", err)
		return err
	}

	return nil
}
