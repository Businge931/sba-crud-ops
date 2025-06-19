package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type testDependencies struct {
	mockService *MockOddsService
	ctx         context.Context
	server      *OddsServer
}

func TestCreateOdds(t *testing.T) {
	// Define test time
	testTime := time.Now().Add(24 * time.Hour)
	testDateStr := testTime.Format(time.RFC3339)

	type testCase struct {
		name   string
		deps   *testDependencies
		before func(t *testing.T, deps *testDependencies) *proto.CreateOddsRequest
		after  func(t *testing.T, deps *testDependencies, response *proto.CreateOddsResponse, err error)
	}

	tests := []testCase{
		{
			name: "Successful creation",
			deps: &testDependencies{
				mockService: new(MockOddsService),
				ctx:         context.Background(),
			},
			before: func(t *testing.T, deps *testDependencies) *proto.CreateOddsRequest {
				req := &proto.CreateOddsRequest{
					League:          "English Premier League",
					HomeTeam:        "Manchester United",
					AwayTeam:        "Liverpool",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        testDateStr,
				}

				deps.mockService.On("CreateOdds", deps.ctx, mock.AnythingOfType("domain.CreateOddsRequest")).
					Return(nil)

				return req
			},
			after: func(t *testing.T, deps *testDependencies, response *proto.CreateOddsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.NotEmpty(t, response.OddsId)

				deps.mockService.AssertExpectations(t)
			},
		},
		{
			name: "Service returns error",
			deps: &testDependencies{
				mockService: new(MockOddsService),
				ctx:         context.Background(),
			},
			before: func(t *testing.T, deps *testDependencies) *proto.CreateOddsRequest {
				req := &proto.CreateOddsRequest{
					League:          "Invalid League",
					HomeTeam:        "Manchester United",
					AwayTeam:        "Liverpool",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        testDateStr,
				}
				deps.mockService.On("CreateOdds", deps.ctx, mock.AnythingOfType("domain.CreateOddsRequest")).Return(domain.ErrInvalidLeague)

				return req
			},
			after: func(t *testing.T, deps *testDependencies, response *proto.CreateOddsResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)

				statusErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.Internal, statusErr.Code())
				assert.Contains(t, statusErr.Message(), domain.ErrInvalidLeague.Error())

				deps.mockService.AssertExpectations(t)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			tc.deps.server = NewOddsServer(tc.deps.mockService)

			// Execute
			req := tc.before(t, tc.deps)
			resp, err := tc.deps.server.CreateOdds(tc.deps.ctx, req)

			// Verify
			tc.after(t, tc.deps, resp, err)
		})
	}
}

type getOddsTestCase struct {
	name   string
	deps   *testDependencies
	before func(t *testing.T, deps *testDependencies) *proto.GetOddsRequest
	after  func(t *testing.T, deps *testDependencies, response *proto.GetOddsResponse, err error)
}

func TestGetOdds(t *testing.T) {
	// Define test time
	testTime := time.Now()

	sampleOdds := []domain.Odds{
		{
			ID:              1,
			League:          "English Premier League",
			HomeTeam:        "Manchester United",
			AwayTeam:        "Liverpool",
			HomeTeamWinOdds: 2.5,
			AwayTeamWinOdds: 2.1,
			DrawOdds:        3.0,
			GameDate:        testTime,
			CreatedAt:       testTime,
			UpdatedAt:       testTime,
		},
	}

	tests := []getOddsTestCase{
		{
			name: "Successful retrieval",
			deps: &testDependencies{
				mockService: new(MockOddsService),
				ctx:         context.Background(),
			},
			before: func(t *testing.T, deps *testDependencies) *proto.GetOddsRequest {
				expectedRequest := domain.ReadOddsRequest{
					League: "English Premier League",
					Date:   time.Now(),
				}
				deps.mockService.On("ReadOdds", deps.ctx, mock.MatchedBy(func(req domain.ReadOddsRequest) bool {
					return req.League == expectedRequest.League &&
						req.Date.Sub(expectedRequest.Date) < time.Second
				})).Return(sampleOdds, nil)
				return &proto.GetOddsRequest{OddsId: "1"}
			},
			after: func(t *testing.T, deps *testDependencies, response *proto.GetOddsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, "Manchester United", response.HomeTeam)
				deps.mockService.AssertExpectations(t)
			},
		},
		{
			name: "Service error",
			deps: &testDependencies{
				mockService: new(MockOddsService),
				ctx:         context.Background(),
			},
			before: func(t *testing.T, deps *testDependencies) *proto.GetOddsRequest {
				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).
					Return([]domain.Odds(nil), domain.ErrOddsNotFound)
				return &proto.GetOddsRequest{OddsId: "1"}
			},
			after: func(t *testing.T, deps *testDependencies, response *proto.GetOddsResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)

				statusErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.Internal, statusErr.Code())
				assert.Contains(t, statusErr.Message(), domain.ErrOddsNotFound.Error())

				deps.mockService.AssertExpectations(t)
			},
		},
		{
			name: "Empty results",
			deps: &testDependencies{
				mockService: new(MockOddsService),
				ctx:         context.Background(),
			},
			before: func(t *testing.T, deps *testDependencies) *proto.GetOddsRequest {
				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).
					Return([]domain.Odds{}, nil)
				return &proto.GetOddsRequest{OddsId: "1"}
			},
			after: func(t *testing.T, deps *testDependencies, response *proto.GetOddsResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)

				statusErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.NotFound, statusErr.Code())
				assert.Contains(t, statusErr.Message(), "odds not found")

				deps.mockService.AssertExpectations(t)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			tc.deps.server = NewOddsServer(tc.deps.mockService)

			// Execute
			req := tc.before(t, tc.deps)
			resp, err := tc.deps.server.GetOdds(tc.deps.ctx, req)

			// Verify results
			tc.after(t, tc.deps, resp, err)
		})
	}
}

