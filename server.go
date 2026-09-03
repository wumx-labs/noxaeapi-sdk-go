package noxaeapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ServerModule covers server-wide operations: liveness, live info, console
// commands, ops, whitelist, restart, logs, entities/chunks, IP bans,
// scoreboards, and chat. Access via Client.Server.
type ServerModule struct{ http *httpEngine }

// Ping is a basic liveness check.
func (m *ServerModule) Ping(ctx context.Context) (map[string]string, error) {
	var out map[string]string
	err := m.http.doJSON(ctx, "GET", "ping", noBody(), nil, &out)
	return out, err
}

// Info returns server info: version, MOTD, TPS, health, player counts, etc.
func (m *ServerModule) Info(ctx context.Context) (ServerInfo, error) {
	var out ServerInfo
	err := m.http.doJSON(ctx, "GET", "server", noBody(), nil, &out)
	return out, err
}

// Exec runs a console command, returning its console output as plain text.
// This is a privileged endpoint — requires a write-enabled API key.
// waitMs is how long to wait for output before returning (server default
// 500ms if 0).
func (m *ServerModule) Exec(ctx context.Context, command string, waitMs int) (string, error) {
	body := url.Values{"command": {command}}
	if waitMs > 0 {
		body.Set("time", strconv.Itoa(waitMs))
	}
	var out string
	err := m.http.doJSON(ctx, "POST", "server/exec", formBody(body), nil, &out)
	return out, err
}

// GetOps lists server operators.
func (m *ServerModule) GetOps(ctx context.Context) ([]WhitelistEntry, error) {
	var out []WhitelistEntry
	err := m.http.doJSON(ctx, "GET", "server/ops", noBody(), nil, &out)
	return out, err
}

// OpPlayer grants operator status to a player.
func (m *ServerModule) OpPlayer(ctx context.Context, uuid string) error {
	body := url.Values{"playerUuid": {uuid}}
	return m.http.doJSON(ctx, "POST", "server/ops", formBody(body), nil, nil)
}

// DeopPlayer revokes operator status from a player. Server-side reads this
// from the query string, not the body.
func (m *ServerModule) DeopPlayer(ctx context.Context, uuid string) error {
	q := url.Values{"playerUuid": {uuid}}
	return m.http.doJSON(ctx, "DELETE", "server/ops", noBody(), q, nil)
}

// GetWhitelist returns the current whitelist.
func (m *ServerModule) GetWhitelist(ctx context.Context) ([]WhitelistEntry, error) {
	var out []WhitelistEntry
	err := m.http.doJSON(ctx, "GET", "server/whitelist", noBody(), nil, &out)
	return out, err
}

// AddToWhitelist adds a player to the whitelist.
func (m *ServerModule) AddToWhitelist(ctx context.Context, uuid, name string) error {
	body := url.Values{"uuid": {uuid}}
	if name != "" {
		body.Set("name", name)
	}
	return m.http.doJSON(ctx, "POST", "server/whitelist", formBody(body), nil, nil)
}

// RemoveFromWhitelist removes a player from the whitelist. Server-side
// reads uuid/name from the query string, not the body.
func (m *ServerModule) RemoveFromWhitelist(ctx context.Context, uuid string) error {
	q := url.Values{"uuid": {uuid}}
	return m.http.doJSON(ctx, "DELETE", "server/whitelist", noBody(), q, nil)
}

// Restart restarts the server. This is a privileged endpoint — requires a write-enabled API key.
func (m *ServerModule) Restart(ctx context.Context) error {
	return m.http.doJSON(ctx, "POST", "server/restart", noBody(), nil, nil)
}

// GetLogs tails the server console log.
func (m *ServerModule) GetLogs(ctx context.Context, lines int) (map[string][]string, error) {
	var out map[string][]string
	err := m.http.doJSON(ctx, "GET", "server/logs", noBody(), limitParam("lines", lines), &out)
	return out, err
}

// GetEntities returns entity counts on the server, optionally scoped to a world.
func (m *ServerModule) GetEntities(ctx context.Context, world string) (map[string]interface{}, error) {
	var out map[string]interface{}
	err := m.http.doJSON(ctx, "GET", "server/entities", noBody(), stringParam("world", world), &out)
	return out, err
}

// GetChunks returns loaded chunk counts, optionally scoped to a world.
func (m *ServerModule) GetChunks(ctx context.Context, world string) (map[string]interface{}, error) {
	var out map[string]interface{}
	err := m.http.doJSON(ctx, "GET", "server/chunks", noBody(), stringParam("world", world), &out)
	return out, err
}

// BanIP bans an IP address.
func (m *ServerModule) BanIP(ctx context.Context, ip, reason string) error {
	body := url.Values{"ip": {ip}}
	if reason != "" {
		body.Set("reason", reason)
	}
	return m.http.doJSON(ctx, "POST", "server/ban-ip", formBody(body), nil, nil)
}

// GetObjective returns a scoreboard objective's scores by objective name.
func (m *ServerModule) GetObjective(ctx context.Context, name string) (map[string]interface{}, error) {
	var out map[string]interface{}
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("scoreboard/%s", name), noBody(), nil, &out)
	return out, err
}

// GetScoreboard lists all scoreboard objectives and tracked entries.
func (m *ServerModule) GetScoreboard(ctx context.Context) (map[string]interface{}, error) {
	var out map[string]interface{}
	err := m.http.doJSON(ctx, "GET", "scoreboard", noBody(), nil, &out)
	return out, err
}

// SetScore sets a score for an entry on an objective.
func (m *ServerModule) SetScore(ctx context.Context, objective, entry string, value int64) error {
	body := url.Values{"entry": {entry}, "value": {strconv.FormatInt(value, 10)}}
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("scoreboard/%s/score", objective), formBody(body), nil, nil)
}

// ResetScore resets (removes) a score for an entry on an objective.
// Server-side reads entry from the query string, not the body.
func (m *ServerModule) ResetScore(ctx context.Context, objective, entry string) error {
	q := url.Values{"entry": {entry}}
	return m.http.doJSON(ctx, "DELETE", fmt.Sprintf("scoreboard/%s/score", objective), noBody(), q, nil)
}

// Broadcast sends a message to every player on the server.
func (m *ServerModule) Broadcast(ctx context.Context, message string) error {
	body := url.Values{"message": {message}}
	return m.http.doJSON(ctx, "POST", "chat/broadcast", formBody(body), nil, nil)
}

// Tell sends a private message to a specific player. Server-side reads the
// form field playerUuid, not uuid.
func (m *ServerModule) Tell(ctx context.Context, uuid, message string) error {
	body := url.Values{"playerUuid": {uuid}, "message": {message}}
	return m.http.doJSON(ctx, "POST", "chat/tell", formBody(body), nil, nil)
}

func limitParam(key string, n int) url.Values {
	if n <= 0 {
		return nil
	}
	return url.Values{key: {strconv.Itoa(n)}}
}

func stringParam(key, val string) url.Values {
	if val == "" {
		return nil
	}
	return url.Values{key: {val}}
}
