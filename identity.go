package noxaeapi

import (
	"context"
	"fmt"
)

// IdentityModule provides player identity resolution and network server discovery.
// Access via Client.Identity.
type IdentityModule struct{ http requester }

// Resolve resolves player name or UUID to identity information.
func (m *IdentityModule) Resolve(ctx context.Context, nameOrUUID string) (PlayerResolveResult, error) {
	var out PlayerResolveResult
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("identity/resolve/%s", nameOrUUID), noBody(), nil, &out)
	return out, err
}

// GetNetworkServers returns all network servers known to the identity service.
func (m *IdentityModule) GetNetworkServers(ctx context.Context) ([]IdentityServerInfo, error) {
	var out []IdentityServerInfo
	err := m.http.doJSON(ctx, "GET", "identity/servers", noBody(), nil, &out)
	return out, err
}

// GetNetworkPlayer returns complete identity information for a player across the network.
func (m *IdentityModule) GetNetworkPlayer(ctx context.Context, nameOrUUID string) (IdentityInfo, error) {
	var out IdentityInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("identity/player/%s", nameOrUUID), noBody(), nil, &out)
	return out, err
}
