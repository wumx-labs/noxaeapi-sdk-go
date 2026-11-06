package noxaeapi

import "encoding/json"

type OnlinePlayer struct {
	UUID          string      `json:"uuid"`
	DisplayName   string      `json:"displayName"`
	Address       *string     `json:"address"`
	Port          *int        `json:"port"`
	Exhaustion    float64     `json:"exhaustion"`
	Exp           float64     `json:"exp"`
	ExpLevel      int64       `json:"expLevel"`
	Whitelisted   bool        `json:"whitelisted"`
	Banned        bool        `json:"banned"`
	Op            bool        `json:"op"`
	Balance       *float64    `json:"balance"`
	Location      *[3]float64 `json:"location"`
	Dimension     *string     `json:"dimension"`
	Health        float64     `json:"health"`
	Hunger        int64       `json:"hunger"`
	Saturation    float64     `json:"saturation"`
	Gamemode      string      `json:"gamemode"`
	LastPlayed    int64       `json:"lastPlayed"`
	Authenticated *bool       `json:"authenticated"`
	Registered    *bool       `json:"registered"`
}

type OfflinePlayer struct {
	UUID        string   `json:"uuid"`
	DisplayName string   `json:"displayName"`
	Whitelisted bool     `json:"whitelisted"`
	Banned      bool     `json:"banned"`
	Op          bool     `json:"op"`
	Balance     *float64 `json:"balance"`
	LastPlayed  int64    `json:"lastPlayed"`
}

// PlayerResolveResult is the result of resolving a player name to their
// UUID via PlayersModule.Resolve.
type PlayerResolveResult struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type ServerHealth struct {
	CPUs        int64 `json:"cpus"`
	Uptime      int64 `json:"uptime"`
	TotalMemory int64 `json:"totalMemory"`
	MaxMemory   int64 `json:"maxMemory"`
	FreeMemory  int64 `json:"freeMemory"`
}

type ServerBan struct {
	Target     string  `json:"target"`
	Source     *string `json:"source"`
	Reason     *string `json:"reason"`
	Expiration *string `json:"expiration"`
}

type WhitelistEntry struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type ServerInfo struct {
	Name               string           `json:"name"`
	Motd               string           `json:"motd"`
	Version            string           `json:"version"`
	BukkitVersion      string           `json:"bukkitVersion"`
	TPS                string           `json:"tps"`
	Health             ServerHealth     `json:"health"`
	BannedIPs          []ServerBan      `json:"bannedIps"`
	BannedPlayers      []ServerBan      `json:"bannedPlayers"`
	WhitelistedPlayers []WhitelistEntry `json:"whitelistedPlayers"`
	MaxPlayers         int64            `json:"maxPlayers"`
	OnlinePlayers      int64            `json:"onlinePlayers"`
}

type World struct {
	Name               string `json:"name"`
	UUID               string `json:"uuid"`
	Time               int64  `json:"time"`
	Storm              bool   `json:"storm"`
	Thundering         bool   `json:"thundering"`
	GenerateStructures bool   `json:"generateStructures"`
	AllowAnimals       bool   `json:"allowAnimals"`
	AllowMonsters      bool   `json:"allowMonsters"`
	Difficulty         string `json:"difficulty"`
	Environment        string `json:"environment"`
	Seed               string `json:"seed"`
}

type InventoryItem struct {
	ID    string `json:"id"`
	Count int64  `json:"count"`
	Slot  int64  `json:"slot"`
}

type Plugin struct {
	Name        string   `json:"name"`
	Enabled     bool     `json:"enabled"`
	Version     string   `json:"version"`
	Website     *string  `json:"website"`
	Authors     []string `json:"authors"`
	Depends     []string `json:"depends"`
	SoftDepends []string `json:"softDepends"`
	APIVersion  *string  `json:"apiVersion"`
	Description *string  `json:"description"`
}

