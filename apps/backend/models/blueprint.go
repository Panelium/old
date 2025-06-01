package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Blueprint struct {
	gorm.Model             `json:"-"`
	FormatVersion          uint           `gorm:"not null" json:"format_version"`  // Version of the blueprint format, used for compatibility checks
	BID                    string         `gorm:"uniqueIndex;not null" json:"bid"` // Unique identifier for the blueprint
	UpdateURL              string         `json:"update_url"`                      // Nullable if not imported from a URL, auto update not possible
	Name                   string         `gorm:"not null" json:"name"`
	Description            string         `gorm:"not null" json:"description"`
	Flags                  datatypes.JSON `gorm:"type:json;not null" json:"flags"` // JSON array of flags that modify the behavior of the blueprint, e.g., eula accept needed for start, server config ui, plugin manager, modpack installer, etc.
	Version                uint           `gorm:"not null" json:"version"`
	DockerImages           datatypes.JSON `gorm:"type:json;not null" json:"docker_images"` // JSON array of Docker images that can be used with this blueprint
	BlockedFiles           datatypes.JSON `gorm:"type:json;not null" json:"blocked_files"` // JSON array of files that the user is not allowed to access or modify
	ServerBinary           string         `json:"server_binary"`                           // Path to the server binary inside the Docker image, e.g., server.jar, server.exe, etc.
	StartCommand           string         `gorm:"not null" json:"start_command"`
	StopCommand            string         `gorm:"not null" json:"stop_command"`
	SetupScriptBase64      string         `gorm:"not null" json:"setup_script_base64"`      // Base64 encoded setup script
	SetupDockerImage       string         `gorm:"not null" json:"setup_docker_image"`       // Docker image used for server setup, can be different from the runtime images
	SetupScriptInterpreter string         `gorm:"not null" json:"setup_script_interpreter"` // Interpreter used for the setup script (e.g., bash, sh, python)
}
