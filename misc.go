package noxaeapi

import (
	"context"
	"net/url"
)

// AdvancementsModule lists advancements known to the server. Access via Client.Advancements.
type AdvancementsModule struct{ http *httpEngine }

// List returns all advancements known to the server.
func (m *AdvancementsModule) List(ctx context.Context) ([]Advancement, error) {
	var out []Advancement
	err := m.http.doJSON(ctx, "GET", "advancements", noBody(), nil, &out)
	return out, err
}

// PlaceholdersModule resolves PlaceholderAPI-style placeholders server-side
// for a given player. Access via Client.Placeholders.
type PlaceholdersModule struct{ http *httpEngine }

// Replace replaces placeholders (e.g. "%player_name%") in message for a
// player, returning the resolved string. Server-side reads form fields
// message and uuid.
func (m *PlaceholdersModule) Replace(ctx context.Context, uuid, message string) (string, error) {
	body := url.Values{"uuid": {uuid}, "message": {message}}
	var out string
	err := m.http.doJSON(ctx, "POST", "placeholders/replace", formBody(body), nil, &out)
	return out, err
}
