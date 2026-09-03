package noxaeapi

import (
	"context"
	"fmt"
)

// LuckPermsModule wraps the /v1/luckperms/* routes. Access via Client.LuckPerms.
//
// These only exist on the server when the LuckPerms mod is loaded —
// calling any method here against a server without it will fail (typically
// a 404). Check Client.Plugins.List for LuckPerms if you need to branch on
// it ahead of time.
//
// Unlike most other modules, these POST/DELETE bodies are sent as real
// JSON (the server reads them with ctx.bodyAsClass(...)), not form-urlencoded.
type LuckPermsModule struct{ http *httpEngine }

// GetPlayerGroups returns the groups a player belongs to.
func (m *LuckPermsModule) GetPlayerGroups(ctx context.Context, uuid string) ([]string, error) {
	var out []string
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("luckperms/player/%s/groups", uuid), noBody(), nil, &out)
	return out, err
}

// GetPlayerPermissions returns a player's effective permission nodes.
func (m *LuckPermsModule) GetPlayerPermissions(ctx context.Context, uuid string) ([]PermissionNode, error) {
	var out []PermissionNode
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("luckperms/player/%s/permissions", uuid), noBody(), nil, &out)
	return out, err
}

// AddPlayerPermission adds a permission node to a player.
func (m *LuckPermsModule) AddPlayerPermission(ctx context.Context, uuid, permission string, value bool) error {
	body := map[string]interface{}{"permission": permission, "value": value}
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("luckperms/player/%s/permission", uuid), jsonBody(body), nil, nil)
}

// CheckPlayerPermission checks whether a player has a given permission.
func (m *LuckPermsModule) CheckPlayerPermission(ctx context.Context, uuid, permission string) (map[string]interface{}, error) {
	body := map[string]interface{}{"permission": permission}
	var out map[string]interface{}
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("luckperms/player/%s/check-permission", uuid), jsonBody(body), nil, &out)
	return out, err
}

// RemovePlayerPermission removes a permission node from a player.
func (m *LuckPermsModule) RemovePlayerPermission(ctx context.Context, uuid, permission string) error {
	body := map[string]interface{}{"permission": permission}
	return m.http.doJSON(ctx, "DELETE", fmt.Sprintf("luckperms/player/%s/permission", uuid), jsonBody(body), nil, nil)
}

// SetPlayerGroup sets a player's primary group.
func (m *LuckPermsModule) SetPlayerGroup(ctx context.Context, uuid, group string) error {
	body := map[string]interface{}{"group": group}
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("luckperms/player/%s/group", uuid), jsonBody(body), nil, nil)
}

// RemovePlayerGroup removes a group from a player.
func (m *LuckPermsModule) RemovePlayerGroup(ctx context.Context, uuid, groupName string) error {
	return m.http.doJSON(ctx, "DELETE", fmt.Sprintf("luckperms/player/%s/group/%s", uuid, groupName), noBody(), nil, nil)
}

// GetGroups lists all known groups.
func (m *LuckPermsModule) GetGroups(ctx context.Context) ([]string, error) {
	var out []string
	err := m.http.doJSON(ctx, "GET", "luckperms/groups", noBody(), nil, &out)
	return out, err
}

// CreateGroup creates a new LuckPerms group. The name is lowercased
// server-side. Returns an *APIError with a 409 status if it already exists.
func (m *LuckPermsModule) CreateGroup(ctx context.Context, name string) (map[string]string, error) {
	body := map[string]interface{}{"name": name}
	var out map[string]string
	err := m.http.doJSON(ctx, "POST", "luckperms/groups", jsonBody(body), nil, &out)
	return out, err
}

// DeleteGroup deletes a LuckPerms group by name. The "default" group can't be deleted.
func (m *LuckPermsModule) DeleteGroup(ctx context.Context, name string) error {
	return m.http.doJSON(ctx, "DELETE", fmt.Sprintf("luckperms/group/%s", name), noBody(), nil, nil)
}

// GetGroupPermissions returns the permissions attached to a specific group.
func (m *LuckPermsModule) GetGroupPermissions(ctx context.Context, name string) (GroupInfo, error) {
	var out GroupInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("luckperms/group/%s/permissions", name), noBody(), nil, &out)
	return out, err
}
