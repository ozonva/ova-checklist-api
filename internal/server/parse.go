package server

import (
	pb "github.com/ozonva/ova-checklist-api/internal/server/generated/service"
	"github.com/ozonva/ova-checklist-api/internal/types"
)

func parseProtoChecklistItem(protoItem *pb.ChecklistItem) types.ChecklistItem {
	return types.ChecklistItem{
		Title:      protoItem.Title,
		IsComplete: protoItem.IsComplete,
	}
}

func parseProtoChecklist(protoChecklist *pb.Checklist, id *string) types.Checklist {
	items := make([]types.ChecklistItem, 0, len(protoChecklist.Items))
	for _, protoItem := range protoChecklist.Items {
		items = append(items, parseProtoChecklistItem(protoItem))
	}
	return types.Checklist{
		ID:          getChecklistId(id),
		UserID:      protoChecklist.UserId,
		Title:       protoChecklist.Title,
		Description: protoChecklist.Description,
		Items:       items,
	}
}

func getChecklistId(id *string) string {
	if id != nil {
		return *id
	}
	return types.NewChecklistID()
}

func parseProtoChecklists(protoChecklists []*pb.Checklist) []types.Checklist {
	checklists := make([]types.Checklist, 0, len(protoChecklists))
	for _, protoChecklist := range protoChecklists {
		if protoChecklist != nil {
			checklists = append(checklists, parseProtoChecklist(protoChecklist, nil))
		}
	}
	return checklists
}

func toProtoChecklistItem(item *types.ChecklistItem) *pb.ChecklistItem {
	return &pb.ChecklistItem{
		Title:      item.Title,
		IsComplete: item.IsComplete,
	}
}

func toProtoChecklist(checklist *types.Checklist) *pb.Checklist {
	items := make([]*pb.ChecklistItem, 0, len(checklist.Items))
	for _, item := range checklist.Items {
		items = append(items, toProtoChecklistItem(&item))
	}
	return &pb.Checklist{
		UserId:      checklist.UserID,
		Title:       checklist.Title,
		Description: checklist.Description,
		Items:       items,
	}
}

func toProtoUserChecklists(checklists []types.Checklist) []*pb.UserChecklist {
	result := make([]*pb.UserChecklist, 0, len(checklists))
	for _, checklist := range checklists {
		nonUserChecklist := toProtoChecklist(&checklist)
		result = append(result, &pb.UserChecklist{
			Checklist:   nonUserChecklist,
			ChecklistId: checklist.ID,
		})
	}
	return result
}