// EconomyInfo describes the connected economy provider. Shape varies by
// provider (Vault vs Impactor); use Raw for provider-specific fields.
type EconomyInfo struct {
	Available bool            `json:"available"`
	Raw       json.RawMessage `json:"-"`
}

func (e *EconomyInfo) UnmarshalJSON(data []byte) error {
	type alias EconomyInfo
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	*e = EconomyInfo(a)
	e.Raw = data
	return nil
}

type PlayerBalance struct {
	UUID    string  `json:"uuid"`
	Balance float64 `json:"balance"`
}

type TopBalanceEntry struct {
	UUID    string  `json:"uuid"`
	Balance float64 `json:"balance"`
}

type GroupInfo struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

type PermissionNode struct {
	Permission string  `json:"permission"`
	Value      bool    `json:"value"`
	Expiry     int64   `json:"expiry"`
	Server     *string `json:"server"`
	World      *string `json:"world"`
}

type Advancement struct {
	Key      string   `json:"key"`
	Criteria []string `json:"criteria"`
}

type NoxAuthPlayerInfo struct {
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	Registered    bool    `json:"registered"`
	Authenticated bool    `json:"authenticated"`
	LastIP        *string `json:"lastIp"`
	LastLoginTime *int64  `json:"lastLoginTime"`
	CountryCode   *string `json:"countryCode"`
	CountryName   *string `json:"countryName"`
}

type PasswordCheckResult struct {
	Name  string `json:"name"`
	Valid bool   `json:"valid"`
}

type PlayerStats struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Kills        int64  `json:"kills"`
	Deaths       int64  `json:"deaths"`
	Playtime     int64  `json:"playtime"`
	BlocksPlaced int64  `json:"blocksPlaced"`
	BlocksBroken int64  `json:"blocksBroken"`
}

type SkillInfo struct {
	UUID       string           `json:"uuid"`
	Skills     map[string]int64 `json:"skills"`
	PowerLevel int64            `json:"powerLevel"`
}

type CurrencyBalance struct {
	UUID     string  `json:"uuid"`
	Name     *string `json:"name"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

type CurrencyTopEntry struct {
	UUID    string  `json:"uuid"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type LeaderboardSourceInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Available   bool   `json:"available"`
	// OnlineOnly is true if this source can only rank currently-online
	// players (e.g. mcMMO).
	OnlineOnly bool `json:"onlineOnly"`
}

