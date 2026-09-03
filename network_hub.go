package noxaeapi

import (
	"context"
	"fmt"
	"net/url"
)

// NetworkHubModule wraps the /v1/network/* routes exposed by the
// **NoxAeApi-Velocity** network hub — a separate plugin/process run on the
// Velocity proxy, listening on its own port (distinct from any backend's
// own REST port). Access via NetworkHubClient.Network.
//
// Backend Paper/Bukkit servers connect out to the hub over WebSocket and
// push register/heartbeat/player-join/player-quit events; the hub answers
// every method below from its own in-memory registry, so calls here are
// cheap and don't block on a live round trip to each backend the way
// NetworkModule (NoxAeApi-main's polling aggregator) does.
//
// Response shapes differ from NetworkModule even where route names match
// — e.g. Players returns one flat proxy-wide list, not a per-server
// breakdown. There's also no hub equivalent of /v1/network/health; per-node
// health is embedded in NetworkHubNode.Health instead.
type NetworkHubModule struct{ http *httpEngine }

// StatusAll returns last-known status for every backend node that has ever registered.
func (m *NetworkHubModule) StatusAll(ctx context.Context) (NetworkHubStatusResponse, error) {
	var out NetworkHubStatusResponse
	err := m.http.doJSON(ctx, "GET", "network/status", noBody(), nil, &out)
	return out, err
}

// StatusByID returns last-known status for a single node by its configured ID.
func (m *NetworkHubModule) StatusByID(ctx context.Context, id string) (NetworkHubNode, error) {
	var out NetworkHubNode
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("network/status/%s", id), noBody(), nil, &out)
	return out, err
}

// Players lists every player currently connected to the proxy, read
// straight from Velocity's own player registry.
func (m *NetworkHubModule) Players(ctx context.Context) (NetworkHubPlayersResponse, error) {
	var out NetworkHubPlayersResponse
	err := m.http.doJSON(ctx, "GET", "network/players", noBody(), nil, &out)
	return out, err
}

// FindPlayer finds which backend server a player is currently on by UUID (proxy-authoritative).
func (m *NetworkHubModule) FindPlayer(ctx context.Context, uuid string) (NetworkHubFindPlayerResponse, error) {
	var out NetworkHubFindPlayerResponse
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("network/players/%s", uuid), noBody(), nil, &out)
	return out, err
}

// Broadcast broadcasts a message directly to every player connected to the
// proxy. Unlike NetworkModule.Broadcast, this doesn't forward to each
// backend's /v1/chat/broadcast — the proxy already has every player in hand.
func (m *NetworkHubModule) Broadcast(ctx context.Context, message string) (NetworkHubBroadcastResponse, error) {
	body := url.Values{"message": {message}}
	var out NetworkHubBroadcastResponse
	err := m.http.doJSON(ctx, "POST", "network/broadcast", formBody(body), nil, &out)
	return out, err
}

// Forward forwards an arbitrary request to a specific backend node's own
// REST API. E.g. Forward(ctx, "survival", "POST", "server/exec",
// url.Values{"command": {"say hi"}}, nil) reaches that backend's
// POST /v1/server/exec directly. Pass a nil formValues for GET/no-body requests.
func (m *NetworkHubModule) Forward(ctx context.Context, id, method, path string, formValues url.Values, query url.Values) ([]byte, error) {
	fullPath := fmt.Sprintf("network/%s/%s", id, trimLeftSlash(path))
	body := noBody()
	if formValues != nil {
		body = formBody(formValues)
	}
	return m.http.doRaw(ctx, method, fullPath, body, query)
}
