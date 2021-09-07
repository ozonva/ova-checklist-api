package application

import (
	"github.com/ozonva/ova-checklist-api/internal/metrics"
	"github.com/rs/zerolog/log"

	"github.com/ozonva/ova-checklist-api/internal/config"
	"github.com/ozonva/ova-checklist-api/internal/repo"
	"github.com/ozonva/ova-checklist-api/internal/saver"
	"github.com/ozonva/ova-checklist-api/internal/server"
)

func runServer(
	cfg *config.ServerConfig,
	storage saver.Saver,
	repository repo.Repo,
	met metrics.Metrics,
) server.Server {
	s := server.New(cfg.Port, storage, repository, met)
	if err := s.Start(); err != nil {
		log.Error().
			Str("reason", "cannot run the server").
			Msgf("%v", err)
		doCrash()
	}
	return s
}

func stopServer(s server.Server) {
	if err := s.Stop(); err != nil {
		log.Warn().
			Str("reason", "cannot stop server gracefully").
			Msgf("%v", err)
	}
}
