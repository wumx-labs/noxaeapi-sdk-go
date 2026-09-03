package noxaeapi

import (
	"context"
	"fmt"
)

// SkillsModule wraps the /v1/skills/* routes (mcMMO / AuraSkills). Access via Client.Skills.
type SkillsModule struct{ http *httpEngine }

// GetMcMMOSkills returns mcMMO skill levels and power level for a player.
// Only works for online players — returns an *APIError with a 404 status if
// the player isn't currently connected, and 503 if mcMMO isn't loaded on
// the target server.
func (m *SkillsModule) GetMcMMOSkills(ctx context.Context, uuid string) (SkillInfo, error) {
	var out SkillInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("skills/mcmmo/player/%s", uuid), noBody(), nil, &out)
	return out, err
}

// GetAuraSkills returns AuraSkills skill levels and power level for a
// player. Works for offline players. Returns an *APIError with a 503
// status if AuraSkills isn't loaded on the target server.
func (m *SkillsModule) GetAuraSkills(ctx context.Context, uuid string) (SkillInfo, error) {
	var out SkillInfo
	err := m.http.doJSON(ctx, "GET", fmt.Sprintf("skills/auraskills/player/%s", uuid), noBody(), nil, &out)
	return out, err
}
