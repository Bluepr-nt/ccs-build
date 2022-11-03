package services

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/registry"
	docker "github.com/docker/docker/client"
	"k8s.io/klog/v2"
)

//go:generate mockgen -source=cri.go -destination=mocks/mock_services.go -package=mocks . Cri
type CntrClient interface {
	RegistryLogin(ctx context.Context, auth types.AuthConfig) (registry.AuthenticateOKBody, error)
}

type CntrSvcI interface {
	Login(username, password, registry string) error
}

type CntrSvc struct {
	client CntrClient
}

func NewCntrSvc(engineType string) (CntrSvcI, error) {
	client, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		klog.Errorf("couldn't create a new docker client %w", err)
		return &CntrSvc{}, err
	}

	containers, err := client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		klog.Errorf("coudn't list containers %w", err)
		return &CntrSvc{}, err
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
	newCri := CntrSvc{client: client}
	return &newCri, nil
}

func (cri *CntrSvc) Login(username, password, registry string) error {
	auth := types.AuthConfig{}
	resp, err := cri.client.RegistryLogin(nil, auth)
	if err != nil {
		return fmt.Errorf("error trying to login to container registry, error: %w, response: %v", err, resp)
	}
	return nil
}

//                                                                                                                             _____
//         _____       _____        ___________             ____    _____              ____________  _____    _____       _____\    \
//    _____\    \_   /      |_      \          \        ____\_  \__|\    \            /            \|\    \   \    \     /    / |    |
//   /     /|     | /         \      \    /\    \      /     /     \\\    \          |\___/\  \\___/|\\    \   |    |   /    /  /___/|
//  /     / /____/||     /\    \      |   \_\    |    /     /\      |\\    \          \|____\  \___|/ \\    \  |    |  |    |__ |___|/
// |     | |____|/ |    |  |    \     |      ___/    |     |  |     | \|    | ______        |  |       \|    \ |    |  |       \
// |     |  _____  |     \/      \    |      \  ____ |     |  |     |  |    |/      \  __  /   / __     |     \|    |  |     __/ __
// |\     \|\    \ |\      /\     \  /     /\ \/    \|     | /     /|  /            | /  \/   /_/  |   /     /\      \ |\    \  /  \
// | \_____\|    | | \_____\ \_____\/_____/ |\______||\     \_____/ | /_____/\_____/||____________/|  /_____/ /______/|| \____\/    |
// | |     /____/| | |     | |     ||     | | |     || \_____\   | / |      | |    |||           | / |      | |     | || |    |____/|
//  \|_____|    ||  \|_____|\|_____||_____|/ \|_____| \ |    |___|/  |______|/|____|/|___________|/  |______|/|_____|/  \|____|   | |
//         |____|/                                     \|____|                                                                |___|/
// 				.;;;, .;;;,                   .;;;, .;;;,
// 				.;;;,;;;;;,;;;;;,.;;;,       .;;;.,;;;;;,;;;;;,;;;.
// 			 ;;;;xXXxXXxXXxXXxXXx;;;. .,. .;;;xXXxXXxXXxXXxXX;;;;;
// 	 .,,.`xXX'             `xXXx,;;;;;,xXXx'            `XXx;;,,.
// 	;;;;xXX'                  `xXXx;xXXx'                 `XXx;;;;
// 	`;;XX'                       `XXX'                      `XX;;'
//  ,;;,XX                         `X'                        XX,;;,
//  ;;;;XX,                                                  ,XX;;;;
// 	``.;XX,                                                ,XX;,''
// 		;;;;XX,                                            ,XX;;;;
// 		 ```.;XX,                                        ,XX;,'''
// 				;;;;XX,                                    ,XX;;;;
// 				 ```,;XX,                                ,XX;,'''
// 						 ;;;;XX,                          ,XX;;;;
// 							````,;XX,                    ,XX;, '''
// 									;;;;;XX,              ,XX;;;;
// 									 `````,;XX,        ,XX;,''''
// 												;;;;;XX,  ,XX;;;;;
// 												 `````;;XX;;'''''
// 															`;;;;'
//
