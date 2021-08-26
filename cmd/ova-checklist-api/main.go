package main

import (
	"context"
	"syscall"

	"github.com/rs/zerolog/log"

	"ova-checklist-api/internal/server"
	"ova-checklist-api/pkg/client"
	"ova-checklist-api/pkg/model"
	"ova-checklist-api/pkg/service"
)

func doSomeCalls(port uint16) {
	const host = "localhost"
	c, err := client.NewClient(host, port)
	log.Error().Msgf("Will send request(s) to: %s:%d", host, port)
	if err != nil {
		log.Error().Msgf("Cannot create a client due to an error: %v", err)
		return
	}
	defer c.Close()

	request := &service.CreateChecklistRequest{
		Checklist: &model.Checklist{
			UserId: 0,
			Title: "Some title",
			Description: "Desc goes here",
		},
	}
	if _, err := c.CreateChecklist(context.Background(), request); err != nil {
		log.Error().Msgf("An error occurred during request: %v", err)
	}
}

func main() {
	const port = uint16(8080) // TODO: read it from config
	s := server.New(port)

	if err := s.Start(); err != nil {
		log.Error().Msgf("Cannot run server due to an error: %v", err)
		syscall.Exit(1)
	}

	go func() {
		doSomeCalls(port)
	}()

	log.Info().Msgf("Server is running on port: %d", port)
	if err := s.Wait(); err != nil {
		log.Error().Msgf("Server stopped with an error: %v", err)
		syscall.Exit(1)
	} else {
		log.Info().Msg("Server stopped gracefully")
	}
}
