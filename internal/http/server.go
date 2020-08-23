package http

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	"github.com/kevinmichaelchen/istio-authz-example/internal/configuration"
)

type Server struct {
	config configuration.Config
}

func NewServer(config configuration.Config) Server {
	return Server{
		config: config,
	}
}

func (s Server) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/secure", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secure"))
	}).Methods(http.MethodGet)

	r.HandleFunc("/insecure", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("insecure"))
	}).Methods(http.MethodGet)

	addr := fmt.Sprintf(":%d", s.config.HTTPPort)
	log.Info().Msgf("Listening for HTTP on %s...\n", addr)

	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal().Err(err).Msg("http server failed")
	}
}
