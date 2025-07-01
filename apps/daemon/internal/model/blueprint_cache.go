package model

import (
	"gorm.io/datatypes"
)

type Blueprint struct {
	BID                    string         `gorm:"uniqueIndex;not null" json:"bid"` // Unique identifier for the blueprint
	Version                uint           `gorm:"not null" json:"version"`
	Flags                  datatypes.JSON `gorm:"type:json;not null" json:"flags"`         // JSON array of flags that modify the behavior of the blueprint, e.g., eula accept needed for start, server config ui, plugin manager, modpack installer, etc.
	DockerImages           datatypes.JSON `gorm:"type:json;not null" json:"docker_images"` // JSON array of Docker images that can be used with this blueprint
	BlockedFiles           datatypes.JSON `gorm:"type:json;not null" json:"blocked_files"` // JSON array of files that the user is not allowed to access or modify
	ServerBinary           string         `json:"server_binary"`                           // Path to the server binary inside the server container, e.g., server.jar, server.exe, etc.
	StartCommand           string         `gorm:"not null" json:"start_command"`
	StopCommand            string         `gorm:"not null" json:"stop_command"`
	SetupScriptBase64      string         `gorm:"not null" json:"setup_script_base64"`      // Base64 encoded setup script
	SetupDockerImage       string         `gorm:"not null" json:"setup_docker_image"`       // Docker image used for server setup, can be different from the runtime images
	SetupScriptInterpreter string         `gorm:"not null" json:"setup_script_interpreter"` // Interpreter used for the setup script (e.g., bash, sh, python)
}
