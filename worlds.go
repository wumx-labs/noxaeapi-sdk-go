package noxaeapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// WorldsModule lists, saves, downloads, and configures worlds: time,
// weather, and entity counts. Access via Client.Worlds.
type WorldsModule struct{ http *httpEngine }

// List returns all worlds.
func (m *WorldsModule) List(ctx context.Context) ([]World, error) {
	var out []World
	err := m.http.doJSON(ctx, "GET", "worlds", noBody(), nil, &out)
	return out, err
}

// SaveAll saves all worlds to disk.
func (m *WorldsModule) SaveAll(ctx context.Context) error {
	return m.http.doJSON(ctx, "POST", "worlds/save", noBody(), nil, nil)
}

// DownloadAll returns a download link/stream reference for all worlds.
func (m *WorldsModule) DownloadAll(ctx context.Context) (map[string]string, error) {
	var out map[string]string
	err := m.http.doJSON(ctx, "GET", "worlds/download", noBody(), nil, &out)
	return out, err
}

// Get returns a single world by UUID.
func (m *WorldsModule) Get(ctx context.Context, uuid string) (World, error) {
	var out World
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("worlds/%s", uuid), noBody(), nil, &out)
	return out, err
}

// Save saves a specific world to disk.
func (m *WorldsModule) Save(ctx context.Context, uuid string) error {
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("worlds/%s/save", uuid), noBody(), nil, nil)
}

// Download returns a download link/stream reference for a specific world.
func (m *WorldsModule) Download(ctx context.Context, uuid string) (map[string]string, error) {
	var out map[string]string
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("worlds/%s/download", uuid), noBody(), nil, &out)
	return out, err
}

// SetTime sets the in-game time for a world (0-24000).
func (m *WorldsModule) SetTime(ctx context.Context, uuid string, ticks int64) error {
	body := url.Values{"time": {strconv.FormatInt(ticks, 10)}}
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("worlds/%s/time", uuid), formBody(body), nil, nil)
}

// SetWeather sets the weather for a world.
func (m *WorldsModule) SetWeather(ctx context.Context, uuid string, weather Weather) error {
	body := url.Values{"weather": {string(weather)}}
	return m.http.doJSON(ctx, "POST", fmt.Sprintf("worlds/%s/weather", uuid), formBody(body), nil, nil)
}

// GetEntities returns entity counts within a specific world.
func (m *WorldsModule) GetEntities(ctx context.Context, uuid string) (map[string]interface{}, error) {
	var out map[string]interface{}
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("worlds/%s/entities", uuid), noBody(), nil, &out)
	return out, err
}
