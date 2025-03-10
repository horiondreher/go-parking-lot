package utils

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Environment        string `validate:"required" koanf:"ENVIRONMENT"`
	HTTPServerAddress  string `validate:"required" koanf:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress  string `validate:"required" koanf:"GRPC_SERVER_ADDRESS"`
	QueueServerAddress string `validate:"required" koanf:"QUEUE_SERVER_ADDRESS"`
}

var configFile string = ".env"

var (
	k        *koanf.Koanf
	instance *Config
	once     sync.Once
)

// should be used before first call of GetConfig (only for testing)
func SetConfigFile(file string) {
	configFile = file
}

// GetConfig returns the configuration instance using once.Do to ensure that the configuration is loaded only once
func GetConfig() *Config {
	once.Do(func() {
		var err error

		k = koanf.New(".")
		validate := validator.New(validator.WithRequiredStructEnabled())

		log.Info().Msg("loading config...")

		fileProvider := file.Provider(configFile)
		envProvider := env.Provider("", ".", nil)

		err = k.Load(fileProvider, dotenv.Parser())
		if err != nil {
			log.Info().Msgf("could not load config file: %s", err.Error())
		}

		err = k.Load(envProvider, nil)
		if err != nil {
			log.Info().Msgf("could not environment variables: %s", err.Error())
		}

		err = k.Unmarshal("", &instance)
		if err != nil {
			log.Panic().Err(err).Msg("error unmarshing config")
		}

		err = validate.Struct(instance)
		if err != nil {
			log.Panic().Err(err).Msg("correct configs were not loaded")
		}
	})

	return instance
}
