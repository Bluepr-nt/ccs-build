package services

import (
	"testing"

	"ccs-build.thephoenixhomelab.com/services/mocks"
	"github.com/golang/mock/gomock"
)

func TestCntrSvc_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		username string
		password string
		registry string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Aragorn",
			args: args{
				username: "watman00paradise",
				password: "biggieSmalls",
				registry: "registry.hub.docker.com/",
			},
		},
		{
			name: "",
			args: args{
				username: "watman00paradise",
				password: "biggieSmalls",
				registry: "registry.hub.docker.com/",
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dockerMockClient := mocks.NewMockCriClient(ctrl)
			dockerMockClient.EXPECT().RegistryLogin(gomock.Any(), gomock.Any())
			cri := &CntrSvc{
				client: dockerMockClient,
			}
			if err := cri.Login(tt.args.username, tt.args.password, tt.args.registry); (err != nil) != tt.wantErr {
				t.Errorf("Cri.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
