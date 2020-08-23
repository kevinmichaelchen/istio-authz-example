package main

import (
	"sync"

	"github.com/rs/zerolog"

	"github.com/rs/zerolog/log"

	"github.com/kevinmichaelchen/istio-authz-example/internal/app"
	"github.com/kevinmichaelchen/istio-authz-example/internal/configuration"
)

func main() {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Pretty, colorized logger
	log.Logger = configuration.GetLogger()

	c := configuration.LoadConfig()
	log.Info().Msgf("Loaded config: %s", c)

	a := app.NewApp(c)

	var wg sync.WaitGroup

	wg.Add(1)
	go a.Run()

	wg.Wait()
}
