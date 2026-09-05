package noxaeapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// LeaderboardModule wraps the generic, pluggable /v1/leaderboards/* routes
// — a unified view over every ranking source registered on the server
// (economy currencies, mcMMO power level, AuraSkills power level, etc).
// Access via Client.Leaderboards.
//
// This module also wraps the player-profile routes
// (/v1/players/{uuid}/profile and /v1/players/profiles), which live
// server-side alongside the leaderboard sources since they're built on top
// of them.
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

// GetPlayerProfile returns a one-call player profile: identity,
// whitelist/ban status, Vault balance, and this player's entry in every
// registered leaderboard source. Built so callers never have to know how
// many leaderboard sources exist or scan top-N lists themselves.
//
// Returns an *APIError with a 404 status if no known player with that UUID
// has ever joined.
func (m *LeaderboardModule) GetPlayerProfile(ctx context.Context, uuid string) (PlayerProfile, error) {
	var out PlayerProfile
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("players/%s/profile", uuid), noBody(), nil, &out)
	return out, err
}

// GetPlayerProfilesOptions configures LeaderboardModule.GetPlayerProfiles.
type GetPlayerProfilesOptions struct {
	// UUIDs fetches an exact, specific set of players. Unknown UUIDs are
	// silently skipped rather than causing an error. When set, Page/Limit/
	// OnlineOnly are ignored.
	UUIDs []string
	// Page is the page number when UUIDs is empty (server default 1).
	Page int
	// Limit is the players per page when UUIDs is empty (server default 50, max 200).
	Limit int
	// OnlineOnly restricts results to currently-online players when UUIDs is empty.
	OnlineOnly bool
}

// GetPlayerProfiles returns bulk player profiles — same shape as
// GetPlayerProfile, for many players at once. Computes each leaderboard
// source's full ranking exactly once for the whole batch rather than once
// per player, so this is far cheaper than calling GetPlayerProfile in a
// loop.
//
// Pass Opts.UUIDs to fetch an exact, specific set of players. Leave it
// empty to page through every known player instead, using Opts.Page/Limit,
// optionally narrowed to only currently-online players via Opts.OnlineOnly.
func (m *LeaderboardModule) GetPlayerProfiles(ctx context.Context, opts GetPlayerProfilesOptions) (PlayerProfilesResponse, error) {
	q := url.Values{}
	if len(opts.UUIDs) > 0 {
		q.Set("uuids", strings.Join(opts.UUIDs, ","))
	} else {
		if opts.Page > 0 {
			q.Set("page", strconv.Itoa(opts.Page))
		}
		if opts.Limit > 0 {
			q.Set("limit", strconv.Itoa(opts.Limit))
		}
		if opts.OnlineOnly {
			q.Set("onlineOnly", "true")
		}
	}

	var out PlayerProfilesResponse
	err := m.http.doJSON(ctx, "GET", "players/profiles", noBody(), q, &out)
	return out, err
}
