package application

import (
	"github.com/ozonva/ova-checklist-api/internal/metrics"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/ozonva/ova-checklist-api/internal/config"
	"github.com/ozonva/ova-checklist-api/internal/eventbus"
	"github.com/ozonva/ova-checklist-api/internal/flusher"
	"github.com/ozonva/ova-checklist-api/internal/repo"
	"github.com/ozonva/ova-checklist-api/internal/saver"
)

func buildRepository(pool *pgxpool.Pool, eventBus eventbus.EventBus) repo.Repo {
	observer := repo.NewWriteObserverOverEventBus(eventBus)
	return repo.NewRepoOverDB(pool, observer)
}

func buildSaver(cfg *config.SettingsConfig, repository repo.Repo) saver.Saver {
	return saver.NewSaver(
		flusher.New(
			uint(cfg.RepoFlushBatchSize),
			repository,
		),
		uint(cfg.InternalBufferSize),
		time.Duration(cfg.FlushPeriodMs)*time.Millisecond,
	)
}

func createMetrics() metrics.Metrics {
	defer func() {
		if err := recover(); err != nil {
			log.Error().
				Str("reason", "cannot register metrics").
				Msgf("%v", err)
			doCrash()
		}
	}()
	return metrics.NewMetrics()
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

	sre := runSreServer(&appConfig.Server)
	defer stopSreServer(sre)

	tracingCloser := startTracing(&appConfig.Trace)
	defer stopTracing(tracingCloser)

	pool := connectToDB(&appConfig.Db)
	defer pool.Close()

	eventBus := createEventBus(&appConfig.Kafka)
	defer closeEventBus(eventBus)

	met := createMetrics()
	repository := buildRepository(pool, eventBus)
	storage := buildSaver(&appConfig.Settings, repository)
	defer storage.Close()

	s := runServer(&appConfig.Server, storage, repository, met)
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