type LeaderboardEntry struct {
	UUID  string  `json:"uuid"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// PlayerProfileBan is the ban status for a player, as reported by
// PlayerProfile.Status.Ban.
type PlayerProfileBan struct {
	Banned bool `json:"banned"`
	// Reason, Source, Created, and Expires are only present when Banned is
	// true. Created/Expires are serialized server-side as a locale-formatted
	// date string, not ISO-8601. A nil Expires with Banned true means a
	// permanent ban.
	Reason  *string `json:"reason,omitempty"`
	Source  *string `json:"source,omitempty"`
	Created *string `json:"created,omitempty"`
	Expires *string `json:"expires,omitempty"`
}

// PlayerProfileStatus is the whitelist + ban status for a player, as
// reported by PlayerProfile.Status.
type PlayerProfileStatus struct {
	Whitelisted bool             `json:"whitelisted"`
	Ban         PlayerProfileBan `json:"ban"`
}

// PlayerProfileEconomy is a player's Vault economy info, as reported by
// PlayerProfile.Economy. Single-currency (Vault), not ExcellentEconomy's
// multi-currency system — see EconomyModule.GetCurrencyBalance for that.
type PlayerProfileEconomy struct {
	Available bool `json:"available"`
	// Balance is only present when Available is true.
	Balance *float64 `json:"balance,omitempty"`
}

// PlayerProfileStatTile is one leaderboard source's contribution to a
// player profile's Stats slice.
type PlayerProfileStatTile struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Status string `json:"status"` // "ranked" | "not_ranked" | "unavailable"
	// Value is only present when Status is "ranked".
	Value *float64 `json:"value,omitempty"`
	// Rank is only present when Status is "ranked" and the source reports a positive rank.
	Rank *int64 `json:"rank,omitempty"`
}

// PlayerProfile is a one-call player profile combining identity,
// whitelist/ban status, Vault balance, and this player's entry in every
// registered leaderboard source. Returned by LeaderboardModule.GetPlayerProfile
// and (as the Players slice) by LeaderboardModule.GetPlayerProfiles.
type PlayerProfile struct {
	Name    *string                 `json:"name"`
	UUID    string                  `json:"uuid"`
	Online  bool                    `json:"online"`
	Status  PlayerProfileStatus     `json:"status"`
	Economy PlayerProfileEconomy    `json:"economy"`
	Stats   []PlayerProfileStatTile `json:"stats"`
}

// PlayerProfilesResponse is the response shape of
// LeaderboardModule.GetPlayerProfiles (bulk player profiles).
type PlayerProfilesResponse struct {
	Count   int64           `json:"count"`
	Players []PlayerProfile `json:"players"`
}

type NetworkServerStatus struct {
	ID      string         `json:"id"`
	Label   string         `json:"label"`
	Online  bool           `json:"online"`
	Server  *ServerInfo    `json:"server"`
	Players []OnlinePlayer `json:"players"`
}

type NetworkPlayersServerEntry struct {
	ID      string         `json:"id"`
	Label   string         `json:"label"`
	Online  bool           `json:"online"`
	Players []OnlinePlayer `json:"players"`
}

type NetworkPlayersResponse struct {
	Total   int64                       `json:"total"`
	Servers []NetworkPlayersServerEntry `json:"servers"`
}

type NetworkFindPlayerResponse struct {
	Found  bool          `json:"found"`
	Server string        `json:"server,omitempty"`
	Player *OnlinePlayer `json:"player,omitempty"`
}

type NetworkHealthServerEntry struct {
	ID     string          `json:"id"`
	Label  string          `json:"label"`
	Online bool            `json:"online"`
	TPS    json.RawMessage `json:"tps,omitempty"`
	Health json.RawMessage `json:"health,omitempty"`
}

type NetworkHealthResponse struct {
	Servers []NetworkHealthServerEntry `json:"servers"`
}

// NetworkBroadcastResponse maps network server ID -> "success" | "error".
type NetworkBroadcastResponse map[string]string

// NetworkHubNode is the last-known state of one backend node, as tracked by
// the NoxAeApi-Velocity network hub.
type NetworkHubNode struct {
	ID            string          `json:"id"`
	Label         string          `json:"label"`
	Online        bool            `json:"online"`
	TPS           string          `json:"tps"`
	OnlinePlayers int64           `json:"onlinePlayers"`
	MaxPlayers    int64           `json:"maxPlayers"`
	Health        json.RawMessage `json:"health"`
	// LastHeartbeatAt is the Unix epoch ms of the last heartbeat received, or 0 if never.
	LastHeartbeatAt int64 `json:"lastHeartbeatAt"`
}

type NetworkHubStatusResponse struct {
	Network []NetworkHubNode `json:"network"`
}

// NetworkHubPlayer is a player as seen directly by the proxy.
type NetworkHubPlayer struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Server string `json:"server,omitempty"`
}

type NetworkHubPlayersResponse struct {
	Total   int64              `json:"total"`
	Players []NetworkHubPlayer `json:"players"`
}

type NetworkHubFindPlayerResponse struct {
	Found  bool              `json:"found"`
	Player *NetworkHubPlayer `json:"player,omitempty"`
}

type NetworkHubBroadcastResponse struct {
	Delivered int64 `json:"delivered"`
}

// Gamemode is a player's active gamemode.
type Gamemode string

const (
	GamemodeSurvival  Gamemode = "survival"
	GamemodeCreative  Gamemode = "creative"
	GamemodeAdventure Gamemode = "adventure"
	GamemodeSpectator Gamemode = "spectator"
)

// Weather is a world's weather state.
type Weather string

const (
	WeatherClear   Weather = "clear"
	WeatherRain    Weather = "rain"
	WeatherThunder Weather = "thunder"
)

// IdentityServerInfo represents a backend server known to the identity service.
type IdentityServerInfo struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

// IdentityInfo is the result of resolving player identity across the network.
type IdentityInfo struct {
	Username        string             `json:"username"`
	IdentityUUID    string             `json:"identityUuid"`
	LocalUUID       *string            `json:"localUuid"`
	Server          IdentityServerInfo `json:"server"`
	OnlineMode      bool               `json:"onlineMode"`
	HasPlayedBefore bool               `json:"hasPlayedBefore"`
}

// AuthMePlayerInfo contains player authentication details from AuthMe Reloaded.
type AuthMePlayerInfo struct {
	Name          string `json:"name"`
	UUID          string `json:"uuid"`
	Registered    bool   `json:"registered"`
	Authenticated bool   `json:"authenticated"`
	LastIP        string `json:"lastIp"`
	LastLoginTime int64  `json:"lastLoginTime"`
}

// FastLoginPlayerInfo contains player authentication details from FastLogin.
type FastLoginPlayerInfo struct {
	Name    string `json:"name"`
	UUID    string `json:"uuid"`
	Premium bool   `json:"premium"`
}

// NLoginPlayerInfo contains player authentication details from nLogin.
type NLoginPlayerInfo struct {
	Name          string `json:"name"`
	UUID          string `json:"uuid"`
	Registered    bool   `json:"registered"`
	Authenticated bool   `json:"authenticated"`
}

// LoginSecurityPlayerInfo contains player authentication details from LoginSecurity.
type LoginSecurityPlayerInfo struct {
	Name          string `json:"name"`
	UUID          string `json:"uuid"`
	Registered    bool   `json:"registered"`
	Authenticated bool   `json:"authenticated"`
}

// BetterTeamsSummary contains summary info for a team in BetterTeams.
type BetterTeamsSummary struct {
	Name        string                 `json:"name"`
	Tag         string                 `json:"tag"`
	Description string                 `json:"description"`
	Score       float64                `json:"score"`
	Money       float64                `json:"money"`
	Level       int64                  `json:"level"`
	MemberCount int64                  `json:"memberCount"`
	Extra       map[string]interface{} `json:"-"`
}

// BetterTeamsDetail contains detailed info for a team including members list.
type BetterTeamsDetail struct {
	BetterTeamsSummary
	Members []string `json:"members"`
}

// BetterTeamsListResponse is the paginated response for teams listing.
type BetterTeamsListResponse struct {
	Page  int64                `json:"page"`
	Limit int64                `json:"limit"`
	Total int64                `json:"total"`
	Teams []BetterTeamsSummary `json:"teams"`
}

// EssentialsHomeLocation contains coordinate info for an Essentials home.
type EssentialsHomeLocation struct {
	Name  string   `json:"name"`
	World string   `json:"world"`
	X     float64  `json:"x"`
	Y     float64  `json:"y"`
	Z     float64  `json:"z"`
	Yaw   *float64 `json:"yaw,omitempty"`
	Pitch *float64 `json:"pitch,omitempty"`
}

// EssentialsPlayerHomesResponse contains homes for a player from Essentials.
type EssentialsPlayerHomesResponse struct {
	Player string                   `json:"player"`
	Homes  []EssentialsHomeLocation `json:"homes"`
}

// EssentialsWarpsResponse contains all configured server warps from Essentials.
type EssentialsWarpsResponse struct {
	Total int64    `json:"total"`
	Warps []string `json:"warps"`
}

// EssentialsPlayerInfo contains player state info from Essentials.
type EssentialsPlayerInfo struct {
	Player   string  `json:"player"`
	Nickname *string `json:"nickname,omitempty"`
	AFK      bool    `json:"afk"`
	Muted    bool    `json:"muted"`
	Jailed   bool    `json:"jailed"`
	GodMode  bool    `json:"godMode"`
}

// GriefPreventionBoundary is a 3D coordinate point.
type GriefPreventionBoundary struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// GriefPreventionClaim represents a land claim in GriefPrevention.
type GriefPreventionClaim struct {
	ID              interface{}              `json:"id,omitempty"`
	Owner           string                   `json:"owner,omitempty"`
	OwnerUUID       string                   `json:"ownerUuid,omitempty"`
	World           string                   `json:"world,omitempty"`
	GreaterBoundary *GriefPreventionBoundary `json:"greaterBoundary,omitempty"`
	LesserBoundary  *GriefPreventionBoundary `json:"lesserBoundary,omitempty"`
	Area            *float64                 `json:"area,omitempty"`
}

// GriefPreventionClaimsResponse contains land claims owned by a player.
type GriefPreventionClaimsResponse struct {
	UUID   string                 `json:"uuid"`
	Total  int64                  `json:"total"`
	Claims []GriefPreventionClaim `json:"claims"`
}

// DiscordSRVLinkedResponse contains Discord user ID linked to a player UUID.
type DiscordSRVLinkedResponse struct {
	UUID      string `json:"uuid"`
	DiscordID string `json:"discordId"`
}

// DiscordSRVBroadcastResponse contains the result of broadcasting via DiscordSRV.
type DiscordSRVBroadcastResponse struct {
	Success bool   `json:"success"`
	Channel string `json:"channel"`
	Message string `json:"message"`
}

// PlayerPointsBalanceResponse contains points balance for a player from PlayerPoints.
type PlayerPointsBalanceResponse struct {
	UUID   string `json:"uuid"`
	Points int64  `json:"points"`
}

// PlayerPointsTransactionResult contains transaction result for PlayerPoints operations.
type PlayerPointsTransactionResult struct {
	UUID           string `json:"uuid"`
	Amount         int64  `json:"amount"`
	Success        bool   `json:"success"`
	CurrentBalance int64  `json:"currentBalance"`
}

// WorldGuardRegionPoint is a 3D point for region bounds.
type WorldGuardRegionPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// WorldGuardRegion represents a region in WorldGuard.
type WorldGuardRegion struct {
	ID       string                 `json:"id"`
	Type     *string                `json:"type,omitempty"`
	Priority *int64                 `json:"priority,omitempty"`
	Owners   []string               `json:"owners,omitempty"`
	Members  []string               `json:"members,omitempty"`
	Flags    map[string]interface{} `json:"flags,omitempty"`
	Min      *WorldGuardRegionPoint `json:"min,omitempty"`
	Max      *WorldGuardRegionPoint `json:"max,omitempty"`
}

// PlanPlayerSummary represents player analytics summary from Plan.
type PlanPlayerSummary struct {
	UUID  string                 `json:"uuid"`
	Extra map[string]interface{} `json:"-"`
}

// GeyserPlayerInfo contains Bedrock info from Geyser/Floodgate.
type GeyserPlayerInfo struct {
	Bedrock  bool    `json:"bedrock"`
	XUID     *string `json:"xuid,omitempty"`
	DeviceOS *string `json:"deviceOs,omitempty"`
	Version  *string `json:"version,omitempty"`
}

// MythicMobsListResponse lists registered custom mobs from MythicMobs.
type MythicMobsListResponse struct {
	Total int64    `json:"total"`
	Mobs  []string `json:"mobs"`
}

// MythicMobSpawnOptions contains options for spawning a custom mob.
type MythicMobSpawnOptions struct {
	Type  string  `json:"type"`
	World string  `json:"world"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Z     float64 `json:"z"`
	Level *int64  `json:"level,omitempty"`
}

