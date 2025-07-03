package suite

import (
	"context"
	"net"
	"strconv"
	"testing"

	"github.com/HollyEllmo/go_auth_server/internal/config"
	ssov1 "github.com/HollyEllmo/my_proto_repo/gen/go/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Cfg *config.Config
	AuthClient ssov1.AuthServiceClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadByPath("../config/local_tests.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	cc, err := grpc.NewClient(grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("failed to create gRPC client: %v", err)
	}

	return ctx, &Suite{
		T:          t,
		Cfg:       cfg,
		AuthClient: ssov1.NewAuthServiceClient(cc),
	}
}

func grpcAddress(cfg *config.Config) string {
	return  net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPC.Port))
}