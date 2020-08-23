package app

import (
	"sync"

	"github.com/kevinmichaelchen/istio-authz-example/internal/http"

	"github.com/rs/zerolog/log"

	"github.com/kevinmichaelchen/istio-authz-example/internal/grpc"

	"github.com/kevinmichaelchen/istio-authz-example/internal/configuration"
)

type App struct {
	config configuration.Config
}

func NewApp(c configuration.Config) App {
	return App{
		config: c,
	}
}

func (a App) Run() {
	config := a.config

	var wg sync.WaitGroup

	log.Info().Msg("Starting HTTP server...")
	wg.Add(1)
	httpServer := http.NewServer(config)
	go httpServer.Run()

	wg.Add(1)
	envoyAuthV2Server := grpc.NewEnvoyV2Server()
	go envoyAuthV2Server.Run(config.EnvoyAuthzV2Port)

	wg.Add(1)
	envoyAuthV3Server := grpc.NewEnvoyV3Server()
	go envoyAuthV3Server.Run(config.EnvoyAuthzV3Port)

	wg.Wait()
}
