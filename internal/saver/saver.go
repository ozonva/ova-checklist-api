package saver

import (
	"context"
	"errors"
	"fmt"
	"github.com/ozonva/ova-checklist-api/internal/tracing"
	"log"
	"sync"
	"time"

	"github.com/ozonva/ova-checklist-api/internal/flusher"
	"github.com/ozonva/ova-checklist-api/internal/types"
)

type Saver interface {
	TrySave(ctx context.Context, checklist types.Checklist) bool

	// TrySaveBatch returns number of successfully saved checklists
	TrySaveBatch(ctx context.Context, checklist []types.Checklist) uint

	Close()
}

// saver implements Saver
type saver struct {
	flusher        flusher.Flusher
	capacity       uint
	flushPeriod    time.Duration
	buffer         []types.Checklist
	waitCompletion sync.WaitGroup
	inputPipe      chan types.Checklist
	stopPipe       chan struct{}
	ctx            context.Context
	ctxCancel      context.CancelFunc
}

func NewSaver(
	flusher flusher.Flusher,
	capacity uint,
	flushPeriod time.Duration,
) Saver {
	result := &saver{
		flusher:        flusher,
		capacity:       capacity,
		flushPeriod:    flushPeriod,
		buffer:         make([]types.Checklist, 0, capacity),
		waitCompletion: sync.WaitGroup{},
		inputPipe:      make(chan types.Checklist, capacity),
		stopPipe:       make(chan struct{}),
	}
	ctx, cancel := context.WithCancel(context.Background())
	result.ctx = ctx
	result.ctxCancel = cancel
	result.runDispatcher()
	return result
}

func (s *saver) TrySave(ctx context.Context, checklist types.Checklist) (ok bool) {
	ok = true
	ctx, span := tracing.RegisterSpan(ctx, "TrySave")
	defer func() {
		if err := recover(); err != nil {
			span.WriteError(errors.New("unable to save checklist, an error occurred"))
			log.Printf("unable to save checklist due to an error: %v", err)
			ok = false
		}
		span.Finish()
	}()
	s.inputPipe <- checklist
	span.WriteInfo("successfully saved checklist")
	return
}

func (s *saver) TrySaveBatch(ctx context.Context, checklists []types.Checklist) (number uint) {
	number = 0
	ctx, span := tracing.RegisterSpan(ctx, "TrySaveBatch")
	defer func() {
		if err := recover(); err != nil {
			span.WriteError(errors.New(fmt.Sprintf("saved only %d of %d checklists, an error occurred", number, len(checklists))))
			log.Printf("unable to save checklists due to an error: %v", err)
		}
		span.Finish()
	}()
	for _, checklist := range checklists {
		s.inputPipe <- checklist
		number++
	}
	span.WriteInfo("successfully saved %d of %d checklists", number, len(checklists))
	return
}

func (s *saver) Close() {
	close(s.inputPipe)
	s.stopPipe <- struct{}{}
	close(s.stopPipe)
	s.waitCompletion.Wait()
	s.ctxCancel()
}

func (s *saver) runDispatcher() {
	s.waitCompletion.Add(1)
	go func() {
		defer s.waitCompletion.Done()
	loop:
		for {
			timer := time.NewTimer(s.flushPeriod)
			select {
			case value, ok := <-s.inputPipe:
				if ok {
					s.buffer = append(s.buffer, value)
					bufferSize := uint(len(s.buffer))
					if bufferSize >= s.capacity {
						s.flush()
					}
				}
				timer.Stop()
			case <-s.stopPipe:
				for value := range s.inputPipe {
					// Save all pending values. NB: inputPipe should be closed by now
					s.buffer = append(s.buffer, value)
				}
				s.flush()
				timer.Stop()
				break loop
			case <-timer.C:
				s.flush()
			}
		}
	}()
}

func (s *saver) flush() {
	if len(s.buffer) > 0 {
		failed := s.flusher.Flush(s.ctx, s.buffer)
		s.buffer = s.buffer[:0]
		s.buffer = append(s.buffer, failed...)
	}
}
