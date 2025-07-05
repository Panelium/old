package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Blueprint struct {
	gorm.Model             `json:"-"`
	FormatVersion          uint           `gorm:"not null" json:"format_version"`  // Version of the blueprint format, used for compatibility checks
	BID                    string         `gorm:"uniqueIndex;not null" json:"bid"` // Unique identifier for the blueprint
	Version                uint           `gorm:"not null" json:"version"`
	UpdateURL              string         `json:"update_url"` // Empty if not imported from a URL -> auto update not possible
	Name                   string         `gorm:"not null" json:"name"`
	Description            string         `gorm:"not null" json:"description"`
	Category               string         `gorm:"not null" json:"category"`                // Category of the blueprint, e.g., Minecraft/Java,Generic,Database,SteamCMD,Games,Storage,Other,...
	Icon                   string         `gorm:"not null" json:"icon"`                    // Base64 encoded icon image for the blueprint, used in UI
	Banner                 string         `gorm:"not null" json:"banner"`                  // Base64 encoded banner image for the blueprint, used in UI
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
