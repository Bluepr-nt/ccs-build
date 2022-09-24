package services

import (
	docker "github.com/docker/docker/client"
)

type CriClient interface {
}

type ContainerRuntimeInterface interface {

}

type Cri struct {
	criClient CriClient
}

func NewCriService(engineType string) ContainerRuntimeInterface {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

func (cri *Cri) Login(username, password, registry string) {
  docker.
}
