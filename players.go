package noxaeapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// PlayersModule manages online and offline players: listing, detailed
// attributes, inventory slots, moderation, teleportation, gamemode
// switching, and stats. Access via Client.Players.
type PlayersModule struct{ http requester }

// Location is a teleport target for PlayersModule.Teleport.
type Location struct {
	X, Y, Z float64
	World   string // optional
}

// List returns all currently online players.
func (m *PlayersModule) List(ctx context.Context) ([]OnlinePlayer, error) {
	var out []OnlinePlayer
	err := m.http.doJSON(ctx, "GET", "players", noBody(), nil, &out)
	return out, err
}

// ListAll returns every player the server has ever seen (online + offline).
func (m *PlayersModule) ListAll(ctx context.Context) ([]OfflinePlayer, error) {
	var out []OfflinePlayer
	err := m.http.doJSON(ctx, "GET", "players/all", noBody(), nil, &out)
	return out, err
}

// Get returns the raw JSON for a single player by UUID (works for online or
// offline players) — unmarshal into OnlinePlayer or OfflinePlayer depending
// on which fields you need.
func (m *PlayersModule) Get(ctx context.Context, uuid string) ([]byte, error) {
	return m.http.doRaw(ctx, "GET", fmt.Sprintf("players/%s", uuid), noBody(), nil)
}

// Resolve resolves a player name to their UUID using the server's own
// local player cache (works on both online-mode and offline-mode servers,
// unlike Mojang's public API — the UUID returned matches whatever this
// server actually uses for that player's stats/economy/etc). Checks
// currently online players first, then falls back to the server's offline
// player cache. Returns an *APIError with a 404 status if no known player
// with that name has ever joined.
func (m *PlayersModule) Resolve(ctx context.Context, name string) (PlayerResolveResult, error) {
	var out PlayerResolveResult
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("players/resolve/%s", name), noBody(), nil, &out)
	return out, err
}

// GetInventory returns a player's inventory in a specific world.
func (m *PlayersModule) GetInventory(ctx context.Context, playerUUID, worldUUID string) ([]InventoryItem, error) {
	var out []InventoryItem
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("players/%s/%s/inventory", playerUUID, worldUUID), noBody(), nil, &out)
	return out, err
}

// Kick disconnects an online player, optionally with a reason.
func (m *PlayersModule) Kick(ctx context.Context, uuid, reason string) error {
	body := url.Values{}
	if reason != "" {
		body.Set("reason", reason)
	}
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("players/%s/kick", uuid), formBody(body), nil, nil)
}

// Ban bans a player, optionally with a reason and an ISO-8601 expiry (e.g.
// "2030-01-01T00:00:00Z"). Pass an empty expiry for a permanent ban.
func (m *PlayersModule) Ban(ctx context.Context, uuid, reason, expiry string) error {
	body := url.Values{}
	if reason != "" {
		body.Set("reason", reason)
	}
	if expiry != "" {
		body.Set("expiry", expiry)
	}
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("players/%s/ban", uuid), formBody(body), nil, nil)
}

// Unban removes a player's ban.
func (m *PlayersModule) Unban(ctx context.Context, uuid string) error {
	return m.http.doJSON(ctx, "DELETE", fmt.Sprintf("players/%s/ban", uuid), noBody(), nil, nil)
}

// Teleport teleports a player to a location.
func (m *PlayersModule) Teleport(ctx context.Context, uuid string, loc Location) error {
	body := url.Values{}
	body.Set("x", strconv.FormatFloat(loc.X, 'f', -1, 64))
	body.Set("y", strconv.FormatFloat(loc.Y, 'f', -1, 64))
	body.Set("z", strconv.FormatFloat(loc.Z, 'f', -1, 64))
	if loc.World != "" {
		body.Set("world", loc.World)
	}
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("players/%s/teleport", uuid), formBody(body), nil, nil)
}

// SetGamemode changes a player's gamemode.
func (m *PlayersModule) SetGamemode(ctx context.Context, uuid string, gamemode Gamemode) error {
	body := url.Values{"gamemode": {string(gamemode)}}
	return m.http.doJSON(ctx, "PUT", fmt.Sprintf("players/%s/gamemode", uuid), formBody(body), nil, nil)
}

// GetStats returns kill/death/playtime/block stats for a player.
func (m *PlayersModule) GetStats(ctx context.Context, uuid string) (PlayerStats, error) {
	var out PlayerStats
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("players/%s/stats", uuid), noBody(), nil, &out)
	return out, err
}
