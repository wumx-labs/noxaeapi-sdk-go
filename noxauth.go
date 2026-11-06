package noxaeapi

import (
	"context"
	"fmt"
)

// NoxAuthModule wraps the /v1/noxauth/* routes. Access via Client.NoxAuth.
//
// Only works when noxauth.enabled: true is set in the server's
// noxaeapi-config.yml and the NoxAuth plugin is installed. CheckPassword's
// body is sent as real JSON.
type NoxAuthModule struct{ http requester }

// GetPlayerAuth returns NoxAuth registration/auth info for a player by name.
func (m *NoxAuthModule) GetPlayerAuth(ctx context.Context, name string) (NoxAuthPlayerInfo, error) {
	var out NoxAuthPlayerInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("noxauth/player/%s", name), noBody(), nil, &out)
	return out, err
}

// CheckPassword checks whether a password matches a player's stored NoxAuth password.
func (m *NoxAuthModule) CheckPassword(ctx context.Context, name, password string) (PasswordCheckResult, error) {
	body := map[string]interface{}{"password": password}
	var out PasswordCheckResult
	err := m.http.doJSON(ctx, "POST", fmt.Sprintf("noxauth/player/%s/check-password", name), jsonBody(body), nil, &out)
	return out, err
}
