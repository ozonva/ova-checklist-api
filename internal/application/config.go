package application

import (
	"os"

	"github.com/akamensky/argparse"
	"github.com/rs/zerolog/log"

	"github.com/ozonva/ova-checklist-api/internal/config"
)

type applicationArguments struct {
	configPath string
}

func readApplicationArguments() *applicationArguments {
	parser := argparse.NewParser("ova-checklist-api-server", "Processes checklist requests")
	configPath := parser.String("c", "config", &argparse.Options{
		Required: true,
		Help: "Application config path",
	})

	if err := parser.Parse(os.Args); err != nil {
		log.Error().
			Str("reason", "unable to parse application arguments").
			Msgf("%v", err)
		doCrash()
	}

	return &applicationArguments{
		configPath: *configPath,
	}
}

func readApplicationConfig(path string) *config.ApplicationConfig {
	cfg, err := config.ReadApplicationConfig(path)
	if err != nil {
		log.Error().
			Str("reason", "unable to open the application config").
			Msgf("%v", err)
		doCrash()
	}
	return cfg
}
