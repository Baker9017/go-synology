package universalsearch

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/synology-community/go-synology/pkg/api"
)

func newClient(t *testing.T) Api {
	c, err := api.New(api.Options{
		Host: os.Getenv("SYNOLOGY_HOST"),
	})
	if err != nil {
		t.Error(err)
		require.NoError(t, err)
	}

	if r, err := c.Login(context.Background(), api.LoginOptions{
		Username: os.Getenv("SYNOLOGY_USER"),
		Password: os.Getenv("SYNOLOGY_PASSWORD"),
	}); err != nil {
		t.Error(err)
		require.NoError(t, err)
	} else {
		t.Log("Login successful")
		t.Logf("[INFO] Session: %s\nDeviceID: %s", r.SessionID, r.DeviceID)
	}

	return New(c)
}

func Test_Client_Search(t *testing.T) {
	type fields struct {
		client Api
	}
	type args struct {
		ctx     context.Context
		keyword string
		from    int
		size    int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Search",
			fields: fields{
				client: newClient(t),
			},
			args: args{
				ctx:     context.Background(),
				keyword: "linux",
				from:    0,
				size:    100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.client.Search(tt.args.ctx, tt.args.keyword, tt.args.from, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.NotNil(t, got)
				t.Logf("Total results: %d", got.Total)
				t.Logf("Hits count: %d", len(got.Hits))
				if len(got.Hits) > 0 {
					t.Logf("First hit: %s", got.Hits[0].SYNOMDFSName)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		client api.Api
	}
	tests := []struct {
		name string
		args args
		want Api
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.client)
			require.NotNil(t, got)
		})
	}
}
