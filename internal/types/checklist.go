package types

import (
	"fmt"
)

// ChecklistItem implements fmt.Stringer
type ChecklistItem struct {
	Title string
	IsComplete bool
}

// Checklist implements fmt.Stringer
type Checklist struct {
	UserID uint64
	Title  string
	Description string
	Items []ChecklistItem
}

func (i *ChecklistItem) determineStatus() string {
	if i.IsComplete {
		return "Complete"
	}
	return "Incomplete"
}

func (i *ChecklistItem) String() string {
	return fmt.Sprintf("%v: %v", i.Title, i.determineStatus())
}

func (c *Checklist) String() string {
	return fmt.Sprintf("%v: %v", c.Title, c.Description)
}

func (c *Checklist) IsEmpty() bool {
	return len(c.Items) == 0
}

func (c *Checklist) IsComplete() bool {
	for _, item := range c.Items {
		if !item.IsComplete {
			return false
		}
	}
	return true
}
