package noxaeapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// AuthMeModule provides access to AuthMe Reloaded integration endpoints.
type AuthMeModule struct{ http requester }

// GetPlayerAuth gets AuthMe player auth details.
func (m *AuthMeModule) GetPlayerAuth(ctx context.Context, name string) (AuthMePlayerInfo, error) {
	var out AuthMePlayerInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("authme/player/%s", url.PathEscape(name)), noBody(), nil, &out)
	return out, err
}

// CheckPassword checks if the password matches the player's AuthMe password.
func (m *AuthMeModule) CheckPassword(ctx context.Context, name, password string) (PasswordCheckResult, error) {
	var out PasswordCheckResult
	body := map[string]string{"password": password}
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("authme/player/%s/check-password", url.PathEscape(name)), jsonBody(body), nil, &out)
	return out, err
}

// FastLoginModule provides access to FastLogin premium status checks.
type FastLoginModule struct{ http requester }

// GetPlayerInfo checks if a player is running on a premium (Mojang-authenticated) account.
func (m *FastLoginModule) GetPlayerInfo(ctx context.Context, name string) (FastLoginPlayerInfo, error) {
	var out FastLoginPlayerInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("fastlogin/player/%s", url.PathEscape(name)), noBody(), nil, &out)
	return out, err
}

// NLoginModule provides access to nLogin integration endpoints.
type NLoginModule struct{ http requester }

// GetPlayerAuth gets nLogin authentication and registration status for a player.
func (m *NLoginModule) GetPlayerAuth(ctx context.Context, name string) (NLoginPlayerInfo, error) {
	var out NLoginPlayerInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("nlogin/player/%s", url.PathEscape(name)), noBody(), nil, &out)
	return out, err
}

// LoginSecurityModule provides access to LoginSecurity integration endpoints.
type LoginSecurityModule struct{ http requester }

// GetPlayerAuth gets LoginSecurity authentication and registration status for a player.
func (m *LoginSecurityModule) GetPlayerAuth(ctx context.Context, name string) (LoginSecurityPlayerInfo, error) {
	var out LoginSecurityPlayerInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("loginsecurity/player/%s", url.PathEscape(name)), noBody(), nil, &out)
	return out, err
}

// BetterTeamsModule provides access to BetterTeams teams and player team lookups.
type BetterTeamsModule struct{ http requester }

// GetPlayerTeam gets team summary for a player by UUID or player name.
func (m *BetterTeamsModule) GetPlayerTeam(ctx context.Context, uuidOrName string) (BetterTeamsSummary, error) {
	var out BetterTeamsSummary
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("betterteams/player/%s", url.PathEscape(uuidOrName)), noBody(), nil, &out)
	return out, err
}

// GetTeam gets full details of a team by team name, including member list.
func (m *BetterTeamsModule) GetTeam(ctx context.Context, name string) (BetterTeamsDetail, error) {
	var out BetterTeamsDetail
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("betterteams/team/%s", url.PathEscape(name)), noBody(), nil, &out)
	return out, err
}

// GetTeams gets a paginated list of all teams.
func (m *BetterTeamsModule) GetTeams(ctx context.Context, page, limit int) (BetterTeamsListResponse, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var out BetterTeamsListResponse
	err := m.http.doJSON(ctx, "GET", "betterteams/teams", noBody(), q, &out)
	return out, err
}

// EssentialsModule provides access to EssentialsX player homes, warps, and info.
type EssentialsModule struct{ http requester }

// GetPlayerHomes gets player's home coordinates and list from Essentials.
func (m *EssentialsModule) GetPlayerHomes(ctx context.Context, name string) (EssentialsPlayerHomesResponse, error) {
	var out EssentialsPlayerHomesResponse
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("essentials/player/%s/homes", url.PathEscape(name)), noBody(), nil, &out)
	return out, err
}

// GetWarps gets all configured server warps from Essentials.
func (m *EssentialsModule) GetWarps(ctx context.Context) (EssentialsWarpsResponse, error) {
	var out EssentialsWarpsResponse
	err := m.http.doJSON(ctx, "GET", "essentials/warps", noBody(), nil, &out)
	return out, err
}

// GetPlayerInfo gets player state info from Essentials.
func (m *EssentialsModule) GetPlayerInfo(ctx context.Context, name string) (EssentialsPlayerInfo, error) {
	var out EssentialsPlayerInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("essentials/player/%s/info", url.PathEscape(name)), noBody(), nil, &out)
	return out, err
}

