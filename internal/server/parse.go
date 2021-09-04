package server

import (
	pb "github.com/ozonva/ova-checklist-api/internal/server/generated/service"
	"github.com/ozonva/ova-checklist-api/internal/types"
)

func parseProtoChecklistItem(protoItem *pb.ChecklistItem) types.ChecklistItem {
	return types.ChecklistItem{
		Title: protoItem.Title,
		IsComplete: protoItem.IsComplete,
	}
}

func parseProtoChecklist(protoChecklist *pb.Checklist) types.Checklist {
	items := make([]types.ChecklistItem, 0, len(protoChecklist.Items))
	for _, protoItem := range protoChecklist.Items {
		items = append(items, parseProtoChecklistItem(protoItem))
	}
	return types.Checklist{
		ID: types.NewChecklistID(),
		UserID: protoChecklist.UserId,
		Title: protoChecklist.Title,
		Description: protoChecklist.Description,
		Items: items,
	}
}

func toProtoChecklistItem(item *types.ChecklistItem) *pb.ChecklistItem {
	return &pb.ChecklistItem{
		Title: item.Title,
		IsComplete: item.IsComplete,
	}
}

func toProtoChecklist(checklist *types.Checklist) *pb.Checklist {
	items := make([]*pb.ChecklistItem, 0, len(checklist.Items))
	for _, item := range checklist.Items {
		items = append(items, toProtoChecklistItem(&item))
	}
	return &pb.Checklist{
		UserId: checklist.UserID,
		Title: checklist.Title,
		Description: checklist.Description,
		Items: items,
	}
}

func toProtoChecklists(checklists []types.Checklist) []*pb.Checklist {
	var result []*pb.Checklist
	for _, checklist := range checklists {
		result = append(result, toProtoChecklist(&checklist))
	}
	return result
}
