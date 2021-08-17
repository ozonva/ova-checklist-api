package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecklistIsEmpty(t *testing.T) {
	checklist := Checklist{
		UserID:      1,
		Title:       "The Wonderful Project",
		Description: "The Sun is bright, water is wet",
		Items: []ChecklistItem{
			{Title: "Task #1", IsComplete: true},
			{Title: "Task #2", IsComplete: false},
		},
	}
	assert.Equal(t, false, checklist.IsEmpty())
	checklist.Items = nil
	assert.Equal(t, true, checklist.IsEmpty())
}

func TestChecklistIsComplete(t *testing.T) {
	checklist := Checklist{
		UserID:      1,
		Title:       "The Wonderful Project",
		Description: "The Sun is bright, water is wet",
		Items: []ChecklistItem{
			{Title: "Task #1", IsComplete: true},
			{Title: "Task #2", IsComplete: false},
		},
	}
	assert.Equal(t, false, checklist.IsComplete())
	checklist.Items[1].IsComplete = true
	assert.Equal(t, true, checklist.IsComplete())
}
