# noxaeapi-sdk-go

Typed Go SDK for NoxAeApi, a REST + WebSocket API plugin/mod for Minecraft servers — ships for both Fabric and Bukkit/Spigot/PaperMC. The REST surface is identical across platforms.

Zero required dependencies in the core module — built entirely on the standard library (`net/http`). Realtime console/event socket support lives in a separate `ws` submodule so the `gorilla/websocket` dependency stays opt-in.

📖 **Docs:** [noxapi.noxlydev.xyz](https://noxapi.noxlydev.xyz) · [pkg.go.dev/github.com/NoxlyDev/noxaeapi-sdk-go](https://pkg.go.dev/github.com/NoxlyDev/noxaeapi-sdk-go)

## Install

```bash
go get github.com/NoxlyDev/noxaeapi-sdk-go

# with realtime console/event socket support:
go get github.com/NoxlyDev/noxaeapi-sdk-go/ws
```

## Usage

```go
package main

import (
	"context"
	"log"

	"github.com/NoxlyDev/noxaeapi-sdk-go"
)

func main() {
	client, err := noxaeapi.NewClient(noxaeapi.ClientOptions{
		BaseURL: "http://localhost:8080",
		APIKey:  "your-api-key",
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	players, err := client.Players.List(ctx)
	if err != nil {
		log.Fatal(err)
	}

	balance, err := client.Economy.GetBalance(ctx, players[0].UUID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s has $%.2f", players[0].DisplayName, balance.Balance)
	client.Server.Broadcast(ctx, "Hello from the Go SDK!")
}
```

### From environment variables

```go
// Reads NOXAEAPI_BASE_URL and NOXAEAPI_KEY from the environment.
client, err := noxaeapi.NewClientFromEnv()
```

### Realtime (console tail / events)

Requires the `ws` submodule (`go get github.com/NoxlyDev/noxaeapi-sdk-go/ws`).

```go
import "github.com/NoxlyDev/noxaeapi-sdk-go/ws"

socket, err := ws.Connect(ws.Options{
	BaseURL: client.BaseURL(),
	APIKey:  client.APIKey(),
	Route:   "console",
})
if err != nil {
	log.Fatal(err)
}
defer socket.Close()

for event := range socket.Events {
	switch event.Type {
	case ws.EventConsole:
		log.Println("[console]", string(event.Payload))
	case ws.EventClose:
		log.Println("disconnected")
	}
}
```

The socket runs on a background goroutine and reconnects automatically with exponential backoff on unexpected disconnects.

## Error handling

Every fallible call returns a Go `error`. Non-2xx HTTP responses come back as `*noxaeapi.APIError`; use the `Is*` helper functions to check the reason:

```go
_, err := client.Server.Restart(ctx)
switch {
case noxaeapi.IsForbidden(err):
	log.Println("this API key isn't allowed to restart the server")
case noxaeapi.IsNotFound(err):
	log.Println("not found")
case noxaeapi.IsRateLimited(err):
	log.Println("rate limited")
case noxaeapi.IsServerError(err):
	log.Println("server-side error (5xx), safe to retry")
case noxaeapi.IsNetworkError(err):
	log.Println("couldn't reach the server at all")
case err != nil:
	log.Printf("unexpected error: %v", err)
}
```

`*APIError` also exposes `.Status`, `.Method`, `.Path`, and the raw `.Body` for anything not covered by a helper.

## Request encoding

The server is a Javalin app; most endpoints read their body with `ctx.formParam(...)` (`application/x-www-form-urlencoded`) rather than JSON. The SDK follows the same split internally — `LuckPerms` and `NoxAuth.CheckPassword` are the only modules that send real JSON bodies (`ctx.bodyAsClass(...)`). You don't need to think about this as a consumer; it's handled per-method.

## Optional modules

Some modules only work depending on the target server's setup:

- `client.LuckPerms` — requires the LuckPerms mod to be loaded on the server
- `client.NoxAuth` — requires `noxauth.enabled: true` in the server's `noxaeapi-config.yml`
- `client.Economy.GetCurrencyBalance` / `.PayCurrency` / `.DebitCurrency` / `.SetCurrencyBalance` / `.GetCurrencyTop` / `.ListCurrencies` — requires ExcellentEconomy (returns an `*APIError` with a 424 status if it isn't installed)
- `client.Network` — requires `network.enabled: true` with at least one backend server configured

Calling these against a server without the corresponding feature enabled will fail (typically 404).

## NoxAeApi-Velocity network hub

If your network runs the **NoxAeApi-Velocity** proxy plugin, use `NetworkHubClient` instead of (or alongside) `Client` — point it at the hub's own REST port, not a backend server's port.

```go
hub, err := noxaeapi.NewNetworkHubClient(noxaeapi.ClientOptions{
	BaseURL: "http://localhost:9090", // the hub's api-port, not a backend's port
	APIKey:  "your-hub-key",
})

status, err := hub.Network.StatusAll(ctx)
players, err := hub.Network.Players(ctx)
hub.Network.Broadcast(ctx, "Hello from the hub!")

// Reach a specific backend's own REST routes through the hub:
hub.Network.Forward(ctx, "survival", "POST", "server/exec", url.Values{"command": {"say hi"}}, nil)
```

Or from environment variables (`NOXAEAPI_HUB_BASE_URL` / `NOXAEAPI_HUB_KEY`):

```go
hub, err := noxaeapi.NewNetworkHubClientFromEnv()
```

## Configuration

```go
retry := noxaeapi.DefaultRetryOptions() // Attempts: 3, BaseDelay: 300ms, MaxDelay: 5s

client, err := noxaeapi.NewClient(noxaeapi.ClientOptions{
	BaseURL: "https://mc.example.com",
	APIKey:  "...",
	Timeout: 10 * time.Second,
	Retry:   &retry, // nil disables retries
	Headers: map[string]string{"X-Extra": "..."},
})
```

## Repo layout

This is a two-module repo:

- `/` — the core `noxaeapi-sdk-go` module (zero dependencies)
- `/ws` — the optional `noxaeapi-sdk-go/ws` module (depends on `gorilla/websocket`)

A `go.work` file ties them together for local development; it has no effect on consumers.

## License

MIT
