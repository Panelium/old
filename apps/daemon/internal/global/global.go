package global

import (
	"github.com/docker/docker/client"
)

var DockerClient *client.Client

func Init() error {
	var err error
	DockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	return nil
}
