package terraform

type Plan struct {
	ResourceChanges []TerraformResourceChange `json:"resource_changes"`
}

type TerraformResourceChange struct {
	Address string           `json:"address"`
	Mode    string           `json:"mode"`
	Type    string           `json:"type"`
	Name    string           `json:"name"`
	Change  TerraformChange  `json:"change"`
}

type TerraformChange struct {
	Actions []string `json:"actions"`
}

// Simplified internal model used by the rest of the CLI
type ResourceChange struct {
	Address string
	Type    string
	Name    string
	Action  string // normalized: create, update, delete, replace, unknown
}