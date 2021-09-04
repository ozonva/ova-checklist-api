package application

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/ozonva/ova-checklist-api/internal/config"
)

func connectToDB(cfg *config.DBConfig) *pgxpool.Pool {
	pool, err := createDBConnection(context.Background(), cfg)
	if err != nil {
		log.Error().
			Str("reason", "unable to connect to the DB").
			Msgf("%v", err)
		doCrash()
	}
	return pool
}

func createDBConnection(ctx context.Context, dbConfig *config.DBConfig) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?pool_max_conns=%d",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DbName,
		dbConfig.MaxConnections,
	)

	pgCfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	retryPeriod := time.Duration(dbConfig.RetryPeriodMs) * time.Millisecond
	tries := int64(dbConfig.ConnectionTries)
	for {
		pool, err := pgxpool.ConnectConfig(ctx, pgCfg)
		if err != nil {
			tries -= 1
			if tries <= 0 {
				log.
					Warn().
					Uint32("tries", max(dbConfig.ConnectionTries, 1)).
					Msg("unable to connect to the DB, giving up.")
				return nil, err
			}
			log.
				Warn().
				Int64("tries_left", tries).
				Uint32("retry_after_ms", dbConfig.RetryPeriodMs).
				Str("reason", "unable to connect to the DB").
				Msgf("%v", err)
			time.Sleep(retryPeriod)
			continue
		}
		return pool, nil
	}
}

func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}
