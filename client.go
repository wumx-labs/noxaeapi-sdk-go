// Package noxaeapi provides a typed Go client for the NoxAeApi REST +
// WebSocket API (Fabric, Bukkit/Spigot/Paper) and the NoxAeApi-Velocity
// network hub. Zero required external dependencies — built entirely on the
// standard library. Realtime console/event socket support lives in the
// separate github.com/NoxlyDev/noxaeapi-sdk-go/ws subpackage, which has its
// own go.mod so importing it is opt-in.
package noxaeapi

import "os"

// Client is a client for a NoxAeApi server (Fabric mod or
// Bukkit/Spigot/Paper plugin).
//
//	client, err := noxaeapi.NewClient(noxaeapi.ClientOptions{
//	    BaseURL: "http://localhost:8080",
//	    APIKey:  "your-api-key",
//	})
//	players, err := client.Players.List(ctx)
//	err = client.Server.Broadcast(ctx, "Hello from the Go SDK!")
type Client struct {
	http *httpEngine

	Players      *PlayersModule
	Identity     *IdentityModule
	Economy      *EconomyModule
	Server       *ServerModule
	Worlds       *WorldsModule
	Plugins      *PluginsModule
	Advancements *AdvancementsModule
	Placeholders *PlaceholdersModule
	LuckPerms    *LuckPermsModule // requires LuckPerms loaded on the target server
	NoxAuth      *NoxAuthModule   // requires noxauth.enabled: true in the server config
	Skills       *SkillsModule    // requires mcMMO and/or AuraSkills loaded on the target server
	Leaderboards *LeaderboardModule
	Network      *NetworkModule // requires network.enabled: true in the server config

	Integrations    *IntegrationsModule
	AuthMe          *AuthMeModule
	FastLogin       *FastLoginModule
	NLogin          *NLoginModule
	LoginSecurity   *LoginSecurityModule
	BetterTeams     *BetterTeamsModule
	Essentials      *EssentialsModule
	GriefPrevention *GriefPreventionModule
	DiscordSRV      *DiscordSRVModule
	PlayerPoints    *PlayerPointsModule
	WorldGuard      *WorldGuardModule
	Plan            *PlanModule
	Geyser          *GeyserModule
	MythicMobs      *MythicMobsModule
	Tab             *TabModule
	CoreProtect     *CoreProtectModule
	ShopGui         *ShopGuiModule
	RankForge       *RankForgeModule
}

// NewClient creates a Client from options.
func NewClient(opts ClientOptions) (*Client, error) {
	h, err := newHTTPEngine(opts)
	if err != nil {
		return nil, err
	}
	integrations := newIntegrationsModule(h)
	return &Client{
		http:            h,
		Players:         &PlayersModule{h},
		Identity:        &IdentityModule{h},
		Economy:         &EconomyModule{h},
		Server:          &ServerModule{h},
		Worlds:          &WorldsModule{h},
		Plugins:         &PluginsModule{h},
		Advancements:    &AdvancementsModule{h},
		Placeholders:    &PlaceholdersModule{h},
		LuckPerms:       &LuckPermsModule{h},
		NoxAuth:         &NoxAuthModule{h},
		Skills:          &SkillsModule{h},
		Leaderboards:    &LeaderboardModule{h},
		Network:         &NetworkModule{h},
		Integrations:    integrations,
		AuthMe:          integrations.AuthMe,
		FastLogin:       integrations.FastLogin,
		NLogin:          integrations.NLogin,
		LoginSecurity:   integrations.LoginSecurity,
		BetterTeams:     integrations.BetterTeams,
		Essentials:      integrations.Essentials,
		GriefPrevention: integrations.GriefPrevention,
		DiscordSRV:      integrations.DiscordSRV,
		PlayerPoints:    integrations.PlayerPoints,
		WorldGuard:      integrations.WorldGuard,
		Plan:            integrations.Plan,
		Geyser:          integrations.Geyser,
		MythicMobs:      integrations.MythicMobs,
		Tab:             integrations.Tab,
		CoreProtect:     integrations.CoreProtect,
		ShopGui:         integrations.ShopGui,
		RankForge:       integrations.RankForge,
	}, nil
}

// NewClientFromEnv builds a Client from environment variables:
// NOXAEAPI_BASE_URL and NOXAEAPI_KEY.
func NewClientFromEnv() (*Client, error) {
	baseURL := os.Getenv("NOXAEAPI_BASE_URL")
	if baseURL == "" {
		return nil, &ConfigError{Message: "NewClientFromEnv: NOXAEAPI_BASE_URL is not set"}
	}
	return NewClient(ClientOptions{BaseURL: baseURL, APIKey: os.Getenv("NOXAEAPI_KEY")})
}

// BaseURL returns the configured base URL, mainly useful for constructing a
// WebSocket connection manually via the ws subpackage.
func (c *Client) BaseURL() string { return c.http.baseURL }

// APIKey returns the configured API key.
func (c *Client) APIKey() string { return c.http.apiKey }

// NetworkHubClient is a client for the **NoxAeApi-Velocity** network hub —
// a separate plugin that runs on the Velocity proxy, not on any individual
// backend server. Point BaseURL at the hub's own REST port, not a
// backend's port.
//
// Only exposes Network — the hub doesn't run any of the other REST modules
// (players, economy, worlds, ...) that a backend Client does. To reach a
// specific backend's own routes through the hub, use
// hub.Network.Forward(id, ...).
type NetworkHubClient struct {
	http    *httpEngine
	Network *NetworkHubModule
}

// NewNetworkHubClient creates a NetworkHubClient from options.
func NewNetworkHubClient(opts ClientOptions) (*NetworkHubClient, error) {
	h, err := newHTTPEngine(opts)
	if err != nil {
		return nil, err
	}
	return &NetworkHubClient{http: h, Network: &NetworkHubModule{h}}, nil
}

// NewNetworkHubClientFromEnv builds a hub client from environment
// variables: NOXAEAPI_HUB_BASE_URL / NOXAEAPI_HUB_KEY.
func NewNetworkHubClientFromEnv() (*NetworkHubClient, error) {
	baseURL := os.Getenv("NOXAEAPI_HUB_BASE_URL")
	if baseURL == "" {
		return nil, &ConfigError{Message: "NewNetworkHubClientFromEnv: NOXAEAPI_HUB_BASE_URL is not set"}
	}
	return NewNetworkHubClient(ClientOptions{BaseURL: baseURL, APIKey: os.Getenv("NOXAEAPI_HUB_KEY")})
}

// Node returns a HubNodeClient scoped to a specific backend server node ID,
// allowing all resource modules (Players, Economy, Integrations, etc.) to be
// called transparently forwarded through the hub proxy.
func (c *NetworkHubClient) Node(id string) *HubNodeClient {
	return newHubNodeClient(c.http, id)
}

// BaseURL returns the configured hub base URL.
func (c *NetworkHubClient) BaseURL() string { return c.http.baseURL }

// APIKey returns the configured hub API key.
func (c *NetworkHubClient) APIKey() string { return c.http.apiKey }
