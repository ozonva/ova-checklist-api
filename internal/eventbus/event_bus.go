package eventbus

import (
	"context"
)

type EventBus interface {
	Send(ctx context.Context, events ...Event) error
	Close() error
}

type Event struct {
	Key   string
	Value []byte
}

type dummyEventBus struct {
}

func NewDummyEventBus() EventBus {
	return &dummyEventBus{}
}

func (d *dummyEventBus) Send(_ context.Context, _ ...Event) error {
	return nil
}

func (d *dummyEventBus) Close() error {
	return nil
}
