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
	Version                uint           `gorm:"not null" json:"version"`
	DockerImages           datatypes.JSON `gorm:"type:json;not null" json:"docker_images"` // JSON array of Docker images that can be used with this blueprint
	BlockedFiles           datatypes.JSON `gorm:"type:json;not null" json:"blocked_files"` // JSON array of files that the user is not allowed to access or modify
	StartCommand           string         `gorm:"not null" json:"start_command"`
	StopCommand            string         `gorm:"not null" json:"stop_command"`
	SetupScriptBase64      string         `gorm:"not null" json:"setup_script_base64"`      // Base64 encoded setup script
	SetupDockerImage       string         `gorm:"not null" json:"setup_docker_image"`       // Docker image used for server setup, can be different from the runtime images
	SetupScriptInterpreter string         `gorm:"not null" json:"setup_script_interpreter"` // Interpreter used for the setup script (e.g., bash, sh, python)
}

/*
Example JSON representation of a Blueprint:
papermc.bp
{
	"$schema": "https://blueprints.panelium.app/blueprint-v1.schema.json",
	"format_version": 1,
	"bid": "<blueprint_id>",
	"update_url": "https://blueprints.panelium.app/minecraft/java/papermc.bp",
	"name": "PaperMC",
	"description": "The most widely used, high performance Minecraft server that aims to fix gameplay and mechanics inconsistencies",
	"version": 1,
	"docker_images": [
		{
			"name": "Java 17",
			"image": "ghcr.io/panelium/langs/java_17"
		},
		{
			"name": "Java 21",
			"image": "ghcr.io/panelium/langs/java_21"
		}
	],
	"blocked_files": [
		{
			"file": "{{$env::SERVER_BINARY}}",
			"viewable": true, // Optional property (default false), if true, the user can see that the file exists, but still cannot see its contents or modify it
			"accessible": false // Optional property (default false), if true, the user can see the file contents, but cannot modify it
		}
	],
	"start_command": "java -Xms128M -Xmx{{$env::SERVER_MEMORY}}M -jar {{$env::SERVER_BINARY}}",
	//"stop_command": "{{$code::sigint}}",
	"stop_command": "stop",
	"setup_script_base64": "<base64_encoded_setup_script>",
	"setup_docker_image": "ghcr.io/panelium/base:alpine",
	"setup_interpreter": "bash"
}
*/

/*
JSON Schema for Blueprint:
blueprint-v1.schema.json
{
	"$schema": "https://blueprints.panelium.app/blueprint-v1.schema.json",
	"type": "object",
	"properties": {
		"format_version": {
			"type": "integer",
			"description": "Version of the blueprint format, used for compatibility checks"
		},
		"bid": {
			"type": "string",
			"format": "uuid",
			"description": "Unique identifier for the blueprint"
		},
		"update_url": {
			"type": "string",
			"format": "uri",
			"description": "URL to check for updates, can be null if not imported from a URL"
		},
		"name": {
			"type": "string",
			"description": "Name of the blueprint"
		},
		"description": {
			"type": "string",
			"description": "Description of the blueprint"
		},
		"version": {
			"type": "integer",
			"description": "Version of the blueprint, used for compatibility checks"
		},
		"docker_images": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"name": {
						"type": "string",
						"description": "Name of the Docker image"
					},
					"image": {
						"type": "string",
						"format": "uri",
						"description": "Docker image URL"
					}
				},
				"required": ["name", "image"]
			},
			"description": "List of Docker images that can be used with this blueprint"
		},
		"blocked_files": {
			"type": "array",
			"items": {
				"type": "object",
				"properties": {
					"file": {
						"type": "string",
						"description": "Path to the file that is blocked"
					},
					"viewable": {
						"type": "boolean",
						"description": "If true, the user can see that the file exists, but still cannot see its contents or modify it",
						"default": false
					},
					"accessible": {
						"type": "boolean",
						"description": "If true, the user can see the file contents, but cannot modify it",
						"default": false
					}
				},
				"required": ["file"]
			},
			"description": "List of files that the user is not allowed to access or modify"
		},
		"start_command": {
			"type": "string",
			"description": "Command to start the server, can use environment variables like {{$env::SERVER_BINARY}} or {{$env::SERVER_MEMORY}}"
		},
		"stop_command": {
			"type": "string",
			"description": "Command to stop the server, can use stop codes like {{$code::sigint}} for graceful shutdown"
		},
		"setup_script_base64": {
			"type": "string",
			"format": "base64",
			"description": "Base64 encoded setup script that will be executed during server setup"
		},
		"setup_docker_image": {
			"type": "string",
			"format": "uri",
			"description": "Docker image used for server setup, can be different from the runtime images"
		},
		"setup_interpreter": {
			"type": "string",
			"description": "Interpreter used for the setup script (e.g., bash, sh, python)"
		}
	},
	"required": [
		"format_version",
		"bid",
		"name",
		"description",
		"version",
		"docker_images",
		"blocked_files",
		"start_command",
		"stop_command",
		"setup_script_base64",
		"setup_docker_image",
		"setup_interpreter"
	],
	"additionalProperties": false
}
*/
