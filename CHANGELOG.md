# Changelog — noxaeapi-sdk-go (Go)

All notable changes to `github.com/NoxlyDev/noxaeapi-sdk-go` will be documented in this file.

---

## [0.7.0] - 2026-11-07

### Added
- **RankForge PvP Ranking & Tier Module (`RankForgeModule`)**: Full REST client methods: `GetPlayerStats`, `GetPoints`, `SetPoints`, `AddPoints`, `RemovePoints`, `GetTiers`, `GetTop`, and `Save`.
- **Integrations Suite (`IntegrationsModule`)**: Comprehensive client modules for 16 external plugins (`AuthMe`, `FastLogin`, `NLogin`, `LoginSecurity`, `BetterTeams`, `Essentials`, `GriefPrevention`, `DiscordSRV`, `PlayerPoints`, `WorldGuard`, `Plan`, `Geyser`, `MythicMobs`, `Tab`, `CoreProtect`, `ShopGui`) and `RankForge`.
- **Identity Module (`IdentityModule`)**: Cross-server identity resolution and server listing (`Resolve`, `GetNetworkServers`, `GetNetworkPlayer`).
- **Velocity Hub Node Forwarding (`HubNodeClient`)**: Added `hub.Node(id)` for transparent request forwarding through Velocity network hubs.

## [0.4.1] - 2026-10-20

### ✨ Features
* Complete idiomatic Go client for NoxAeApi REST & WebSocket endpoints.
* **Core Modules**: Players, Server, Economy (including ExcellentEconomy multi-currency), Worlds, Plugins, Advancements, Placeholders.
* **Permissions & Auth**: LuckPerms node management (`luckperms.go`), NoxAuth integration (`noxauth.go`).
* **Cross-Server**: Identity resolution (`identity.go`), Network polling aggregator (`network.go`), and Velocity Network Hub push client (`network_hub.go`).
* Built-in retry mechanism with exponential backoff and full jitter.
