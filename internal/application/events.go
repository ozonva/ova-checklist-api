package application

import (
	"github.com/rs/zerolog/log"

	"github.com/ozonva/ova-checklist-api/internal/config"
	"github.com/ozonva/ova-checklist-api/internal/eventbus"
)

func createEventBus(cfg *config.KafkaConfig) eventbus.EventBus {
	if !cfg.Enabled {
		return eventbus.NewDummyEventBus()
	}
	return eventbus.NewEventBusOverKafka(cfg)
}

func closeEventBus(bus eventbus.EventBus) {
	if err := bus.Close(); err != nil {
		log.Warn().
			Str("reason", "cannot close event bus gracefully").
			Msgf("%v", err)
	}
}
