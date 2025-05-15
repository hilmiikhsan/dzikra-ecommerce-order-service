package cmd

import (
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-order-service/internal/infrastructure/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func RunServeGRPC() {
	envs := config.Envs

	logLevel, err := zerolog.ParseLevel(envs.App.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	infrastructure.InitializeLogger(
		envs.App.Environtment,
		envs.App.LogFile,
		logLevel,
	)

	lis, err := net.Listen("tcp", ":"+envs.App.GrpcPort)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen on gRPC port")
	}

	grpcServer := grpc.NewServer()

	// Sync adapters
	err = adapter.Adapters.Sync(
		adapter.WithGRPCServer(grpcServer),
		adapter.WithDzikraPostgres(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to sync adapters")
	}

	// Run gRPC server in a goroutine
	go func() {
		log.Info().Msgf("gRPC server is running on port %s", envs.App.GrpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve gRPC")
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	if runtime.GOOS == "windows" {
		shutdownSignals = []os.Signal{os.Interrupt}
	}
	signal.Notify(quit, shutdownSignals...)
	<-quit

	log.Info().Msg("gRPC server is shutting down ...")
	grpcServer.GracefulStop()

	err = adapter.Adapters.Unsync()
	if err != nil {
		log.Error().Err(err).Msg("Error while closing adapters")
	}

	log.Info().Msg("gRPC server gracefully stopped")
}
