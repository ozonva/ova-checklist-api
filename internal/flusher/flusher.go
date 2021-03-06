package flusher

import (
	"context"
	"log"

	"github.com/ozonva/ova-checklist-api/internal/repo"
	"github.com/ozonva/ova-checklist-api/internal/types"
	"github.com/ozonva/ova-checklist-api/internal/utils"
)

// Flusher is an interface which flushes entities of type types.Checklist to a storage
type Flusher interface {
	Flush(ctx context.Context, checklists []types.Checklist) []types.Checklist
}

type flusher struct {
	chunkSize  uint
	repository repo.Repo
}

// New creates a new Flusher
func New(
	chunkSize uint,
	repository repo.Repo,
) Flusher {
	return &flusher{
		chunkSize:  chunkSize,
		repository: repository,
	}
}

// Flush tries to push checklists into a storage and returns a slice of
// checklists which it failed to push
func (f *flusher) Flush(ctx context.Context, checklists []types.Checklist) []types.Checklist {
	notFlushed := make([]types.Checklist, 0)
	chunks := utils.SplitToChunks(checklists, f.chunkSize)
	for _, chunk := range chunks {
		if err := f.repository.AddChecklists(ctx, chunk); err != nil {
			log.Printf("Unable to flush chunk of checklists to a repository due to an error: %v", err)
			notFlushed = append(notFlushed, chunk...)
		}
	}
	return notFlushed
}