// GriefPreventionModule provides access to GriefPrevention claims.
type GriefPreventionModule struct{ http requester }

// GetPlayerClaims gets all claims owned by a player by UUID.
func (m *GriefPreventionModule) GetPlayerClaims(ctx context.Context, uuid string) (GriefPreventionClaimsResponse, error) {
	var out GriefPreventionClaimsResponse
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("claims/player/%s", url.PathEscape(uuid)), noBody(), nil, &out)
	return out, err
}

// GetClaimAt gets claim details at a specific coordinate location.
func (m *GriefPreventionModule) GetClaimAt(ctx context.Context, world string, x, z float64) (*GriefPreventionClaim, error) {
	q := url.Values{
		"world": {world},
		"x":     {strconv.FormatFloat(x, 'f', -1, 64)},
		"z":     {strconv.FormatFloat(z, 'f', -1, 64)},
	}
	var out GriefPreventionClaim
	err := m.http.doJSON(ctx, "GET", "claims", noBody(), q, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// DiscordSRVModule provides access to DiscordSRV linked accounts and broadcasts.
type DiscordSRVModule struct{ http requester }

// GetLinked gets Discord user ID linked to a player UUID.
func (m *DiscordSRVModule) GetLinked(ctx context.Context, uuid string) (DiscordSRVLinkedResponse, error) {
	var out DiscordSRVLinkedResponse
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("discordsrv/linked/%s", url.PathEscape(uuid)), noBody(), nil, &out)
	return out, err
}

// Broadcast broadcasts a message to a Discord channel configured in DiscordSRV.
func (m *DiscordSRVModule) Broadcast(ctx context.Context, message, channel string) (DiscordSRVBroadcastResponse, error) {
	if channel == "" {
		channel = "global"
	}
	body := map[string]string{
		"message": message,
		"channel": channel,
	}
	var out DiscordSRVBroadcastResponse
	err := m.http.doJSON(ctx, "POST", "discordsrv/broadcast", jsonBody(body), nil, &out)
	return out, err
}

// PlayerPointsModule provides access to PlayerPoints balance and transactions.
type PlayerPointsModule struct{ http requester }

// GetBalance gets points balance for a player UUID.
func (m *PlayerPointsModule) GetBalance(ctx context.Context, uuid string) (PlayerPointsBalanceResponse, error) {
	var out PlayerPointsBalanceResponse
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("playerpoints/%s/balance", url.PathEscape(uuid)), noBody(), nil, &out)
	return out, err
}

// Pay gives points to a player UUID.
func (m *PlayerPointsModule) Pay(ctx context.Context, uuid string, amount int64) (PlayerPointsTransactionResult, error) {
	body := map[string]int64{"amount": amount}
	var out PlayerPointsTransactionResult
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("playerpoints/%s/pay", url.PathEscape(uuid)), jsonBody(body), nil, &out)
	return out, err
}

// Take deducts points from a player UUID.
func (m *PlayerPointsModule) Take(ctx context.Context, uuid string, amount int64) (PlayerPointsTransactionResult, error) {
	body := map[string]int64{"amount": amount}
	var out PlayerPointsTransactionResult
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("playerpoints/%s/take", url.PathEscape(uuid)), jsonBody(body), nil, &out)
	return out, err
}

// Set sets points balance for a player UUID to an exact amount.
func (m *PlayerPointsModule) Set(ctx context.Context, uuid string, amount int64) (PlayerPointsTransactionResult, error) {
	body := map[string]int64{"amount": amount}
	var out PlayerPointsTransactionResult
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("playerpoints/%s/set", url.PathEscape(uuid)), jsonBody(body), nil, &out)
	return out, err
}

// WorldGuardModule provides access to WorldGuard regions.
type WorldGuardModule struct{ http requester }

// GetAllRegions gets all WorldGuard regions across all loaded worlds.
func (m *WorldGuardModule) GetAllRegions(ctx context.Context) (map[string][]WorldGuardRegion, error) {
	var out map[string][]WorldGuardRegion
	err := m.http.doJSON(ctx, "GET", "worldguard/regions", noBody(), nil, &out)
	return out, err
}

