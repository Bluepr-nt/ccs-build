package services

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	docker "github.com/docker/docker/client"
)

//go:generate mockgen -source=cri.go -destination=mocks/mock_services.go -package=mocks . Cri
type CriClient interface {
	RegistryLogin(ctx context.Context, auth types.AuthConfig) (registry.AuthenticateOKBody, error)
}

type ContainerRuntimeInterface interface {
	Login(username, password, registry string) error
}

type Cri struct {
	client CriClient
}

func NewCriService(engineType string) (ContainerRuntimeInterface, error) {
	client, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
	newCri := Cri{client: client}
	return &newCri, nil
}

func (cri *Cri) Login(username, password, registry string) error {
	auth := types.AuthConfig{}
	resp, err := cri.client.RegistryLogin(nil, auth)
	if err != nil {
		return fmt.Errorf("error trying to login to container registry, error: %w, response: %v", err, resp)
	}
	return nil
}
