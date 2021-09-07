package application

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/ozonva/ova-checklist-api/internal/config"
	"github.com/ozonva/ova-checklist-api/internal/server"
)

// sreServer implements server.Server
type sreServer struct {
	server *http.Server
	wait   sync.WaitGroup
	err    error
}

func runSreServer(cfg *config.ServerConfig) server.Server {
	srv := sreServer{
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.SrePort),
		},
	}

	http.Handle("/metrics", promhttp.Handler())

	if err := srv.Start(); err != nil {
		log.Error().
			Str("reason", "cannot run the SRE server").
			Msgf("%v", err)
		doCrash()
	}
	return &srv
}

func stopSreServer(s server.Server) {
	if err := s.Stop(); err != nil {
		log.Warn().
			Str("reason", "cannot stop sre server gracefully").
			Msgf("%v", err)
	}
}

func (s *sreServer) Start() error {
	s.wait.Add(1)
	go func() {
		defer s.wait.Done()
		if err := s.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.err = err
			}
		}
	}()
	return nil
}

func (s *sreServer) Wait() error {
	s.wait.Wait()
	return s.err
}

func (s *sreServer) Stop() error {
	if s.server != nil {
		err := s.server.Shutdown(context.Background())
		s.wait.Wait()
		return err
	}
	return nil
}
