package configuration

import (
	"encoding/json"

	"github.com/rs/zerolog/log"

	flag "github.com/spf13/pflag"

	"github.com/spf13/viper"
)

const (
	flagForHTTPPort         = "http_port"
	flagForEnvoyAuthzPortV2 = "envoy_authz_v2_port"
	flagForEnvoyAuthzPortV3 = "envoy_authz_v3_port"
	flagForHealthPort       = "health_port"
)

type Config struct {
	// HTTPPort controls what port our HTTP server runs on.
	HTTPPort int

	// Should match with the port we use for envoy.ext_authz
	EnvoyAuthzV2Port int

	// Should match with the port we use for envoy.ext_authz
	EnvoyAuthzV3Port int

	// HealthPort is the port for app health (e.g., liveness probe)
	HealthPort int
}

func (c Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not marshal config to string")
	}
	return string(b)
}

func LoadConfig() Config {
	c := Config{
		HTTPPort:         8081,
		EnvoyAuthzV2Port: 8082,
		EnvoyAuthzV3Port: 8083,
		HealthPort:       8084,
	}

	flag.Int(flagForHTTPPort, c.HTTPPort, "HTTP port")
	flag.Int(flagForEnvoyAuthzPortV2, c.EnvoyAuthzV2Port, "envoy authz v2 port")
	flag.Int(flagForEnvoyAuthzPortV3, c.EnvoyAuthzV3Port, "envoy authz v3 port")
	flag.Int(flagForHealthPort, c.HealthPort, "health port")

	flag.Parse()

	viper.BindPFlag(flagForHTTPPort, flag.Lookup(flagForHTTPPort))
	viper.BindPFlag(flagForEnvoyAuthzPortV2, flag.Lookup(flagForEnvoyAuthzPortV2))
	viper.BindPFlag(flagForEnvoyAuthzPortV3, flag.Lookup(flagForEnvoyAuthzPortV3))
	viper.BindPFlag(flagForHealthPort, flag.Lookup(flagForHealthPort))

	viper.AutomaticEnv()

	c.HTTPPort = viper.GetInt(flagForHTTPPort)
	c.EnvoyAuthzV2Port = viper.GetInt(flagForEnvoyAuthzPortV2)
	c.EnvoyAuthzV3Port = viper.GetInt(flagForEnvoyAuthzPortV3)
	c.HealthPort = viper.GetInt(flagForHealthPort)

	return c
}
