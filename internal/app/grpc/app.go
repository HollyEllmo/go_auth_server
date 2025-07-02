package app

import (
	"fmt"
	"log/slog"
	"net"

	authrpc "github.com/HollyEllmo/go_auth_server/internal/grpc/auth" // Importing the auth package for gRPC server registration
	"google.golang.org/grpc"
)


type App struct {
	log *slog.Logger
	gRPCServer *grpc.Server
	port int
}

func New(
	log *slog.Logger, 
	authService authrpc.Auth,
	port int,
	) *App {
	gRPCServer := grpc.NewServer()

	authrpc.Register(gRPCServer, authService) // Registering the auth service with the gRPC server

	return &App{
		log: log,
		gRPCServer: gRPCServer,
		port: port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(slog.String("op", op),
	 slog.Int("port", a.port),
	)
	
	i, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return  fmt.Errorf("%s: failed to listen on port %d: %w", op, a.port, err)
	}

	log.Info("grpc server is running", slog.String("address", i.Addr().String()))

    if err := a.gRPCServer.Serve(i); err != nil {
		return fmt.Errorf("%s: failed to serve gRPC server: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}