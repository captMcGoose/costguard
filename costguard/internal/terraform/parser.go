package terraform

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ParsePlanFile reads a terraform plan JSON file and returns normalized ResourceChange items.
// It uses local raw types to avoid redeclaring types that live in models.go.
func ParsePlanFile(path string) ([]ResourceChange, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open plan file %q: %w", path, err)
	}
	defer f.Close()

	// local types matching only the JSON shape we need
	type rawChange struct {
		Address string `json:"address"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Mode    string `json:"mode"`
		Change  struct {
			Actions []string `json:"actions"`
		} `json:"change"`
	}

	var raw struct {
		ResourceChanges []rawChange `json:"resource_changes"`
	}

	dec := json.NewDecoder(f)
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode JSON plan %q: %w", path, err)
	}

	out := make([]ResourceChange, 0, len(raw.ResourceChanges))
	for i, rc := range raw.ResourceChanges {
		// Validate required fields
		if strings.TrimSpace(rc.Address) == "" {
			return nil, fmt.Errorf("resource_changes[%d] missing address", i)
		}
		if strings.TrimSpace(rc.Type) == "" {
			return nil, fmt.Errorf("resource_changes[%d] missing type", i)
		}
		if strings.TrimSpace(rc.Name) == "" {
			return nil, fmt.Errorf("resource_changes[%d] missing name", i)
		}

		action := normalizeAction(rc.Change.Actions)
		out = append(out, ResourceChange{
			Address: rc.Address,
			Type:    rc.Type,
			Name:    rc.Name,
			Action:  action,
		})
	}

	return out, nil
}

func normalizeAction(actions []string) string {
	if len(actions) == 0 {
		return "unknown"
	}

	seen := map[string]bool{}
	for _, a := range actions {
		seen[strings.ToLower(strings.TrimSpace(a))] = true
	}

	switch {
	case len(seen) == 1 && seen["create"]:
		return "create"
	case len(seen) == 1 && seen["update"]:
		return "update"
	case len(seen) == 1 && seen["delete"]:
		return "delete"
	case seen["delete"] && seen["create"]:
		return "replace"
	default:
		return "unknown"
	}
}