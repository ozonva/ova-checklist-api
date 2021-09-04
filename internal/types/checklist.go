package types

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// ChecklistItem implements fmt.Stringer
type ChecklistItem struct {
	Title      string `json:"title"`
	IsComplete bool   `json:"is_complete"`
}

// Checklist implements fmt.Stringer
type Checklist struct {
	ID          string          `json:"id"`
	UserID      uint64          `json:"user_id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Items       []ChecklistItem `json:"items"`
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

func (c *Checklist) ToJSON() (string, error) {
	result, err := json.Marshal(*c)
	return string(result), err
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

func ChecklistFromJSON(serialized string) (Checklist, error) {
	result := Checklist{}
	err := json.Unmarshal([]byte(serialized), &result)
	return result, err
}

func NewChecklistID() string {
	return uuid.NewString()
}
