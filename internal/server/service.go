package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ozonva/ova-checklist-api/internal/server/generated/service"
)

func (s *service) handleCreateChecklist(_ context.Context, request *pb.CreateChecklistRequest) (*pb.CreateChecklistResponse, error) {
	if request.Checklist == nil {
		return nil, status.Error(codes.InvalidArgument, "checklist parameter is absent")
	}
	checklist := parseProtoChecklist(request.Checklist)
	if !s.storage.TrySave(checklist) {
		return nil, status.Error(codes.Internal, "unable to save checklist due to an unknown reason")
	}
	return &pb.CreateChecklistResponse{
		ChecklistId: checklist.ID,
	}, nil
}

func (s *service) handleDescribeChecklist(ctx context.Context, request *pb.DescribeChecklistRequest) (*pb.DescribeChecklistResponse, error) {
	if len(request.ChecklistId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "checklist_id parameter is absent")
	}
	checklist, err := s.repository.DescribeChecklist(ctx, request.ChecklistId)
	if err != nil {
		msg := fmt.Sprintf("cannot find a checklist by id %s due to an error: %v", request.ChecklistId, err)
		return nil, status.Error(codes.Internal, msg)
	}
	if checklist == nil {
		msg := fmt.Sprintf("there is no any checklists with id %s", request.ChecklistId)
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
		Checklists: toProtoChecklists(checklists),
	}, nil
}

func (s *service) handleRemoveChecklist(ctx context.Context, request *pb.RemoveChecklistRequest) (*pb.RemoveChecklistResponse, error) {
	if len(request.ChecklistId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "checklist_id parameter is absent")
	}
	if err := s.repository.RemoveChecklist(ctx, request.ChecklistId); err != nil {
		msg := fmt.Sprintf("cannot remove a checklist by id %s due to an error: %v", request.ChecklistId, err)
		return nil, status.Error(codes.Internal, msg)
	}
	return &pb.RemoveChecklistResponse{}, nil
}
