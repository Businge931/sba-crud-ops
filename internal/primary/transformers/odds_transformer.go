package transformers

import (
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/proto"
)

// ProtoToCreateOddsRequest transforms a proto CreateOddsRequest to a domain CreateOddsRequest
func ProtoToCreateOddsRequest(req *proto.CreateOddsRequest) (domain.CreateOddsRequest, error) {
	gameDate, err := time.Parse(time.RFC3339, req.GetGameDate())
	if err != nil {
		return domain.CreateOddsRequest{}, status.Errorf(codes.InvalidArgument, "invalid game date format: %v", err)
	}

	return domain.CreateOddsRequest{
		League:          req.GetLeague(),
		HomeTeam:        req.GetHomeTeam(),
		AwayTeam:        req.GetAwayTeam(),
		HomeTeamWinOdds: req.GetHomeTeamWinOdds(),
		AwayTeamWinOdds: req.GetAwayTeamWinOdds(),
		DrawOdds:        req.GetDrawOdds(),
		GameDate:        gameDate,
	}, nil
}

// OddsToProtoResponse transforms a domain Odds to a proto GetOddsResponse
func OddsToProtoResponse(odds domain.Odds) *proto.GetOddsResponse {
	return &proto.GetOddsResponse{
		OddsId:          strconv.FormatInt(odds.ID, 10),
		League:          odds.League,
		HomeTeam:        odds.HomeTeam,
		AwayTeam:        odds.AwayTeam,
		HomeTeamWinOdds: odds.HomeTeamWinOdds,
		AwayTeamWinOdds: odds.AwayTeamWinOdds,
		DrawOdds:        odds.DrawOdds,
		GameDate:        odds.GameDate.Format(time.RFC3339),
	}
}

// ProtoToUpdateOddsRequest transforms a proto UpdateOddsRequest to a domain CreateOddsRequest
func ProtoToUpdateOddsRequest(req *proto.UpdateOddsRequest) (domain.CreateOddsRequest, error) {
	gameDate, err := time.Parse(time.RFC3339, req.GetGameDate())
	if err != nil {
		return domain.CreateOddsRequest{}, status.Errorf(codes.InvalidArgument, "invalid game date format: %v", err)
	}

	return domain.CreateOddsRequest{
		League:          req.GetLeague(),
		HomeTeam:        req.GetHomeTeam(),
		AwayTeam:        req.GetAwayTeam(),
		HomeTeamWinOdds: req.GetHomeTeamWinOdds(),
		AwayTeamWinOdds: req.GetAwayTeamWinOdds(),
		DrawOdds:        req.GetDrawOdds(),
		GameDate:        gameDate,
	}, nil
}

// ProtoToDeleteOddsRequest transforms a proto DeleteOddsRequest with odds details to a domain DeleteOddsRequest
func ProtoToDeleteOddsRequest(odds domain.Odds) domain.DeleteOddsRequest {
	return domain.DeleteOddsRequest{
		League:   odds.League,
		HomeTeam: odds.HomeTeam,
		AwayTeam: odds.AwayTeam,
		GameDate: odds.GameDate,
	}
}

// OddsListToProtoResponses transforms a slice of domain Odds to a slice of proto GetOddsResponse
func OddsListToProtoResponses(oddsList []domain.Odds) []*proto.GetOddsResponse {
	responses := make([]*proto.GetOddsResponse, len(oddsList))
	for i, odds := range oddsList {
		responses[i] = OddsToProtoResponse(odds)
	}
	return responses
}

// FindOddsByID finds an odds item by its ID in a slice
func FindOddsByID(oddsList []domain.Odds, id string) *domain.Odds {
	for i, odds := range oddsList {
		if strconv.FormatInt(odds.ID, 10) == id {
			return &oddsList[i]
		}
	}
	return nil
}
