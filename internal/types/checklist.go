package types

import (
	"encoding/json"
)

type ChecklistItem struct {
	Title string
	IsComplete bool
}

type Checklist struct {
	UserId uint64
	Title string
	Description string
	Items []ChecklistItem
}

func (c *Checklist) String() (string, error) {
	bytes, err := json.Marshal(*c)
	return string(bytes), err
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

func ChecklistFromString(value string) (Checklist, error) {
	var checklist Checklist
	err := json.Unmarshal([]byte(value), &checklist)
	return checklist, err
}
