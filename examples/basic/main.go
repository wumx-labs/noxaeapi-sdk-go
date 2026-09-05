package main

import (
	"context"
	"fmt"
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
	fmt.Printf("%d players online\n", len(players))

	if len(players) > 0 {
		balance, err := client.Economy.GetBalance(ctx, players[0].UUID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s has $%.2f\n", players[0].DisplayName, balance.Balance)
	}

	if err := client.Server.Broadcast(ctx, "Hello from the Go SDK!"); err != nil {
		log.Fatal(err)
	}

	// Resolve a name to a UUID, then get their one-call profile (identity,
	// ban/whitelist status, Vault balance, and every leaderboard entry at once).
	resolved, err := client.Players.Resolve(ctx, "NoxlyDev")
	if err != nil {
		log.Fatal(err)
	}
	profile, err := client.Leaderboards.GetPlayerProfile(ctx, resolved.UUID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("banned: %v, whitelisted: %v\n", profile.Status.Ban.Banned, profile.Status.Whitelisted)
}
