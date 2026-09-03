package noxaeapi

import (
	"context"
	"fmt"
)

// LeaderboardModule wraps the generic, pluggable /v1/leaderboards/* routes
// — a unified view over every ranking source registered on the server
// (economy currencies, mcMMO power level, AuraSkills power level, etc).
// Access via Client.Leaderboards.
type LeaderboardModule struct{ http *httpEngine }

// List returns every registered leaderboard source and its availability.
func (m *LeaderboardModule) List(ctx context.Context) ([]LeaderboardSourceInfo, error) {
	var out []LeaderboardSourceInfo
	err := m.http.doJSON(ctx, "GET", "leaderboards", noBody(), nil, &out)
	return out, err
}

// GetTop returns ranked entries for one leaderboard source (see List for valid IDs).
func (m *LeaderboardModule) GetTop(ctx context.Context, id string, limit int) ([]LeaderboardEntry, error) {
	var out []LeaderboardEntry
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("leaderboards/%s/top", id), noBody(), limitQuery(limit), &out)
	return out, err
}
