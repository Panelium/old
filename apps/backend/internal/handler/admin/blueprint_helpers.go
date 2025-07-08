package admin

import (
	"encoding/json"
	"panelium/backend/internal/model"
	"panelium/proto_gen_go/backend/admin"
)

func BlueprintModelToProto(b *model.Blueprint) *admin.Blueprint {
	if b == nil {
		return nil
	}
	var flags []string
	_ = json.Unmarshal(b.Flags, &flags)
	var dockerImages []*admin.DockerImage
	_ = json.Unmarshal(b.DockerImages, &dockerImages)
	var blockedFiles []*admin.BlockedFile
	_ = json.Unmarshal(b.BlockedFiles, &blockedFiles)
	return &admin.Blueprint{
		FormatVersion:          uint32(b.FormatVersion),
		Bid:                    b.BID,
		Version:                uint32(b.Version),
		UpdateUrl:              b.UpdateURL,
		Name:                   b.Name,
		Description:            b.Description,
		Category:               b.Category,
		Icon:                   b.Icon,
		Banner:                 b.Banner,
		Flags:                  flags,
		DockerImages:           dockerImages,
		BlockedFiles:           blockedFiles,
		ServerBinary:           b.ServerBinary,
		StartCommand:           b.StartCommand,
		StopCommand:            b.StopCommand,
		SetupScriptBase64:      b.SetupScriptBase64,
		SetupDockerImage:       b.SetupDockerImage,
		SetupScriptInterpreter: b.SetupScriptInterpreter,
	}
}

func BlueprintProtoToModel(b *admin.Blueprint) *model.Blueprint {
	if b == nil {
		return nil
	}
	flags, _ := json.Marshal(b.Flags)
	dockerImages, _ := json.Marshal(b.DockerImages)
	blockedFiles, _ := json.Marshal(b.BlockedFiles)
	return &model.Blueprint{
		FormatVersion:          uint(b.FormatVersion),
		BID:                    b.Bid,
		Version:                uint(b.Version),
		UpdateURL:              b.UpdateUrl,
		Name:                   b.Name,
		Description:            b.Description,
		Category:               b.Category,
		Icon:                   b.Icon,
		Banner:                 b.Banner,
		Flags:                  flags,
		DockerImages:           dockerImages,
		BlockedFiles:           blockedFiles,
		ServerBinary:           b.ServerBinary,
		StartCommand:           b.StartCommand,
		StopCommand:            b.StopCommand,
		SetupScriptBase64:      b.SetupScriptBase64,
		SetupDockerImage:       b.SetupDockerImage,
		SetupScriptInterpreter: b.SetupScriptInterpreter,
	}
}
