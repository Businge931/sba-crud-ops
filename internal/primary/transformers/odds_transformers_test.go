package transformers

import (
	"strconv"
	"testing"
	"time"

	"github.com/Businge931/sba-crud-ops/internal/core/domain"
	"github.com/Businge931/sba-crud-ops/proto"
	"github.com/stretchr/testify/require"
)

func TestProtoToCreateOddsRequest(t *testing.T) {
	type dependencies struct{}
	type args struct {
		input *proto.CreateOddsRequest
	}
	testCases := []struct {
		name    string
		deps    dependencies
		args    args
		before  func(t *testing.T, deps *dependencies, args *args)
		after   func(t *testing.T, deps *dependencies, args *args)
		want    domain.CreateOddsRequest
		wantErr bool
	}{
		{
			name: "valid request",
			args: args{
				input: &proto.CreateOddsRequest{
					League:          "EPL",
					HomeTeam:        "Chelsea",
					AwayTeam:        "Arsenal",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 3.1,
					DrawOdds:        2.9,
					GameDate:        time.Date(2025, 6, 11, 13, 0, 0, 0, time.UTC).Format(time.RFC3339),
				},
			},
			want: domain.CreateOddsRequest{
				League:          "EPL",
				HomeTeam:        "Chelsea",
				AwayTeam:        "Arsenal",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 3.1,
				DrawOdds:        2.9,
				GameDate:        time.Date(2025, 6, 11, 13, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "invalid date format",
			args: args{
				input: &proto.CreateOddsRequest{
					League:          "EPL",
					HomeTeam:        "Chelsea",
					AwayTeam:        "Arsenal",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 3.1,
					DrawOdds:        2.9,
					GameDate:        "not-a-date",
				},
			},
			want:    domain.CreateOddsRequest{},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.deps, &tc.args)
			}
			got, err := ProtoToCreateOddsRequest(tc.args.input)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want.League, got.League)
				require.Equal(t, tc.want.HomeTeam, got.HomeTeam)
				require.Equal(t, tc.want.AwayTeam, got.AwayTeam)
				require.Equal(t, tc.want.HomeTeamWinOdds, got.HomeTeamWinOdds)
				require.Equal(t, tc.want.AwayTeamWinOdds, got.AwayTeamWinOdds)
				require.Equal(t, tc.want.DrawOdds, got.DrawOdds)
				require.True(t, tc.want.GameDate.Equal(got.GameDate))
			}
			if tc.after != nil {
				tc.after(t, &tc.deps, &tc.args)
			}
		})
	}
}

func TestProtoToUpdateOddsRequest(t *testing.T) {
	type dependencies struct{}
	type args struct {
		input *proto.UpdateOddsRequest
	}
	testCases := []struct {
		name    string
		deps    dependencies
		args    args
		before  func(t *testing.T, deps *dependencies, args *args)
		after   func(t *testing.T, deps *dependencies, args *args)
		want    domain.CreateOddsRequest
		wantErr bool
	}{
		{
			name: "valid update request",
			args: args{
				input: &proto.UpdateOddsRequest{
					League:          "EPL",
					HomeTeam:        "Chelsea",
					AwayTeam:        "Arsenal",
					HomeTeamWinOdds: 2.7,
					AwayTeamWinOdds: 3.3,
					DrawOdds:        2.8,
					GameDate:        time.Date(2025, 6, 12, 15, 0, 0, 0, time.UTC).Format(time.RFC3339),
				},
			},
			want: domain.CreateOddsRequest{
				League:          "EPL",
				HomeTeam:        "Chelsea",
				AwayTeam:        "Arsenal",
				HomeTeamWinOdds: 2.7,
				AwayTeamWinOdds: 3.3,
				DrawOdds:        2.8,
				GameDate:        time.Date(2025, 6, 12, 15, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "invalid date format",
			args: args{
				input: &proto.UpdateOddsRequest{
					League:          "EPL",
					HomeTeam:        "Chelsea",
					AwayTeam:        "Arsenal",
					HomeTeamWinOdds: 2.7,
					AwayTeamWinOdds: 3.3,
					DrawOdds:        2.8,
					GameDate:        "invalid-date",
				},
			},
			want:    domain.CreateOddsRequest{},
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.deps, &tc.args)
			}
			got, err := ProtoToUpdateOddsRequest(tc.args.input)
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
			if tc.after != nil {
				tc.after(t, &tc.deps, &tc.args)
			}
		})
	}
}

func TestProtoToDeleteOddsRequest(t *testing.T) {
	type dependencies struct{}
	type args struct {
		odds domain.Odds
	}
	testCases := []struct {
		name   string
		deps   dependencies
		args   args
		before func(t *testing.T, deps *dependencies, args *args)
		after  func(t *testing.T, deps *dependencies, args *args)
		want   domain.DeleteOddsRequest
	}{
		{
			name: "basic delete request",
			args: args{
				odds: domain.Odds{
					League:   "EPL",
					HomeTeam: "Chelsea",
					AwayTeam: "Arsenal",
					GameDate: time.Date(2025, 6, 12, 15, 0, 0, 0, time.UTC),
				},
			},
			want: domain.DeleteOddsRequest{
				League:   "EPL",
				HomeTeam: "Chelsea",
				AwayTeam: "Arsenal",
				GameDate: time.Date(2025, 6, 12, 15, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.deps, &tc.args)
			}
			got := ProtoToDeleteOddsRequest(tc.args.odds)
			require.Equal(t, tc.want, got)
			if tc.after != nil {
				tc.after(t, &tc.deps, &tc.args)
			}
		})
	}
}

