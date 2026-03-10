package main

import (
	"os"
	"runtime"

	"github.com/G0tem/go-service-auth/internal/config"
	grpcServer "github.com/G0tem/go-service-auth/internal/grpc"
	"github.com/G0tem/go-service-auth/internal/service/factory"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title Local-Template-Auth Swagger
// @version 1.0
// @description This is an API of auth-service
// @schemes http https

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

// @BasePath /api/v1
func main() {
	// Initialize Zerolog logger with output to stdout
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg := config.LoadConfig()
	zerolog.SetGlobalLevel(zerolog.Level(cfg.LogLevel))

	limitToTwoThreads()

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		if err := grpcServer.StartGrpcServer(&cfg); err != nil {
			log.Error().Msgf("gRPC server error: %v", err)
		}
	}()

	err := factory.StartHttpService(&cfg)
	if err != nil {
		log.Error().Msgf("Attempt to start application fail with error %v", err)
	}
}

func limitToTwoThreads() {
	currentThreadsCount := runtime.GOMAXPROCS(0)

	// Если больше 2, уменьшаем до 2
	if currentThreadsCount > 2 {
		runtime.GOMAXPROCS(2)
	}
}
