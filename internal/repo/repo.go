package repo

import (
	"context"

	"github.com/ozonva/ova-checklist-api/internal/types"
)

// Repo is an interface of a storage which stores entities of type types.Checklist
type Repo interface {
	AddChecklists(ctx context.Context, checklists []types.Checklist) error
	ListChecklists(ctx context.Context, userId, limit, offset uint64) ([]types.Checklist, error)
	DescribeChecklist(ctx context.Context, checklistId string) (*types.Checklist, error)
	RemoveChecklist(ctx context.Context, checklistId string) error
}
