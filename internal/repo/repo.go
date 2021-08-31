package repo

import "github.com/ozonva/ova-checklist-api/internal/types"

// Repo is an interface of a storage which stores entities of type types.Checklist
type Repo interface {
	AddChecklists(checklists []types.Checklist) error
	ListChecklists(limit, offset uint64) ([]types.Checklist, error)
	DescribeChecklist(checklistId uint64) (*types.Checklist, error)
}
