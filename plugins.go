package noxaeapi

import (
	"context"
	"fmt"
	"net/url"
)

// PluginsModule lists installed plugins/mods, installs by URL, and toggles
// enable/disable state. Access via Client.Plugins.
type PluginsModule struct{ http *httpEngine }

// List returns all installed plugins/mods.
func (m *PluginsModule) List(ctx context.Context) ([]Plugin, error) {
	var out []Plugin
	err := m.http.doJSON(ctx, "GET", "plugins", noBody(), nil, &out)
	return out, err
}

// Install installs a plugin by downloading it from a direct URL. This is a
// privileged endpoint — requires a write-enabled API key. Server-side reads
// the form field downloadUrl, not source.
func (m *PluginsModule) Install(ctx context.Context, downloadURL string) error {
	body := url.Values{"downloadUrl": {downloadURL}}
	return m.http.doJSON(ctx, "POST", "plugins", formBody(body), nil, nil)
}

// Enable enables a plugin by name.
func (m *PluginsModule) Enable(ctx context.Context, name string) error {
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("plugins/%s/enable", name), noBody(), nil, nil)
}

// Disable disables a plugin by name.
func (m *PluginsModule) Disable(ctx context.Context, name string) error {
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("plugins/%s/disable", name), noBody(), nil, nil)
}
