package grpc

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/internal/core/ports"
	"github.com/Businge931/sba-crud-ops/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OddsServer struct {
	oddsService ports.OddsService
	proto.UnimplementedOddsServiceServer
}

func NewOddsServer(oddsService ports.OddsService) *OddsServer {
	return &OddsServer{oddsService: oddsService}
}

func (s *OddsServer) CreateOdds(ctx context.Context, req *proto.CreateOddsRequest) (*proto.CreateOddsResponse, error) {
	gameDate, err := time.Parse(time.RFC3339, req.GameDate)
	if err != nil {
		log.Printf("Failed to parse game date: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid game date format: %v", err)
	}

	// Convert proto request to domain request
	createRequest := domain.CreateOddsRequest{
		League:          req.League,
		HomeTeam:        req.HomeTeam,
		AwayTeam:        req.AwayTeam,
		HomeTeamWinOdds: req.HomeTeamWinOdds,
		AwayTeamWinOdds: req.AwayTeamWinOdds,
		DrawOdds:        req.DrawOdds,
		GameDate:        gameDate,
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
		Date:   time.Now(),              // Using current date as we don't have a specific ID lookup
	}

	oddsItems, err := s.oddsService.ReadOdds(ctx, readRequest)
	if err != nil {
		log.Printf("ReadOdds failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get odds: %v", err)
	}

	// If we have results, return the first one as a simplification
	// In a real implementation, you would find the specific odds by ID
	if len(oddsItems) > 0 {
		odds := oddsItems[0]
		return &proto.GetOddsResponse{
			OddsId:          strconv.FormatInt(odds.ID, 10),
			League:          odds.League,
			HomeTeam:        odds.HomeTeam,
			AwayTeam:        odds.AwayTeam,
			HomeTeamWinOdds: odds.HomeTeamWinOdds,
			AwayTeamWinOdds: odds.AwayTeamWinOdds,
			DrawOdds:        odds.DrawOdds,
			GameDate:        odds.GameDate.Format(time.RFC3339),
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, "odds not found")
}

func (s *OddsServer) UpdateOdds(ctx context.Context, req *proto.UpdateOddsRequest) (*proto.UpdateOddsResponse, error) {
	gameDate, err := time.Parse(time.RFC3339, req.GameDate)
	if err != nil {
		log.Printf("Failed to parse game date: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid game date format: %v", err)
	}

	// Convert proto request to domain request
	updateRequest := domain.CreateOddsRequest{ // Reusing CreateOddsRequest for update as per service interface
		League:          req.League,
		HomeTeam:        req.HomeTeam,
		AwayTeam:        req.AwayTeam,
		HomeTeamWinOdds: req.HomeTeamWinOdds,
		AwayTeamWinOdds: req.AwayTeamWinOdds,
		DrawOdds:        req.DrawOdds,
		GameDate:        gameDate,
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
		Date:   time.Now(),              // Using current date as we don't have a specific ID lookup
	}

	oddsItems, err := s.oddsService.ReadOdds(ctx, readRequest)
	if err != nil {
		log.Printf("ReadOdds failed while preparing for delete: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to prepare for delete: %v", err)
	}

	// Find the odds with the matching ID
	// This is a simplification - in a real implementation, you would have a direct lookup
	var targetOdds *domain.Odds
	for _, odds := range oddsItems {
		if strconv.FormatInt(odds.ID, 10) == req.OddsId {
			targetOdds = &odds
			break
		}
	}

	if targetOdds == nil {
		return nil, status.Errorf(codes.NotFound, "odds not found for deletion")
	}

	// Create the delete request with the found information
	deleteRequest := domain.DeleteOddsRequest{
		League:   targetOdds.League,
		HomeTeam: targetOdds.HomeTeam,
		AwayTeam: targetOdds.AwayTeam,
		GameDate: targetOdds.GameDate,
	}

	// Call the service
	err = s.oddsService.DeleteOdds(ctx, deleteRequest)
	if err != nil {
		log.Printf("DeleteOdds failed: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete odds: %v", err)
	}

	return &proto.DeleteOddsResponse{Success: true}, nil
}
