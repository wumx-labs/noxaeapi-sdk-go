package noxaeapi

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// forwardingEngine adapts requests to forward via the network hub.
type forwardingEngine struct {
	http   *httpEngine
	nodeID string
}

func (f *forwardingEngine) doJSON(ctx context.Context, method, path string, body requestBody, query url.Values, out interface{}) error {
	fullPath := fmt.Sprintf("network/%s/%s", url.PathEscape(f.nodeID), strings.TrimLeft(path, "/"))
	return f.http.doJSON(ctx, method, fullPath, body, query, out)
}

func (f *forwardingEngine) doRaw(ctx context.Context, method, path string, body requestBody, query url.Values) ([]byte, error) {
	fullPath := fmt.Sprintf("network/%s/%s", url.PathEscape(f.nodeID), strings.TrimLeft(path, "/"))
	return f.http.doRaw(ctx, method, fullPath, body, query)
}

// HubNodeClient is a per-node view of a NetworkHubClient, returned by its
// Node(id) method. Exposes the exact same resource modules as a direct
// Client (Players, Economy, Server, Integrations, ...), but every call is
// routed through this hub's network/forward endpoint.
type HubNodeClient struct {
	NodeID string

	Players      *PlayersModule
	Identity     *IdentityModule
	Economy      *EconomyModule
	Server       *ServerModule
	Worlds       *WorldsModule
	Plugins      *PluginsModule
	Advancements *AdvancementsModule
	Placeholders *PlaceholdersModule
	LuckPerms    *LuckPermsModule
	NoxAuth      *NoxAuthModule
	Skills       *SkillsModule
	Leaderboards *LeaderboardModule
	Network      *NetworkModule

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

func newHubNodeClient(http *httpEngine, nodeID string) *HubNodeClient {
	fwd := &forwardingEngine{http: http, nodeID: nodeID}
	integrations := newIntegrationsModule(fwd)
	return &HubNodeClient{
		NodeID:          nodeID,
		Players:         &PlayersModule{fwd},
		Identity:        &IdentityModule{fwd},
		Economy:         &EconomyModule{fwd},
		Server:          &ServerModule{fwd},
		Worlds:          &WorldsModule{fwd},
		Plugins:         &PluginsModule{fwd},
		Advancements:    &AdvancementsModule{fwd},
		Placeholders:    &PlaceholdersModule{fwd},
		LuckPerms:       &LuckPermsModule{fwd},
		NoxAuth:         &NoxAuthModule{fwd},
		Skills:          &SkillsModule{fwd},
		Leaderboards:    &LeaderboardModule{fwd},
		Network:         &NetworkModule{fwd},
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
	}
}
