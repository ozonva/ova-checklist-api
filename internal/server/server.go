package server

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "ova-checklist-api/pkg/service"
)

type Server interface {
	Start() error
	Wait() error
	Stop() error
}

type service struct {
	pb.UnimplementedChecklistStorageServer
}

// server implements Server
type server struct {
	impl *grpc.Server
	port uint16
	wait sync.WaitGroup
	err  error
}

func (s *service) CreateChecklist(ctx context.Context, request *pb.CreateChecklistRequest) (*pb.CreateChecklistResponse, error) {
	log.Debug().
		Str("handler", "CreateChecklist").
		Str("params", request.String()).
		Send()
	return &pb.CreateChecklistResponse{}, nil
}

func (s *service) DescribeChecklist(ctx context.Context, request *pb.DescribeChecklistRequest) (*pb.DescribeChecklistResponse, error) {
	log.Debug().
		Str("handler", "DescribeChecklist").
		Str("params", request.String()).
		Send()
	return &pb.DescribeChecklistResponse{}, nil
}

func (s *service) ListChecklists(ctx context.Context, request *pb.ListChecklistsRequest) (*pb.ListChecklistsResponse, error) {
	log.Debug().
		Str("handler", "ListChecklists").
		Str("params", request.String()).
		Send()
	return &pb.ListChecklistsResponse{}, nil
}

func (s *service) RemoveChecklist(ctx context.Context, request *pb.RemoveChecklistRequest) (*pb.RemoveChecklistResponse, error) {
	log.Debug().
		Str("handler", "RemoveChecklist").
		Str("params", request.String()).
		Send()
	return &pb.RemoveChecklistResponse{}, nil
}

func New(port uint16) Server {
	s := &server{
		impl: grpc.NewServer(),
		port: port,
	}
	pb.RegisterChecklistStorageServer(s.impl, &service{})
	return s
}

func (s *server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	s.wait.Add(1)
	go func() {
		defer s.wait.Done()
		if err := s.impl.Serve(listener); err != nil {
			s.err = err
		}
	}()

	return nil
}

func (s *server) Wait() error {
	s.wait.Wait()
	return s.err
}

func (s *server) Stop() error {
	s.impl.GracefulStop()
	s.wait.Wait()
	return s.err
}
