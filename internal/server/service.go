package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonva/ova-checklist-api/internal/server/generated/service"
	"github.com/ozonva/ova-checklist-api/internal/types"
)

func (s *service) handleCreateChecklist(ctx context.Context, request *pb.CreateChecklistRequest) (*pb.CreateChecklistResponse, error) {
	if request.Checklist == nil {
		return nil, status.Error(codes.InvalidArgument, "checklist parameter is absent")
	}
	checklist := parseProtoChecklist(request.Checklist, nil)
	if err := s.repository.AddChecklists(ctx, []types.Checklist{checklist}); err != nil {
		msg := fmt.Sprintf("unable to save checklist due to an error: %v", err)
		return nil, status.Error(codes.Internal, msg)
	}
	return &pb.CreateChecklistResponse{
		ChecklistId: checklist.ID,
	}, nil
}

func (s *service) handleMultiCreateChecklist(ctx context.Context, request *pb.MultiCreateChecklistRequest) (*pb.MultiCreateChecklistResponse, error) {
	checklists := parseProtoChecklists(request.Checklists)
	if len(checklists) == 0 {
		return nil, status.Error(codes.InvalidArgument, "the list of checklists is empty")
	}
	totalSaved := s.storage.TrySaveBatch(ctx, checklists)
	if totalSaved == 0 {
		return nil, status.Error(codes.Internal, "unable to save checklists due to an unknown reason")
	}
	return &pb.MultiCreateChecklistResponse{
		TotalSaved: uint32(totalSaved),
	}, nil
}

func (s *service) handleDescribeChecklist(ctx context.Context, request *pb.DescribeChecklistRequest) (*pb.DescribeChecklistResponse, error) {
	if len(request.ChecklistId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "checklist_id parameter is absent")
	}
	checklist, err := s.repository.DescribeChecklist(ctx, request.UserId, request.ChecklistId)
	if err != nil {
		msg := fmt.Sprintf("cannot find a checklist of user %d with id %s due to an error: %v", request.UserId, request.ChecklistId, err)
		return nil, status.Error(codes.Internal, msg)
	}
	if checklist == nil {
		msg := fmt.Sprintf("there is no any checklists of user %d with id %s", request.UserId, request.ChecklistId)
		return nil, status.Error(codes.NotFound, msg)
	}
	return &pb.DescribeChecklistResponse{
		Checklist: toProtoChecklist(checklist),
	}, nil
}

func (s *service) handleListChecklists(ctx context.Context, request *pb.ListChecklistsRequest) (*pb.ListChecklistsResponse, error) {
	checklists, err := s.repository.ListChecklists(ctx, request.UserId, request.Limit, request.Offset)
	if err != nil {
		msg := fmt.Sprintf("cannot find checklists for user %d due to an error: %v", request.UserId, err)
		return nil, status.Error(codes.Internal, msg)
	}
	return &pb.ListChecklistsResponse{
		Checklists: toProtoUserChecklists(checklists),
	}, nil
}

func (s *service) handleRemoveChecklist(ctx context.Context, request *pb.RemoveChecklistRequest) (*pb.RemoveChecklistResponse, error) {
	if len(request.ChecklistId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "checklist_id parameter is absent")
	}
	if err := s.repository.RemoveChecklist(ctx, request.UserId, request.ChecklistId); err != nil {
		msg := fmt.Sprintf("cannot remove a checklist by id %s due to an error: %v", request.ChecklistId, err)
		return nil, status.Error(codes.Internal, msg)
	}
	return &pb.RemoveChecklistResponse{}, nil
}

func (s *service) handleUpdateChecklist(ctx context.Context, request *pb.UpdateChecklistRequest) (*pb.UpdateChecklistResponse, error) {
	checklist := parseProtoChecklist(request.Checklist, &request.ChecklistId)
	if err := s.repository.UpdateChecklist(ctx, checklist); err != nil {
		msg := fmt.Sprintf("cannot update checklist by id %s for user %d due to an error: %v", checklist.ID, checklist.UserID, err)
		return nil, status.Error(codes.Internal, msg)
	}
	return &pb.UpdateChecklistResponse{}, nil
}
