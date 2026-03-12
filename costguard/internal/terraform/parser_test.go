package terraform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeAction(t *testing.T) {
	cases := []struct {
		in   []string
		want string
	}{
		{[]string{"create"}, "create"},
		{[]string{"update"}, "update"},
		{[]string{"delete"}, "delete"},
		{[]string{"delete", "create"}, "replace"},
		{[]string{"create", "delete"}, "replace"},
		{[]string{}, "unknown"},
		{[]string{"FOO"}, "unknown"},
	}

	for _, c := range cases {
		if got := normalizeAction(c.in); got != c.want {
			t.Fatalf("normalizeAction(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestParsePlanFile(t *testing.T) {
	planPath := filepath.Join("..", "..", "examples", "plan.json")
	got, err := ParsePlanFile(planPath)
	if err != nil {
		t.Fatalf("ParsePlanFile returned error: %v", err)
	}

	wantActions := []string{"create", "create", "update", "replace", "delete", "unknown"}
	wantAddresses := []string{
		"aws_db_instance.prod_db",
		"aws_nat_gateway.main",
		"aws_ebs_volume.data",
		"aws_instance.worker",
		"aws_instance.old_worker",
		"aws_custom.unknown_action",
	}

	if len(got) != len(wantActions) {
		t.Fatalf("expected %d resource changes, got %d", len(wantActions), len(got))
	}

	for i := range wantActions {
		if got[i].Action != wantActions[i] {
			t.Fatalf("item %d: action = %q, want %q", i, got[i].Action, wantActions[i])
		}
		if got[i].Address != wantAddresses[i] {
			t.Fatalf("item %d: address = %q, want %q", i, got[i].Address, wantAddresses[i])
		}
	}
}

func TestParsePlanFile_FileNotFound(t *testing.T) {
	_, err := ParsePlanFile("does-not-exist.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParsePlanFile_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	fp := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(fp, []byte("{ this is not json"), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	_, err := ParsePlanFile(fp)
	if err == nil {
		t.Fatal("expected JSON decode error, got nil")
	}
}