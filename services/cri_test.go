package services

import (
	"ccs-build.thephoenixhomelab.com/services/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestCri_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		client CriClient
	}
	type args struct {
		username string
		password string
		registry string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Aragorn",
			fields: fields{
				client: mocks.NewMockCriClient(ctrl),
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cri := &Cri{
				client: tt.fields.client,
			}
			if err := cri.Login(tt.args.username, tt.args.password, tt.args.registry); (err != nil) != tt.wantErr {
				t.Errorf("Cri.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