// GetRegion gets a specific WorldGuard region in a world.
func (m *WorldGuardModule) GetRegion(ctx context.Context, world, id string) (WorldGuardRegion, error) {
	var out WorldGuardRegion
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("worldguard/regions/%s/%s", url.PathEscape(world), url.PathEscape(id)), noBody(), nil, &out)
	return out, err
}

// PlanModule provides access to Player Analytics (Plan) player summaries.
type PlanModule struct{ http requester }

// GetPlayerSummary gets analytics summary for a player from Player Analytics (Plan).
func (m *PlanModule) GetPlayerSummary(ctx context.Context, uuid string) (PlanPlayerSummary, error) {
	var out PlanPlayerSummary
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("plan/player/%s/summary", url.PathEscape(uuid)), noBody(), nil, &out)
	return out, err
}

// GeyserModule provides access to Geyser/Floodgate Bedrock status info.
type GeyserModule struct{ http requester }

// GetPlayerInfo checks Bedrock status, XUID, and device info for a player via Geyser/Floodgate.
func (m *GeyserModule) GetPlayerInfo(ctx context.Context, uuid string) (GeyserPlayerInfo, error) {
	var out GeyserPlayerInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("geyser/player/%s", url.PathEscape(uuid)), noBody(), nil, &out)
	return out, err
}

// MythicMobsModule provides access to MythicMobs mob listing and spawning.
type MythicMobsModule struct{ http requester }

// GetMobs gets all registered custom mob types from MythicMobs.
func (m *MythicMobsModule) GetMobs(ctx context.Context) (MythicMobsListResponse, error) {
	var out MythicMobsListResponse
	err := m.http.doJSON(ctx, "GET", "mythicmobs/mobs", noBody(), nil, &out)
	return out, err
}

// Spawn spawns a custom MythicMob at target location coordinates.
func (m *MythicMobsModule) Spawn(ctx context.Context, opts MythicMobSpawnOptions) (MythicMobSpawnResponse, error) {
	var out MythicMobSpawnResponse
	err := m.http.doJSON(ctx, "POST", "mythicmobs/spawn", jsonBody(opts), nil, &out)
	return out, err
}

// TabModule provides access to TAB nametag prefix and suffix.
type TabModule struct{ http requester }

// GetNametag gets custom nametag prefix/suffix for a player from TAB.
func (m *TabModule) GetNametag(ctx context.Context, uuid string) (TabNametagResponse, error) {
	var out TabNametagResponse
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("tab/player/%s/nametag", url.PathEscape(uuid)), noBody(), nil, &out)
	return out, err
}

// SetNametag sets custom nametag prefix/suffix for a player via TAB.
func (m *TabModule) SetNametag(ctx context.Context, uuid, prefix, suffix string) (TabNametagResponse, error) {
	body := map[string]string{
		"prefix": prefix,
		"suffix": suffix,
	}
	var out TabNametagResponse
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("tab/player/%s/nametag", url.PathEscape(uuid)), jsonBody(body), nil, &out)
	return out, err
}

// CoreProtectModule provides access to CoreProtect block and action audit lookups.
type CoreProtectModule struct{ http requester }

// Lookup performs block and action history lookup via CoreProtect.
func (m *CoreProtectModule) Lookup(ctx context.Context, opts CoreProtectLookupOptions) (CoreProtectLookupResponse, error) {
	q := url.Values{}
	if opts.Time > 0 {
		q.Set("time", strconv.FormatInt(opts.Time, 10))
	}
	if opts.User != "" {
		q.Set("user", opts.User)
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.FormatInt(opts.Limit, 10))
	}
	var out CoreProtectLookupResponse
	err := m.http.doJSON(ctx, "GET", "coreprotect/lookup", noBody(), q, &out)
	return out, err
}

// ShopGuiModule provides access to ShopGUIPlus catalogs.
type ShopGuiModule struct{ http requester }

// GetShops gets all configured shops from ShopGUIPlus.
func (m *ShopGuiModule) GetShops(ctx context.Context) (ShopGuiListResponse, error) {
	var out ShopGuiListResponse
	err := m.http.doJSON(ctx, "GET", "shopgui/shops", noBody(), nil, &out)
	return out, err
}

