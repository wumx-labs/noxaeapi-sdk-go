package noxaeapi

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// EconomyModule manages player currency balances via Vault (single
// currency) or ExcellentEconomy (multi-currency). Access via Client.Economy.
type EconomyModule struct{ http *httpEngine }

// Info returns info about the connected economy provider (Impactor on
// Fabric, Vault on Bukkit/Spigot/Paper).
func (m *EconomyModule) Info(ctx context.Context) (EconomyInfo, error) {
	var out EconomyInfo
	err := m.http.doJSON(ctx, "GET", "economy", noBody(), nil, &out)
	return out, err
}

// GetBalance returns a player's default (Vault) balance.
func (m *EconomyModule) GetBalance(ctx context.Context, uuid string) (PlayerBalance, error) {
	var out PlayerBalance
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("economy/balance/%s", uuid), noBody(), nil, &out)
	return out, err
}

// GetTopBalance returns the top richest players for the default currency.
// Pass limit <= 0 to use the server default.
func (m *EconomyModule) GetTopBalance(ctx context.Context, limit int) ([]TopBalanceEntry, error) {
	var out []TopBalanceEntry
	err := m.http.doJSON(ctx, "GET", "economy/top", noBody(), limitQuery(limit), &out)
	return out, err
}

// Pay credits/deposits an amount into a player's default balance.
func (m *EconomyModule) Pay(ctx context.Context, uuid string, amount float64) error {
	body := url.Values{"uuid": {uuid}, "amount": {formatFloat(amount)}}
	return m.http.doJSON(ctx, "POST", "economy/pay", formBody(body), nil, nil)
}

// Debit debits funds from a player's default balance.
func (m *EconomyModule) Debit(ctx context.Context, uuid string, amount float64) error {
	body := url.Values{"uuid": {uuid}, "amount": {formatFloat(amount)}}
	return m.http.doJSON(ctx, "POST", "economy/debit", formBody(body), nil, nil)
}

// --- ExcellentEconomy multi-currency (native API, not Vault) --------------
//
// These talk directly to ExcellentEconomy's Developer API rather than
// Vault, so they work with any currency configured on the server. They
// return an *APIError with a 424 status if ExcellentEconomy isn't
// installed on the target server, and 404 if the currency ID doesn't exist.

// ListCurrencies lists all currency IDs configured on ExcellentEconomy.
func (m *EconomyModule) ListCurrencies(ctx context.Context) ([]string, error) {
	var out []string
	err := m.http.doJSON(ctx, "GET", "economy/currencies", noBody(), nil, &out)
	return out, err
}

// GetCurrencyBalance returns a player's balance for a specific ExcellentEconomy currency.
func (m *EconomyModule) GetCurrencyBalance(ctx context.Context, currency, uuid string) (CurrencyBalance, error) {
	var out CurrencyBalance
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("economy/currency/%s/balance/%s", currency, uuid), noBody(), nil, &out)
	return out, err
}

// PayCurrency pays a player in a specific currency (adds amount). amount must be > 0.
func (m *EconomyModule) PayCurrency(ctx context.Context, currency, uuid string, amount float64) (string, error) {
	body := url.Values{"uuid": {uuid}, "amount": {formatFloat(amount)}}
	var out string
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("economy/currency/%s/pay", currency), formBody(body), nil, &out)
	return out, err
}

// DebitCurrency debits a player in a specific currency (subtracts amount). amount must be > 0.
func (m *EconomyModule) DebitCurrency(ctx context.Context, currency, uuid string, amount float64) (string, error) {
	body := url.Values{"uuid": {uuid}, "amount": {formatFloat(amount)}}
	var out string
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("economy/currency/%s/debit", currency), formBody(body), nil, &out)
	return out, err
}

// SetCurrencyBalance sets a player's balance for a specific currency to an exact amount (>= 0).
func (m *EconomyModule) SetCurrencyBalance(ctx context.Context, currency, uuid string, amount float64) (string, error) {
	body := url.Values{"uuid": {uuid}, "amount": {formatFloat(amount)}}
	var out string
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("economy/currency/%s/set", currency), formBody(body), nil, &out)
	return out, err
}

// GetCurrencyTop returns the top richest players for a specific currency.
func (m *EconomyModule) GetCurrencyTop(ctx context.Context, currency string, limit int) ([]CurrencyTopEntry, error) {
	var out []CurrencyTopEntry
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("economy/currency/%s/top", currency), noBody(), limitQuery(limit), &out)
	return out, err
}

func limitQuery(limit int) url.Values {
	if limit <= 0 {
		return nil
	}
	return url.Values{"limit": {strconv.Itoa(limit)}}
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
