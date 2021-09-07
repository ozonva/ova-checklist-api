package server

import (
	"context"
	"fmt"
	"github.com/ozonva/ova-checklist-api/internal/tracing"
	"net"
	"sync"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	gref "google.golang.org/grpc/reflection"

	"github.com/ozonva/ova-checklist-api/internal/metrics"
	"github.com/ozonva/ova-checklist-api/internal/repo"
	"github.com/ozonva/ova-checklist-api/internal/saver"
	pb "github.com/ozonva/ova-checklist-api/internal/server/generated/service"
)

type Server interface {
	Start() error
	Wait() error
	Stop() error
}

type service struct {
	pb.UnimplementedChecklistStorageServer

	met        metrics.Metrics
	storage    saver.Saver
	repository repo.Repo
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
	ctx, span := tracing.RegisterSpan(ctx, "CreateChecklist")
	defer span.Finish()
	response, err := s.handleCreateChecklist(ctx, request)
	if err != nil {
		s.met.CreateChecklistError()
	} else {
		s.met.CreateChecklistSuccess()
	}
	return response, err
}

func (s *service) MultiCreateChecklist(ctx context.Context, request *pb.MultiCreateChecklistRequest) (*pb.MultiCreateChecklistResponse, error) {
	log.Debug().
		Str("handler", "MultiCreateChecklist").
		Str("params", request.String()).
		Send()
	ctx, span := tracing.RegisterSpan(ctx, "MultiCreateChecklist")
	defer span.Finish()
	response, err := s.handleMultiCreateChecklist(ctx, request)
	if err != nil {
		s.met.MultiCreateChecklistError()
	} else {
		s.met.MultiCreateChecklistSuccess()
	}
	return response, err
}

func (s *service) DescribeChecklist(ctx context.Context, request *pb.DescribeChecklistRequest) (*pb.DescribeChecklistResponse, error) {
	log.Debug().
		Str("handler", "DescribeChecklist").
		Str("params", request.String()).
		Send()
	ctx, span := tracing.RegisterSpan(ctx, "DescribeChecklist")
	defer span.Finish()
	return s.handleDescribeChecklist(ctx, request)
}

func (s *service) ListChecklists(ctx context.Context, request *pb.ListChecklistsRequest) (*pb.ListChecklistsResponse, error) {
	log.Debug().
		Str("handler", "ListChecklists").
		Str("params", request.String()).
		Send()
	ctx, span := tracing.RegisterSpan(ctx, "ListChecklists")
	defer span.Finish()
	return s.handleListChecklists(ctx, request)
}

func (s *service) RemoveChecklist(ctx context.Context, request *pb.RemoveChecklistRequest) (*pb.RemoveChecklistResponse, error) {
	log.Debug().
		Str("handler", "RemoveChecklist").
		Str("params", request.String()).
		Send()
	ctx, span := tracing.RegisterSpan(ctx, "RemoveChecklist")
	defer span.Finish()
	response, err := s.handleRemoveChecklist(ctx, request)
	if err != nil {
		s.met.RemoveChecklistError()
	} else {
		s.met.RemoveChecklistSuccess()
	}
	return response, err
}

func (s *service) UpdateChecklist(ctx context.Context, request *pb.UpdateChecklistRequest) (*pb.UpdateChecklistResponse, error) {
	log.Debug().
		Str("handler", "UpdateChecklist").
		Str("params", request.String()).
		Send()
	ctx, span := tracing.RegisterSpan(ctx, "UpdateChecklist")
	defer span.Finish()
	response, err := s.handleUpdateChecklist(ctx, request)
	if err != nil {
		s.met.UpdateChecklistError()
	} else {
		s.met.UpdateChecklistSuccess()
	}
	return response, err
}

func New(
	port uint16,
	storage saver.Saver,
	repository repo.Repo,
	met metrics.Metrics,
) Server {
	srv := &server{
		impl: grpc.NewServer(),
		port: port,
	}
	svc := &service{
		met:        met,
		storage:    storage,
		repository: repository,
	}
	pb.RegisterChecklistStorageServer(srv.impl, svc)
	gref.Register(srv.impl)
	return srv
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
