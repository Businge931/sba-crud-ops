package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockOddsService is a mock implementation of the OddsService interface for testing
type MockOddsService struct {
	mock.Mock
}

func (m *MockOddsService) CreateOdds(ctx context.Context, request domain.CreateOddsRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockOddsService) ReadOdds(ctx context.Context, request domain.ReadOddsRequest) ([]domain.Odds, error) {
	args := m.Called(ctx, request)
	return args.Get(0).([]domain.Odds), args.Error(1)
}

func (m *MockOddsService) UpdateOdds(ctx context.Context, request domain.CreateOddsRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func (m *MockOddsService) DeleteOdds(ctx context.Context, request domain.DeleteOddsRequest) error {
	args := m.Called(ctx, request)
	return args.Error(0)
}

func TestCreateOdds(t *testing.T) {
	// Define test time
	testTime := time.Now().Add(24 * time.Hour)
	testDateStr := testTime.Format(time.RFC3339)

	tests := []struct {
		name    string
		deps    struct {
			mockService *MockOddsService
			ctx        context.Context
		}
		args struct {
			request *proto.CreateOddsRequest
		}
		before  func(deps *struct {
			mockService *MockOddsService
			ctx        context.Context
		})
		after   func()
		want    *proto.CreateOddsResponse
		wantErr error
	}{
		{
			name: "Successful creation",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
			}{mockService: new(MockOddsService), ctx: context.Background()},
			args: struct{ request *proto.CreateOddsRequest }{request: &proto.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        testDateStr,
			}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
			}) {
				deps.mockService.On("CreateOdds", deps.ctx, mock.AnythingOfType("domain.CreateOddsRequest")).Return(nil)
			},
			after:   nil,
			want:    &proto.CreateOddsResponse{OddsId: "created-successfully"},
			wantErr: nil,
		},
		{
			name: "Service error",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
			}{mockService: new(MockOddsService), ctx: context.Background()},
			args: struct{ request *proto.CreateOddsRequest }{request: &proto.CreateOddsRequest{
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        testDateStr,
			}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
			}) {
				deps.mockService.On("CreateOdds", deps.ctx, mock.AnythingOfType("domain.CreateOddsRequest")).Return(domain.ErrInvalidLeague)
			},
			after:   nil,
			want:    nil,
			wantErr: status.Errorf(codes.Internal, "failed to create odds: %v", domain.ErrInvalidLeague),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(&tt.deps)
			}
			server := NewOddsServer(tt.deps.mockService)
			result, err := server.CreateOdds(tt.deps.ctx, tt.args.request)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if statusErr, ok := status.FromError(err); ok {
					expectedStatusErr, _ := status.FromError(tt.wantErr)
					assert.Equal(t, expectedStatusErr.Code(), statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.OddsId, result.OddsId)
			}
			tt.deps.mockService.AssertExpectations(t)
			if tt.after != nil {
				tt.after()
			}
		})
	}
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

	tests := []struct {
		name    string
		deps    struct {
			mockService *MockOddsService
			ctx        context.Context
			odds       []domain.Odds
		}
		args struct {
			request *proto.GetOddsRequest
		}
		before  func(deps *struct {
			mockService *MockOddsService
			ctx        context.Context
			odds       []domain.Odds
		})
		after   func()
		want    *proto.GetOddsResponse
		wantErr error
	}{
		{
			name: "Successful retrieval",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}{mockService: new(MockOddsService), ctx: context.Background(), odds: sampleOdds},
			args: struct{ request *proto.GetOddsRequest }{request: &proto.GetOddsRequest{OddsId: "1"}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}) {
				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return(deps.odds, nil)
			},
			after:   nil,
			want:    &proto.GetOddsResponse{HomeTeam: sampleOdds[0].HomeTeam}, // Only checking HomeTeam for brevity
			wantErr: nil,
		},
		{
			name: "Service error",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}{mockService: new(MockOddsService), ctx: context.Background(), odds: []domain.Odds{}},
			args: struct{ request *proto.GetOddsRequest }{request: &proto.GetOddsRequest{OddsId: "1"}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}) {
				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return([]domain.Odds{}, domain.ErrOddsNotFound)
			},
			after:   nil,
			want:    nil,
			wantErr: status.Errorf(codes.Internal, "failed to get odds: %v", domain.ErrOddsNotFound),
		},
		{
			name: "Empty results",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}{mockService: new(MockOddsService), ctx: context.Background(), odds: []domain.Odds{}},
			args: struct{ request *proto.GetOddsRequest }{request: &proto.GetOddsRequest{OddsId: "1"}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}) {
				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return([]domain.Odds{}, nil)
			},
			after:   nil,
			want:    nil,
			wantErr: status.Errorf(codes.NotFound, "odds not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(&tt.deps)
			}
			server := NewOddsServer(tt.deps.mockService)
			result, err := server.GetOdds(tt.deps.ctx, tt.args.request)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if statusErr, ok := status.FromError(err); ok {
					expectedStatusErr, _ := status.FromError(tt.wantErr)
					assert.Equal(t, expectedStatusErr.Code(), statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.want.HomeTeam, result.HomeTeam)
			}
			tt.deps.mockService.AssertExpectations(t)
			if tt.after != nil {
				tt.after()
			}
		})
	}
}