type updateOddsTestCase struct {
	name  string
	deps  *testDependencies
	setup func(t *testing.T, deps *testDependencies) *proto.UpdateOddsRequest
	after func(t *testing.T, deps *testDependencies, response *proto.UpdateOddsResponse, err error)
}

func TestUpdateOdds(t *testing.T) {
	testTime := time.Now()
	testDateStr := testTime.Format(time.RFC3339)

	tests := []updateOddsTestCase{
		{
			name: "Successful update",
			deps: &testDependencies{
				mockService: new(MockOddsService),
				ctx:         context.Background(),
			},
			setup: func(t *testing.T, deps *testDependencies) *proto.UpdateOddsRequest {
				req := &proto.UpdateOddsRequest{
					OddsId:          "1",
					League:          "English Premier League",
					HomeTeam:        "Manchester United",
					AwayTeam:        "Liverpool",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        testDateStr,
				}

				deps.mockService.On("UpdateOdds", deps.ctx, mock.AnythingOfType("domain.CreateOddsRequest")).
					Return(nil)

				return req
			},
			after: func(t *testing.T, deps *testDependencies, response *proto.UpdateOddsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Success)
				deps.mockService.AssertExpectations(t)
			},
		},
		{
			name: "Service error",
			deps: &testDependencies{
				mockService: new(MockOddsService),
				ctx:         context.Background(),
			},
			setup: func(t *testing.T, deps *testDependencies) *proto.UpdateOddsRequest {
				req := &proto.UpdateOddsRequest{
					OddsId:          "1",
					League:          "English Premier League",
					HomeTeam:        "Manchester United",
					AwayTeam:        "Liverpool",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        testDateStr,
				}

				deps.mockService.On("UpdateOdds", deps.ctx, mock.AnythingOfType("domain.CreateOddsRequest")).
					Return(domain.ErrOddsNotFound)

				return req
			},
			after: func(t *testing.T, deps *testDependencies, response *proto.UpdateOddsResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)

				statusErr, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, codes.Internal, statusErr.Code())
				assert.Contains(t, statusErr.Message(), domain.ErrOddsNotFound.Error())

				deps.mockService.AssertExpectations(t)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			tc.deps.server = NewOddsServer(tc.deps.mockService)

			// Execute
			req := tc.setup(t, tc.deps)
			resp, err := tc.deps.server.UpdateOdds(tc.deps.ctx, req)

			// Verify
			tc.after(t, tc.deps, resp, err)
		})
	}
}

type deleteOddsTestCase struct {
	name  string
	deps  *testDependencies
	setup func(t *testing.T, deps *testDependencies) *proto.DeleteOddsRequest
	after func(t *testing.T, deps *testDependencies, response *proto.DeleteOddsResponse, err error)
}

