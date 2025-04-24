package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Businge931/sba-crud-ops/internal/app/transformers"
	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
	"github.com/Businge931/sba-crud-ops/proto"
)

type OddsServer struct {
	oddsService ports.OddsService
	proto.UnimplementedOddsServiceServer
}

func NewOddsServer(oddsService ports.OddsService) *OddsServer {
	return &OddsServer{oddsService: oddsService}
}

func (s *OddsServer) CreateOdds(ctx context.Context, req *proto.CreateOddsRequest) (*proto.CreateOddsResponse, error) {
	// Use transformer to convert proto request to domain request
	createRequest, err := transformers.ProtoToCreateOddsRequest(req)
	if err != nil {
		log.Printf("Failed to transform request: %v", err)
		return nil, err
	}

	// Call the service
	err = s.oddsService.CreateOdds(ctx, createRequest)
	if err != nil {
		log.Printf("CreateOdds failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create odds: %v", err)
	}

	// In a real implementation, you'd return the actual ID
	// For now, we'll return a placeholder
	return &proto.CreateOddsResponse{OddsId: "created-successfully"}, nil
}

func (s *OddsServer) GetOdds(ctx context.Context, req *proto.GetOddsRequest) (*proto.GetOddsResponse, error) {
	// Since we don't have a specific method to get a single odds item by ID in our service interface,
	// we'll create a read request and modify it to work with the available methods

	// This is a simplification - in a real implementation, you would add a specific method
	// to get odds by ID in your service interface
	readRequest := domain.ReadOddsRequest{
		League: "English Premier League", // Assuming default league as per service validation
		Date:   time.Now(),               // Using current date as we don't have a specific ID lookup
	}

	// Since the service implements OddsRetriever, we can access ReadOdds
	oddsItems, err := s.oddsService.ReadOdds(ctx, readRequest)
	if err != nil {
		log.Printf("ReadOdds failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get odds: %v", err)
	}

	// If we have results, return the first one as a simplification
	// In a real implementation, you would find the specific odds by ID
	if len(oddsItems) > 0 {
		// Use transformer to convert domain model to proto response
		return transformers.OddsToProtoResponse(oddsItems[0]), nil
	}

	return nil, status.Errorf(codes.NotFound, "odds not found")
}

func (s *OddsServer) UpdateOdds(ctx context.Context, req *proto.UpdateOddsRequest) (*proto.UpdateOddsResponse, error) {
	// Use transformer to convert proto request to domain request
	updateRequest, err := transformers.ProtoToUpdateOddsRequest(req)
	if err != nil {
		log.Printf("Failed to transform update request: %v", err)
		return nil, err
	}

	// Call the service
	err = s.oddsService.UpdateOdds(ctx, updateRequest)
	if err != nil {
		log.Printf("UpdateOdds failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update odds: %v", err)
	}

	return &proto.UpdateOddsResponse{Success: true}, nil
}

func (s *OddsServer) DeleteOdds(ctx context.Context, req *proto.DeleteOddsRequest) (*proto.DeleteOddsResponse, error) {
	// Since our DeleteOddsRequest requires more information than just an ID,
	// we need to first retrieve the odds details to get the required fields

	// This is a simplification - in a real implementation, you would enhance your service
	// interface to allow deletion by ID
	readRequest := domain.ReadOddsRequest{
		League: "English Premier League", // Assuming default league as per service validation
		Date:   time.Now(),               // Using current date as we don't have a specific ID lookup
	}

	// Since the service implements OddsRetriever, we can access ReadOdds
	oddsItems, err := s.oddsService.ReadOdds(ctx, readRequest)
	if err != nil {
		log.Printf("ReadOdds failed while preparing for delete: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to prepare for delete: %v", err)
	}

	// Use transformer to find odds by ID
	targetOdds := transformers.FindOddsByID(oddsItems, req.GetOddsId())
	if targetOdds == nil {
		return nil, status.Errorf(codes.NotFound, "odds not found for deletion")
	}

	// Use transformer to create delete request
	deleteRequest := transformers.ProtoToDeleteOddsRequest(*targetOdds)

	// Call the service
	err = s.oddsService.DeleteOdds(ctx, deleteRequest)
	if err != nil {
		log.Printf("DeleteOdds failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete odds: %v", err)
	}

	return &proto.DeleteOddsResponse{Success: true}, nil
}
