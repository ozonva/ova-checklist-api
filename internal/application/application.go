package application

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ozonva/ova-checklist-api/internal/config"
	"github.com/ozonva/ova-checklist-api/internal/flusher"
	"github.com/ozonva/ova-checklist-api/internal/repo"
	"github.com/ozonva/ova-checklist-api/internal/saver"
)

func buildSaver(cfg *config.SettingsConfig, repository repo.Repo) saver.Saver {
	return saver.NewSaver(
		flusher.New(
			uint(cfg.RepoFlushBatchSize),
			repository,
		),
		uint(cfg.InternalBufferSize),
		time.Duration(cfg.FlushPeriodMs) * time.Millisecond,
	)
}

func doCrash() {
	log.Error().
		Str("reason", "fatal error occurred").
		Msg("crashing the application")
	panic("fatal error")
}

func Run() {
	args := readApplicationArguments()
	appConfig := readApplicationConfig(args.configPath)

	pool := connectToDB(&appConfig.Db)
	defer pool.Close()

	repository := repo.NewRepoOverDB(pool)
	storage := buildSaver(&appConfig.Settings, repository)
	defer storage.Close()

	s := runServer(&appConfig.Server, storage, repository)
	defer stopServer(s)

	log.Info().
		Uint16("port", appConfig.Server.Port).
		Msg("server is running")
	if err := s.Wait(); err != nil {
		log.Error().
			Str("reason", "server was unexpectedly stopped").
			Msgf("%v", err)
		doCrash()
	} else {
		log.Info().
			Str("reason", "server stopped gracefully").
			Send()
	}
}