func TestDeleteOdds(t *testing.T) {
	testTime := time.Now()
	tests := []deleteOddsTestCase{
		{
			name: "Successful deletion",
			deps: &testDependencies{
				mockService: new(MockOddsService),
				ctx:         context.Background(),
			},
			setup: func(t *testing.T, deps *testDependencies) *proto.DeleteOddsRequest {
				// Set up ReadOdds mock to return a sample odds item
				sampleOdds := []domain.Odds{{
					ID:              1,
					League:          "English Premier League",
					HomeTeam:        "Manchester United",
					AwayTeam:        "Liverpool",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        testTime,
					CreatedAt:       testTime,
					UpdatedAt:       testTime,
				}}

				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).
					Return(sampleOdds, nil)

				req := &proto.DeleteOddsRequest{OddsId: "1"}
				deleteRequest := domain.DeleteOddsRequest{
					League:   "English Premier League",
					HomeTeam: "Manchester United",
					AwayTeam: "Liverpool",
					GameDate: testTime,
				}
				deps.mockService.On("DeleteOdds", deps.ctx, deleteRequest).Return(nil)
				return req
			},
			after: func(t *testing.T, deps *testDependencies, response *proto.DeleteOddsResponse, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Success)
				deps.mockService.AssertExpectations(t)
			},
		},
		{
			name: "Not found error",
			deps: &testDependencies{
				mockService: new(MockOddsService),
				ctx:         context.Background(),
			},
			setup: func(t *testing.T, deps *testDependencies) *proto.DeleteOddsRequest {
				// Set up ReadOdds mock to return a sample odds item
				sampleOdds := []domain.Odds{{
					ID:              1,
					League:          "English Premier League",
					HomeTeam:        "Manchester United",
					AwayTeam:        "Liverpool",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        testTime,
					CreatedAt:       testTime,
					UpdatedAt:       testTime,
				}}

				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).
					Return(sampleOdds, nil)

				req := &proto.DeleteOddsRequest{OddsId: "1"}
				deleteRequest := domain.DeleteOddsRequest{
					League:   "English Premier League",
					HomeTeam: "Manchester United",
					AwayTeam: "Liverpool",
					GameDate: testTime,
				}
				deps.mockService.On("DeleteOdds", deps.ctx, deleteRequest).Return(domain.ErrOddsNotFound)
				return req
			},
			after: func(t *testing.T, deps *testDependencies, response *proto.DeleteOddsResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)

				statusErr, ok := status.FromError(err)
				assert.True(t, ok, "expected gRPC status error")
				assert.Equal(t, codes.Internal, statusErr.Code())
				assert.Contains(t, statusErr.Message(), "failed to delete odds: odds not found")

				deps.mockService.AssertExpectations(t)
			},
		},
		{
			name: "Internal server error",
			deps: &testDependencies{
				mockService: new(MockOddsService),
				ctx:         context.Background(),
			},
			setup: func(t *testing.T, deps *testDependencies) *proto.DeleteOddsRequest {
				// Set up ReadOdds mock to return a sample odds item
				sampleOdds := []domain.Odds{{
					ID:              1,
					League:          "English Premier League",
					HomeTeam:        "Manchester United",
					AwayTeam:        "Liverpool",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 2.1,
					DrawOdds:        3.0,
					GameDate:        testTime,
					CreatedAt:       testTime,
					UpdatedAt:       testTime,
				}}

				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).
					Return(sampleOdds, nil)

				req := &proto.DeleteOddsRequest{OddsId: "1"}
				deleteRequest := domain.DeleteOddsRequest{
					League:   "English Premier League",
					HomeTeam: "Manchester United",
					AwayTeam: "Liverpool",
					GameDate: testTime,
				}
				deps.mockService.On("DeleteOdds", deps.ctx, deleteRequest).
					Return(errors.New("internal server error"))
				return req
			},
			after: func(t *testing.T, deps *testDependencies, response *proto.DeleteOddsResponse, err error) {
				assert.Error(t, err)
				assert.Nil(t, response)

				statusErr, ok := status.FromError(err)
				assert.True(t, ok, "expected gRPC status error")
				assert.Equal(t, codes.Internal, statusErr.Code())
				assert.Contains(t, statusErr.Message(), "failed to delete odds: internal server error")

				deps.mockService.AssertExpectations(t)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			tc.deps.server = NewOddsServer(tc.deps.mockService)

			// Execute
			req := tc.setup(t, tc.deps)
			resp, err := tc.deps.server.DeleteOdds(tc.deps.ctx, req)

			// Verify
			tc.after(t, tc.deps, resp, err)
		})
	}
}