// GetShop gets details of a specific shop by shop ID from ShopGUIPlus.
func (m *ShopGuiModule) GetShop(ctx context.Context, id string) (ShopGuiShopDetail, error) {
	var out ShopGuiShopDetail
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("shopgui/shop/%s", url.PathEscape(id)), noBody(), nil, &out)
	return out, err
}

// RankForgeModule provides access to RankForge PvP ranking and tiers.
type RankForgeModule struct{ http requester }

// GetPlayerStats gets full PvP stats for a player (points, kills, deaths, KD, tier info).
func (m *RankForgeModule) GetPlayerStats(ctx context.Context, uuid string) (RankForgePlayerStats, error) {
	var out RankForgePlayerStats
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("rankforge/player/%s", url.PathEscape(uuid)), noBody(), nil, &out)
	return out, err
}

// GetPoints gets current RankForge points for a player.
func (m *RankForgeModule) GetPoints(ctx context.Context, uuid string) (RankForgePointsResponse, error) {
	var out RankForgePointsResponse
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("rankforge/player/%s/points", url.PathEscape(uuid)), noBody(), nil, &out)
	return out, err
}

// SetPoints sets RankForge points for a player to an exact amount.
func (m *RankForgeModule) SetPoints(ctx context.Context, uuid string, amount int64) (RankForgePointsTransactionResult, error) {
	body := map[string]int64{"amount": amount}
	var out RankForgePointsTransactionResult
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("rankforge/player/%s/points/set", url.PathEscape(uuid)), jsonBody(body), nil, &out)
	return out, err
}

// AddPoints adds RankForge points to a player.
func (m *RankForgeModule) AddPoints(ctx context.Context, uuid string, amount int64) (RankForgePointsTransactionResult, error) {
	body := map[string]int64{"amount": amount}
	var out RankForgePointsTransactionResult
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("rankforge/player/%s/points/add", url.PathEscape(uuid)), jsonBody(body), nil, &out)
	return out, err
}

// RemovePoints removes RankForge points from a player.
func (m *RankForgeModule) RemovePoints(ctx context.Context, uuid string, amount int64) (RankForgePointsTransactionResult, error) {
	body := map[string]int64{"amount": amount}
	var out RankForgePointsTransactionResult
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("rankforge/player/%s/points/remove", url.PathEscape(uuid)), jsonBody(body), nil, &out)
	return out, err
}

// GetTiers gets all registered RankForge tiers.
func (m *RankForgeModule) GetTiers(ctx context.Context) ([]RankForgeTier, error) {
	var out []RankForgeTier
	err := m.http.doJSON(ctx, "GET", "rankforge/tiers", noBody(), nil, &out)
	return out, err
}

// GetTop gets top players by RankForge points.
func (m *RankForgeModule) GetTop(ctx context.Context, limit int) ([]RankForgeTopEntry, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	var out []RankForgeTopEntry
	err := m.http.doJSON(ctx, "GET", "rankforge/top", noBody(), q, &out)
	return out, err
}

// Save triggers async save of all RankForge player data.
func (m *RankForgeModule) Save(ctx context.Context) (RankForgeSaveResponse, error) {
	var out RankForgeSaveResponse
	err := m.http.doJSON(ctx, "POST", "rankforge/save", noBody(), nil, &out)
	return out, err
}

// IntegrationsModule aggregates all external plugin integrations supported by NoxAeApi.
type IntegrationsModule struct {
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

// newIntegrationsModule creates a new IntegrationsModule.
func newIntegrationsModule(h requester) *IntegrationsModule {
	return &IntegrationsModule{
		AuthMe:          &AuthMeModule{h},
		FastLogin:       &FastLoginModule{h},
		NLogin:          &NLoginModule{h},
		LoginSecurity:   &LoginSecurityModule{h},
		BetterTeams:     &BetterTeamsModule{h},
		Essentials:      &EssentialsModule{h},
		GriefPrevention: &GriefPreventionModule{h},
		DiscordSRV:      &DiscordSRVModule{h},
		PlayerPoints:    &PlayerPointsModule{h},
		WorldGuard:      &WorldGuardModule{h},
		Plan:            &PlanModule{h},
		Geyser:          &GeyserModule{h},
		MythicMobs:      &MythicMobsModule{h},
		Tab:             &TabModule{h},
		CoreProtect:     &CoreProtectModule{h},
		ShopGui:         &ShopGuiModule{h},
		RankForge:       &RankForgeModule{h},
	}
}
