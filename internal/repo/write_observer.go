package repo

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"

	"github.com/ozonva/ova-checklist-api/internal/event"
	"github.com/ozonva/ova-checklist-api/internal/eventbus"
	"github.com/ozonva/ova-checklist-api/internal/types"
)

type WriteObserver interface {
	OnAddSuccess(ctx context.Context, checklists []types.Checklist)
	OnRemoveSuccess(ctx context.Context, userId uint64, checklistId string)
	OnUpdateSuccess(ctx context.Context, checklist types.Checklist)
}

type eventBusWriteObserver struct {
	bus eventbus.EventBus
}

func NewWriteObserverOverEventBus(bus eventbus.EventBus) WriteObserver {
	return &eventBusWriteObserver{
		bus: bus,
	}
}

func (e *eventBusWriteObserver) OnAddSuccess(ctx context.Context, checklists []types.Checklist) {
	events := makeEvents(event.EventType_CREATED, checklists...)
	if err := e.bus.Send(ctx, events...); err != nil {
		log.Error().
			Str("reason", "cannot send CREATED event").
			Msgf("%v", err)
	}
}

func (e *eventBusWriteObserver) OnRemoveSuccess(ctx context.Context, userId uint64, checklistId string) {
	ev := makeEvent(event.EventType_REMOVED, userId, checklistId)
	if err := e.bus.Send(ctx, ev); err != nil {
		log.Error().
			Str("reason", "cannot send REMOVED event").
			Msgf("%v", err)
	}
}

func (e *eventBusWriteObserver) OnUpdateSuccess(ctx context.Context, checklist types.Checklist) {
	events := makeEvents(event.EventType_UPDATED, checklist)
	if err := e.bus.Send(ctx, events...); err != nil {
		log.Error().
			Str("reason", "cannot send UPDATED event").
			Msgf("%v", err)
	}
}

func makeEvent(eventType event.EventType, userId uint64, checklistId string) eventbus.Event {
	ev := event.Event{
		UserId:      userId,
		ChecklistId: checklistId,
	}
	serialized, _ := proto.Marshal(&ev)
	return eventbus.Event{
		Key:   event.EventType_name[int32(eventType)],
		Value: serialized,
	}
}

func makeEvents(eventType event.EventType, checklists ...types.Checklist) []eventbus.Event {
	events := make([]eventbus.Event, 0, len(checklists))
	for _, checklist := range checklists {
		events = append(events, makeEvent(eventType, checklist.UserID, checklist.ID))
	}
	return events
}