// MythicMobSpawnResponse contains the result of spawning a custom mob.
type MythicMobSpawnResponse struct {
	Success bool    `json:"success"`
	Type    string  `json:"type"`
	World   string  `json:"world"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Z       float64 `json:"z"`
	Level   int64   `json:"level"`
}

// TabNametagResponse contains TAB nametag prefix and suffix.
type TabNametagResponse struct {
	UUID    string `json:"uuid"`
	Prefix  string `json:"prefix"`
	Suffix  string `json:"suffix"`
	Success *bool  `json:"success,omitempty"`
}

// CoreProtectLookupOptions contains query parameters for CoreProtect lookup.
type CoreProtectLookupOptions struct {
	Time  int64  `json:"time,omitempty"`
	User  string `json:"user,omitempty"`
	Limit int64  `json:"limit,omitempty"`
}

// CoreProtectLookupEntry represents a single block/action history entry.
type CoreProtectLookupEntry struct {
	Time   *int64      `json:"time,omitempty"`
	User   *string     `json:"user,omitempty"`
	Action interface{} `json:"action,omitempty"`
	Type   *string     `json:"type,omitempty"`
	World  *string     `json:"world,omitempty"`
	X      *float64    `json:"x,omitempty"`
	Y      *float64    `json:"y,omitempty"`
	Z      *float64    `json:"z,omitempty"`
}

// CoreProtectLookupResponse contains the lookup results from CoreProtect.
type CoreProtectLookupResponse struct {
	Total             int64                    `json:"total"`
	TimeWindowSeconds int64                    `json:"timeWindowSeconds"`
	UserFilter        string                   `json:"userFilter"`
	Results           []CoreProtectLookupEntry `json:"results"`
}

// ShopGuiSummary represents basic info for a ShopGUIPlus shop.
type ShopGuiSummary struct {
	ID   string  `json:"id"`
	Name *string `json:"name,omitempty"`
}

// ShopGuiListResponse lists all shops from ShopGUIPlus.
type ShopGuiListResponse struct {
	Total int64            `json:"total"`
	Shops []ShopGuiSummary `json:"shops"`
}

// ShopGuiShopDetail represents full details of a specific shop.
type ShopGuiShopDetail struct {
	ID    string                 `json:"id"`
	Extra map[string]interface{} `json:"-"`
}

// RankForgeTier represents a PvP tier in RankForge.
type RankForgeTier struct {
	Name      string                 `json:"name"`
	MinPoints int64                  `json:"minPoints"`
	MaxPoints *int64                 `json:"maxPoints,omitempty"`
	Extra     map[string]interface{} `json:"-"`
}

// RankForgePlayerStats contains PvP ranking stats for a player.
type RankForgePlayerStats struct {
	UUID        string         `json:"uuid"`
	Points      int64          `json:"points"`
	Kills       int64          `json:"kills"`
	Deaths      int64          `json:"deaths"`
	KD          float64        `json:"kd"`
	DailyKills  int64          `json:"dailyKills"`
	CurrentTier *RankForgeTier `json:"currentTier"`
	NextTier    *RankForgeTier `json:"nextTier"`
}

// RankForgePointsResponse contains points for a player from RankForge.
type RankForgePointsResponse struct {
	UUID   string `json:"uuid"`
	Points int64  `json:"points"`
}

// RankForgePointsTransactionResult contains transaction result for RankForge points.
type RankForgePointsTransactionResult struct {
	UUID    string `json:"uuid"`
	Points  int64  `json:"points"`
	Success bool   `json:"success"`
}

// RankForgeTopEntry represents a leaderboard entry in RankForge.
type RankForgeTopEntry struct {
	Rank   int64          `json:"rank"`
	UUID   string         `json:"uuid"`
	Name   *string        `json:"name,omitempty"`
	Points int64          `json:"points"`
	Tier   *RankForgeTier `json:"tier,omitempty"`
}

// RankForgeTiersResponse contains all configured tiers from RankForge.
type RankForgeTiersResponse struct {
	Tiers []RankForgeTier `json:"tiers"`
}

// RankForgeSaveResponse contains the result of saving RankForge data.
type RankForgeSaveResponse struct {
	Success bool `json:"success"`
}