func TestUpdateOdds(t *testing.T) {
	testTime := time.Now()
	testDateStr := testTime.Format(time.RFC3339)

	tests := []struct {
		name    string
		deps    struct {
			mockService *MockOddsService
			ctx        context.Context
		}
		args struct {
			request *proto.UpdateOddsRequest
		}
		before  func(deps *struct {
			mockService *MockOddsService
			ctx        context.Context
		})
		after   func()
		want    *proto.UpdateOddsResponse
		wantErr error
	}{
		{
			name: "Successful update",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
			}{mockService: new(MockOddsService), ctx: context.Background()},
			args: struct{ request *proto.UpdateOddsRequest }{request: &proto.UpdateOddsRequest{
				OddsId:          "1",
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        testDateStr,
			}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
			}) {
				deps.mockService.On("UpdateOdds", deps.ctx, mock.AnythingOfType("domain.CreateOddsRequest")).Return(nil)
			},
			after:   nil,
			want:    &proto.UpdateOddsResponse{Success: true},
			wantErr: nil,
		},
		{
			name: "Service error",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
			}{mockService: new(MockOddsService), ctx: context.Background()},
			args: struct{ request *proto.UpdateOddsRequest }{request: &proto.UpdateOddsRequest{
				OddsId:          "1",
				League:          "English Premier League",
				HomeTeam:        "Manchester United",
				AwayTeam:        "Liverpool",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 2.1,
				DrawOdds:        3.0,
				GameDate:        testDateStr,
			}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
			}) {
				deps.mockService.On("UpdateOdds", deps.ctx, mock.AnythingOfType("domain.CreateOddsRequest")).Return(domain.ErrOddsNotFound)
			},
			after:   nil,
			want:    nil,
			wantErr: status.Errorf(codes.Internal, "failed to update odds: %v", domain.ErrOddsNotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(&tt.deps)
			}
			server := NewOddsServer(tt.deps.mockService)
			result, err := server.UpdateOdds(tt.deps.ctx, tt.args.request)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if statusErr, ok := status.FromError(err); ok {
					expectedStatusErr, _ := status.FromError(tt.wantErr)
					assert.Equal(t, expectedStatusErr.Code(), statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Success, result.Success)
			}
			tt.deps.mockService.AssertExpectations(t)
			if tt.after != nil {
				tt.after()
			}
		})
	}
}

func TestDeleteOdds(t *testing.T) {
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

	tests := []struct {
		name    string
		deps    struct {
			mockService *MockOddsService
			ctx        context.Context
			odds       []domain.Odds
		}
		args struct {
			request *proto.DeleteOddsRequest
		}
		before  func(deps *struct {
			mockService *MockOddsService
			ctx        context.Context
			odds       []domain.Odds
		})
		after   func()
		want    *proto.DeleteOddsResponse
		wantErr error
	}{
		{
			name: "Successful deletion",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}{mockService: new(MockOddsService), ctx: context.Background(), odds: sampleOdds},
			args: struct{ request *proto.DeleteOddsRequest }{request: &proto.DeleteOddsRequest{OddsId: "1"}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}) {
				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return(deps.odds, nil)
				deps.mockService.On("DeleteOdds", deps.ctx, mock.AnythingOfType("domain.DeleteOddsRequest")).Return(nil)
			},
			after:   nil,
			want:    &proto.DeleteOddsResponse{Success: true},
			wantErr: nil,
		},
		{
			name: "Not found for deletion",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}{mockService: new(MockOddsService), ctx: context.Background(), odds: []domain.Odds{}},
			args: struct{ request *proto.DeleteOddsRequest }{request: &proto.DeleteOddsRequest{OddsId: "1"}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}) {
				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return([]domain.Odds{}, nil)
			},
			after:   nil,
			want:    nil,
			wantErr: status.Errorf(codes.NotFound, "odds not found for deletion"),
		},
		{
			name: "Error fetching for deletion",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}{mockService: new(MockOddsService), ctx: context.Background(), odds: []domain.Odds{}},
			args: struct{ request *proto.DeleteOddsRequest }{request: &proto.DeleteOddsRequest{OddsId: "1"}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}) {
				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return([]domain.Odds{}, domain.ErrInvalidLeague)
			},
			after:   nil,
			want:    nil,
			wantErr: status.Errorf(codes.Internal, "failed to prepare for delete: %v", domain.ErrInvalidLeague),
		},
		{
			name: "Error during deletion",
			deps: struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}{mockService: new(MockOddsService), ctx: context.Background(), odds: sampleOdds},
			args: struct{ request *proto.DeleteOddsRequest }{request: &proto.DeleteOddsRequest{OddsId: "1"}},
			before: func(deps *struct {
				mockService *MockOddsService
				ctx        context.Context
				odds       []domain.Odds
			}) {
				deps.mockService.On("ReadOdds", deps.ctx, mock.AnythingOfType("domain.ReadOddsRequest")).Return(deps.odds, nil)
				deps.mockService.On("DeleteOdds", deps.ctx, mock.AnythingOfType("domain.DeleteOddsRequest")).Return(domain.ErrInvalidConfiguration)
			},
			after:   nil,
			want:    nil,
			wantErr: status.Errorf(codes.Internal, "failed to delete odds: %v", domain.ErrInvalidConfiguration),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(&tt.deps)
			}
			server := NewOddsServer(tt.deps.mockService)
			result, err := server.DeleteOdds(tt.deps.ctx, tt.args.request)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if statusErr, ok := status.FromError(err); ok {
					expectedStatusErr, _ := status.FromError(tt.wantErr)
					assert.Equal(t, expectedStatusErr.Code(), statusErr.Code())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Success, result.Success)
			}
			tt.deps.mockService.AssertExpectations(t)
			if tt.after != nil {
				tt.after()
			}
		})
	}
}


