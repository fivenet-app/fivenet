package cmd

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/alecthomas/kong"
	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type SyncCmd struct {
	Status SyncStatusCmd `cmd:"" help:"Check sync status"`
}

type SyncStatusCmd struct {
	Host     string `help:"Host to sync with"`
	Port     int    `default:"8080" help:"Port to sync with"`
	Insecure bool   `help:"Skip TLS verification"`
	APIToken string `help:"API token for authentication"`
}

func (c *SyncStatusCmd) Run(_ *kong.Context) error {
	gCli, err := c.createGRPCClient(c.Host, c.Port, c.Insecure, c.APIToken)
	if err != nil {
		return err
	}
	cli := pbsync.NewSyncServiceClient(gCli)

	ctx := context.Background()

	resp, err := cli.GetStatus(ctx, &pbsync.GetStatusRequest{})
	if err != nil {
		return err
	}

	out, err := protoutils.MarshalToPrettyJSON(resp)
	if err != nil {
		return err
	}
	fmt.Println("Sync Status:\n", string(out))

	return nil
}

func (s *SyncStatusCmd) createGRPCClient(host string, port int, skipTlsVerify bool, apiToken string) (*grpc.ClientConn, error) {
	// Create GRPC client for sync if destination is given
	transportCreds := insecure.NewCredentials()
	if !skipTlsVerify {
		transportCreds = credentials.NewTLS(&tls.Config{
			ClientAuth: tls.NoClientCert,
		})
	}

	cli, err := grpc.NewClient(fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(transportCreds),
		// Require transport security for release mode
		grpc.WithPerRPCCredentials(auth.NewClientTokenAuth(apiToken, !skipTlsVerify)),
	)
	if err != nil {
		return nil, err
	}

	return cli, nil
}
