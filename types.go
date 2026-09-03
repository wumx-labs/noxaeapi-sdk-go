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
