package docker

import (
	"github.com/docker/docker/client"
	"sync"
)

var (
	initOnce     sync.Once
	dockerClient *client.Client
)

func Init() error {
	var err error
	initOnce.Do(func() {
		dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()) // TODO: config
	})

	return err
}

func Instance() *client.Client {
	if dockerClient == nil {
		panic("docker client not initialized, call Init() first")
	}
	return dockerClient
}
