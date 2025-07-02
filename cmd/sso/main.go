package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/HollyEllmo/go_auth_server/internal/app"
	"github.com/HollyEllmo/go_auth_server/internal/config"
	"github.com/HollyEllmo/go_auth_server/internal/lib/logger/handlers/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: инициализировать объект конфига
	cfg := config.MustLoad()

	// TODO: инициализировать логгер
	log := setupLogger(cfg.Env)

	log.Info("Starting SSO service",
		slog.String("env", cfg.Env),
		slog.String("storage_path", cfg.StoragePath),
		slog.Duration("token_ttl", cfg.TokenTTL),
		slog.Int("grpc_port", cfg.GRPC.Port),
		slog.Duration("grpc_timeout", cfg.GRPC.Timeout),
	)	


	// TODO: инициализировать приложение (app)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)


	// TODO: запустить gRPC-сервер приложения
	
	go application.GRPCServer.MustRun()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping application", slog.String("signal", sign.String()))

	application.GRPCServer.Stop()
	log.Info("SSO service stopped gracefully")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
		case envLocal:
			log = setupPrettySlog()
		case envDev:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelDebug,
			}))
		case envProd:
			log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}