func TestOddsListToProtoResponses(t *testing.T) {
	type dependencies struct{}
	type args struct {
		oddsList []domain.Odds
	}
	testCases := []struct {
		name   string
		deps   dependencies
		args   args
		before func(t *testing.T, deps *dependencies, args *args)
		after  func(t *testing.T, deps *dependencies, args *args)
		want   []*proto.GetOddsResponse
	}{
		{
			name: "multiple odds",
			args: args{
				oddsList: []domain.Odds{
					{
						ID:              1,
						League:          "EPL",
						HomeTeam:        "Chelsea",
						AwayTeam:        "Arsenal",
						HomeTeamWinOdds: 2.5,
						AwayTeamWinOdds: 3.1,
						DrawOdds:        2.9,
						GameDate:        time.Date(2025, 6, 11, 13, 0, 0, 0, time.UTC),
					},
					{
						ID:              2,
						League:          "EPL",
						HomeTeam:        "Man Utd",
						AwayTeam:        "Liverpool",
						HomeTeamWinOdds: 2.2,
						AwayTeamWinOdds: 3.4,
						DrawOdds:        3.0,
						GameDate:        time.Date(2025, 6, 12, 15, 0, 0, 0, time.UTC),
					},
				},
			},
			want: []*proto.GetOddsResponse{
				{
					OddsId:          strconv.FormatInt(1, 10),
					League:          "EPL",
					HomeTeam:        "Chelsea",
					AwayTeam:        "Arsenal",
					HomeTeamWinOdds: 2.5,
					AwayTeamWinOdds: 3.1,
					DrawOdds:        2.9,
					GameDate:        time.Date(2025, 6, 11, 13, 0, 0, 0, time.UTC).Format(time.RFC3339),
				},
				{
					OddsId:          strconv.FormatInt(2, 10),
					League:          "EPL",
					HomeTeam:        "Man Utd",
					AwayTeam:        "Liverpool",
					HomeTeamWinOdds: 2.2,
					AwayTeamWinOdds: 3.4,
					DrawOdds:        3.0,
					GameDate:        time.Date(2025, 6, 12, 15, 0, 0, 0, time.UTC).Format(time.RFC3339),
				},
			},
		},
		{
			name: "empty list",
			args: args{
				oddsList: []domain.Odds{},
			},
			want: []*proto.GetOddsResponse{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.deps, &tc.args)
			}
			got := OddsListToProtoResponses(tc.args.oddsList)
			require.Equal(t, tc.want, got)
			if tc.after != nil {
				tc.after(t, &tc.deps, &tc.args)
			}
		})
	}
}

func TestFindOddsByID(t *testing.T) {
	type dependencies struct{}
	type args struct {
		oddsList []domain.Odds
		id       string
	}
	testCases := []struct {
		name   string
		deps   dependencies
		args   args
		before func(t *testing.T, deps *dependencies, args *args)
		after  func(t *testing.T, deps *dependencies, args *args)
		wantID int64
		found  bool
	}{
		{
			name: "found first",
			args: args{
				oddsList: []domain.Odds{{ID: 1, League: "EPL"}, {ID: 2, League: "EPL"}},
				id:       "1",
			},
			wantID: 1,
			found:  true,
		},
		{
			name: "found second",
			args: args{
				oddsList: []domain.Odds{{ID: 1, League: "EPL"}, {ID: 2, League: "EPL"}},
				id:       "2",
			},
			wantID: 2,
			found:  true,
		},
		{
			name: "not found",
			args: args{
				oddsList: []domain.Odds{{ID: 1, League: "EPL"}, {ID: 2, League: "EPL"}},
				id:       "3",
			},
			wantID: 0,
			found:  false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.before != nil {
				tc.before(t, &tc.deps, &tc.args)
			}
			got := FindOddsByID(tc.args.oddsList, tc.args.id)
			if tc.found {
				require.NotNil(t, got)
				require.Equal(t, tc.wantID, got.ID)
			} else {
				require.Nil(t, got)
			}
			if tc.after != nil {
				tc.after(t, &tc.deps, &tc.args)
			}
		})
	}
}

func TestOddsToProtoResponse(t *testing.T) {
	testCases := []struct {
		name  string
		input domain.Odds
		want  *proto.GetOddsResponse
	}{
		{
			name: "valid odds",
			input: domain.Odds{
				ID:              123,
				League:          "EPL",
				HomeTeam:        "Chelsea",
				AwayTeam:        "Arsenal",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 3.1,
				DrawOdds:        2.9,
				GameDate:        time.Date(2025, 6, 11, 13, 0, 0, 0, time.UTC),
			},
			want: &proto.GetOddsResponse{
				OddsId:          strconv.FormatInt(123, 10),
				League:          "EPL",
				HomeTeam:        "Chelsea",
				AwayTeam:        "Arsenal",
				HomeTeamWinOdds: 2.5,
				AwayTeamWinOdds: 3.1,
				DrawOdds:        2.9,
				GameDate:        time.Date(2025, 6, 11, 13, 0, 0, 0, time.UTC).Format(time.RFC3339),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := OddsToProtoResponse(tc.input)
			require.Equal(t, tc.want, got)
		})
	}
}
