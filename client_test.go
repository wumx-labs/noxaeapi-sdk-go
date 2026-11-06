package noxaeapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientRequiresBaseURL(t *testing.T) {
	_, err := NewClient(ClientOptions{})
	if err == nil {
		t.Fatal("expected error for empty BaseURL")
	}
}

func TestPlayersListOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/players" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Header.Get("key") != "test-key" {
			t.Errorf("missing/incorrect key header: %q", r.Header.Get("key"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]OnlinePlayer{{UUID: "abc", DisplayName: "Steve"}})
	}))
	defer srv.Close()

	client, err := NewClient(ClientOptions{BaseURL: srv.URL, APIKey: "test-key"})
	if err != nil {
		t.Fatal(err)
	}

	players, err := client.Players.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(players) != 1 || players[0].UUID != "abc" {
		t.Fatalf("unexpected players: %+v", players)
	}
}

func Test404ReturnsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"not found"}`))
	}))
	defer srv.Close()

	client, _ := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := client.Players.Get(context.Background(), "nope")
	if !IsNotFound(err) {
		t.Fatalf("expected IsNotFound, got %v", err)
	}
}

func Test401ReturnsAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	client, _ := NewClient(ClientOptions{BaseURL: srv.URL})
	_, err := client.Server.Info(context.Background())
	if !IsUnauthorized(err) {
		t.Fatalf("expected IsUnauthorized, got %v", err)
	}
}

func TestFormEncodingForPay(t *testing.T) {
	var gotContentType, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotContentType = r.Header.Get("Content-Type")
		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)
		gotBody = string(buf)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client, _ := NewClient(ClientOptions{BaseURL: srv.URL})
	if err := client.Economy.Pay(context.Background(), "uuid-1", 100); err != nil {
		t.Fatal(err)
	}
	if gotContentType != "application/x-www-form-urlencoded" {
		t.Errorf("unexpected content type: %s", gotContentType)
	}
	if gotBody != "amount=100&uuid=uuid-1" {
		t.Errorf("unexpected body: %s", gotBody)
	}
}

func TestRetryOn503(t *testing.T) {
	attempts := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer srv.Close()

	retry := DefaultRetryOptions()
	client, _ := NewClient(ClientOptions{BaseURL: srv.URL, Retry: &retry})
	_, err := client.Server.Ping(context.Background())
	if err != nil {
		t.Fatalf("expected success after retries, got %v", err)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestIdentityAndRankForge(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v1/identity/resolve/Steve":
			json.NewEncoder(w).Encode(PlayerResolveResult{Name: "Steve", UUID: "uuid-steve"})
		case "/v1/rankforge/player/uuid-steve":
			json.NewEncoder(w).Encode(RankForgePlayerStats{
				UUID:   "uuid-steve",
				Points: 1250,
				Kills:  50,
				Deaths: 10,
				KD:     5.0,
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	client, err := NewClient(ClientOptions{BaseURL: srv.URL})
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Identity.Resolve(context.Background(), "Steve")
	if err != nil {
		t.Fatalf("identity resolve failed: %v", err)
	}
	if res.UUID != "uuid-steve" {
		t.Errorf("unexpected uuid: %s", res.UUID)
	}

	stats, err := client.RankForge.GetPlayerStats(context.Background(), "uuid-steve")
	if err != nil {
		t.Fatalf("rankforge stats failed: %v", err)
	}
	if stats.Points != 1250 {
		t.Errorf("unexpected points: %d", stats.Points)
	}
}

func TestHubNodeClientForwarding(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/v1/network/survival/rankforge/player/uuid-1" {
			json.NewEncoder(w).Encode(RankForgePlayerStats{UUID: "uuid-1", Points: 300})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	hub, err := NewNetworkHubClient(ClientOptions{BaseURL: srv.URL})
	if err != nil {
		t.Fatal(err)
	}

	node := hub.Node("survival")
	stats, err := node.RankForge.GetPlayerStats(context.Background(), "uuid-1")
	if err != nil {
		t.Fatalf("node forward failed: %v", err)
	}
	if stats.Points != 300 {
		t.Errorf("unexpected points: %d", stats.Points)
	}
}
