package noxaeapi

import (
	"context"
	"fmt"
	"net/url"
)

// NetworkModule wraps the /v1/network/* routes. Access via Client.Network.
//
// Only exists when network.enabled: true is set with at least one backend
// server configured — calling any method here against a server without it
// will 404. This is NoxAeApi-main's built-in polling aggregator: it lives
// on the same backend server you're connected to and fans requests out to
// the other backends listed in that server's own config.
type NetworkModule struct{ http *httpEngine }

// StatusAll returns status (server info + online players) for every
// configured network server, keyed by server ID.
func (m *NetworkModule) StatusAll(ctx context.Context) (map[string]NetworkServerStatus, error) {
	var out map[string]NetworkServerStatus
	err := m.http.doJSON(ctx, "GET", "network/status", noBody(), nil, &out)
	return out, err
}

// StatusByID returns status for a single network server by its configured ID.
func (m *NetworkModule) StatusByID(ctx context.Context, id string) (NetworkServerStatus, error) {
	var out NetworkServerStatus
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("network/status/%s", id), noBody(), nil, &out)
	return out, err
}

// Players aggregates online players across every network server.
func (m *NetworkModule) Players(ctx context.Context) (NetworkPlayersResponse, error) {
	var out NetworkPlayersResponse
	err := m.http.doJSON(ctx, "GET", "network/players", noBody(), nil, &out)
	return out, err
}

// FindPlayer finds which network server a player is currently on by UUID.
func (m *NetworkModule) FindPlayer(ctx context.Context, uuid string) (NetworkFindPlayerResponse, error) {
	var out NetworkFindPlayerResponse
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("network/players/%s", uuid), noBody(), nil, &out)
	return out, err
}

// Health returns aggregate health (TPS/memory) from every network server.
func (m *NetworkModule) Health(ctx context.Context) (NetworkHealthResponse, error) {
	var out NetworkHealthResponse
	err := m.http.doJSON(ctx, "GET", "network/health", noBody(), nil, &out)
	return out, err
}

// Broadcast broadcasts a message to every server on the network. Returns a
// map of server ID -> "success" | "error" (per-server delivery result).
func (m *NetworkModule) Broadcast(ctx context.Context, message string) (NetworkBroadcastResponse, error) {
	body := url.Values{"message": {message}}
	var out NetworkBroadcastResponse
	err := m.http.doJSON(ctx, "POST", "network/broadcast", formBody(body), nil, &out)
	return out, err
}

// Forward forwards an arbitrary request to a specific network server's own
// REST API. E.g. Forward(ctx, "survival", "POST", "server/exec",
// url.Values{"command": {"say hi"}}, nil) reaches that server's
// POST /v1/server/exec directly. Pass a nil formValues for GET/no-body
// requests. Useful for endpoints the aggregator doesn't have a dedicated
// method for.
func (m *NetworkModule) Forward(ctx context.Context, id, method, path string, formValues url.Values, query url.Values) ([]byte, error) {
	fullPath := fmt.Sprintf("network/%s/%s", id, trimLeftSlash(path))
	body := noBody()
	if formValues != nil {
		body = formBody(formValues)
	}
	return m.http.doRaw(ctx, method, fullPath, body, query)
}

func trimLeftSlash(s string) string {
	for len(s) > 0 && s[0] == '/' {
		s = s[1:]
	}
	return s
}
