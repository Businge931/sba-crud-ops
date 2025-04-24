package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Businge931/sba-crud-ops/internal/app/transformers"
	"github.com/Businge931/sba-crud-ops/internal/app/usecase"
	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
	"github.com/Businge931/sba-crud-ops/proto"
)

// ImprovedOddsServer uses application use cases for better separation of concerns
type ImprovedOddsServer struct {
	createOddsUseCase *usecase.CreateOddsUseCase
	readOddsUseCase   *usecase.ReadOddsUseCase
	updateOddsUseCase *usecase.UpdateOddsUseCase
	deleteOddsUseCase *usecase.DeleteOddsUseCase
	proto.UnimplementedOddsServiceServer
}

// NewImprovedOddsServer creates a new instance of ImprovedOddsServer with use cases
func NewImprovedOddsServer(
	oddsService ports.OddsService,
) *ImprovedOddsServer {
	return &ImprovedOddsServer{
		createOddsUseCase: usecase.NewCreateOddsUseCase(oddsService),
		readOddsUseCase:   usecase.NewReadOddsUseCase(oddsService),
		updateOddsUseCase: usecase.NewUpdateOddsUseCase(oddsService),
		deleteOddsUseCase: usecase.NewDeleteOddsUseCase(oddsService),
	}
}

// CreateOdds handles the CreateOdds RPC call
func (s *ImprovedOddsServer) CreateOdds(ctx context.Context, req *proto.CreateOddsRequest) (*proto.CreateOddsResponse, error) {
	// Transform request
	createRequest, err := transformers.ProtoToCreateOddsRequest(req)
	if err != nil {
		log.Printf("Failed to transform request: %v", err)
		return nil, err
	}

	// Execute use case
	err = s.createOddsUseCase.Execute(ctx, createRequest)
	if err != nil {
		log.Printf("CreateOdds use case failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create odds: %v", err)
	}

	return &proto.CreateOddsResponse{OddsId: "created-successfully"}, nil
}

// GetOdds handles the GetOdds RPC call
func (s *ImprovedOddsServer) GetOdds(ctx context.Context, req *proto.GetOddsRequest) (*proto.GetOddsResponse, error) {
	// Create read request
	readRequest := domain.ReadOddsRequest{
		League: "English Premier League", // Default league as a simplification
		Date:   time.Now(),
	}

	// Execute use case
	oddsItems, err := s.readOddsUseCase.Execute(ctx, readRequest)
	if err != nil {
		log.Printf("ReadOdds use case failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get odds: %v", err)
	}

	// Return first result as a simplification
	if len(oddsItems) > 0 {
		return transformers.OddsToProtoResponse(oddsItems[0]), nil
	}

	return nil, status.Errorf(codes.NotFound, "odds not found")
}

// UpdateOdds handles the UpdateOdds RPC call
func (s *ImprovedOddsServer) UpdateOdds(ctx context.Context, req *proto.UpdateOddsRequest) (*proto.UpdateOddsResponse, error) {
	// Transform request
	updateRequest, err := transformers.ProtoToUpdateOddsRequest(req)
	if err != nil {
		log.Printf("Failed to transform update request: %v", err)
		return nil, err
	}

	// Execute use case
	err = s.updateOddsUseCase.Execute(ctx, updateRequest)
	if err != nil {
		log.Printf("UpdateOdds use case failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update odds: %v", err)
	}

	return &proto.UpdateOddsResponse{Success: true}, nil
}

// DeleteOdds handles the DeleteOdds RPC call
func (s *ImprovedOddsServer) DeleteOdds(ctx context.Context, req *proto.DeleteOddsRequest) (*proto.DeleteOddsResponse, error) {
	// Since we need full details for delete, we'll first get the odds details
	readRequest := domain.ReadOddsRequest{
		League: "English Premier League",
		Date:   time.Now(),
	}

	// Execute read use case
	oddsItems, err := s.readOddsUseCase.Execute(ctx, readRequest)
	if err != nil {
		log.Printf("ReadOdds use case failed while preparing for delete: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to prepare for delete: %v", err)
	}

	// Find odds by ID
	targetOdds := transformers.FindOddsByID(oddsItems, req.GetOddsId())
	if targetOdds == nil {
		return nil, status.Errorf(codes.NotFound, "odds not found for deletion")
	}

	// Create delete request
	deleteRequest := transformers.ProtoToDeleteOddsRequest(*targetOdds)

	// Execute delete use case
	err = s.deleteOddsUseCase.Execute(ctx, deleteRequest)
	if err != nil {
		log.Printf("DeleteOdds use case failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete odds: %v", err)
	}

	return &proto.DeleteOddsResponse{Success: true}, nil
}
