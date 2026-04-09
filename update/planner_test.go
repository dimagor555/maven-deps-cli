package update

import (
	"testing"
)

func TestPlan_FilterByLevel_Patch(t *testing.T) {
	upgrades := []Upgrade{
		{Alias: "a", Section: "libraries", Current: "1.0.0", Latest: "1.0.1", Type: "patch", InlineVersion: true},
		{Alias: "b", Section: "libraries", Current: "1.0.0", Latest: "1.1.0", Type: "minor", InlineVersion: true},
		{Alias: "c", Section: "libraries", Current: "1.0.0", Latest: "2.0.0", Type: "major", InlineVersion: true},
	}
	plan := Plan(upgrades, LevelPatch)
	if len(plan.InlineUpdates) != 1 || plan.InlineUpdates[0].Alias != "a" {
		t.Errorf("got %+v", plan.InlineUpdates)
	}
}

func TestPlan_FilterByLevel_Minor(t *testing.T) {
	upgrades := []Upgrade{
		{Alias: "a", Section: "libraries", Current: "1.0.0", Latest: "1.0.1", Type: "patch", InlineVersion: true},
		{Alias: "b", Section: "libraries", Current: "1.0.0", Latest: "1.1.0", Type: "minor", InlineVersion: true},
		{Alias: "c", Section: "libraries", Current: "1.0.0", Latest: "2.0.0", Type: "major", InlineVersion: true},
	}
	plan := Plan(upgrades, LevelMinor)
	if len(plan.InlineUpdates) != 2 {
		t.Errorf("got %d updates, want 2", len(plan.InlineUpdates))
	}
}

func TestPlan_FilterByLevel_Major(t *testing.T) {
	upgrades := []Upgrade{
		{Alias: "a", Section: "libraries", Current: "1.0.0", Latest: "1.0.1", Type: "patch", InlineVersion: true},
		{Alias: "b", Section: "libraries", Current: "1.0.0", Latest: "1.1.0", Type: "minor", InlineVersion: true},
		{Alias: "c", Section: "libraries", Current: "1.0.0", Latest: "2.0.0", Type: "major", InlineVersion: true},
	}
	plan := Plan(upgrades, LevelMajor)
	if len(plan.InlineUpdates) != 3 {
		t.Errorf("got %d updates, want 3", len(plan.InlineUpdates))
	}
}

func TestPlan_SharedRef_MinOfLatest(t *testing.T) {
	upgrades := []Upgrade{
		{Alias: "ktor-core", VersionRef: "ktor", Current: "3.0.0", Latest: "3.2.0", Type: "minor"},
		{Alias: "ktor-client", VersionRef: "ktor", Current: "3.0.0", Latest: "3.1.0", Type: "minor"},
	}
	plan := Plan(upgrades, LevelMinor)
	if len(plan.VersionRefUpdates) != 1 {
		t.Fatalf("got %d ref updates, want 1", len(plan.VersionRefUpdates))
	}
	ru := plan.VersionRefUpdates[0]
	if ru.Ref != "ktor" || ru.NewVersion != "3.1.0" {
		t.Errorf("got %+v, want ref=ktor new=3.1.0", ru)
	}
	if len(ru.Aliases) != 2 {
		t.Errorf("aliases = %v, want 2", ru.Aliases)
	}
}

func TestPlan_SharedRef_OneFilteredOut_NoBump(t *testing.T) {
	upgrades := []Upgrade{
		{Alias: "a", VersionRef: "shared", Current: "1.0.0", Latest: "1.1.0", Type: "minor"},
		{Alias: "b", VersionRef: "shared", Current: "1.0.0", Latest: "2.0.0", Type: "major"},
	}
	plan := Plan(upgrades, LevelMinor)
	if len(plan.VersionRefUpdates) != 0 {
		t.Errorf("expected no ref updates, got %+v", plan.VersionRefUpdates)
	}
}

func TestPlan_EmptyInput_EmptyPlan(t *testing.T) {
	plan := Plan(nil, LevelMinor)
	if len(plan.VersionRefUpdates) != 0 || len(plan.InlineUpdates) != 0 {
		t.Errorf("expected empty, got %+v", plan)
	}
}

func TestPlan_NoneUpgrade_Ignored(t *testing.T) {
	upgrades := []Upgrade{
		{Alias: "a", Section: "libraries", InlineVersion: true, Current: "1.0.0", Latest: "1.0.0", Type: "none"},
	}
	plan := Plan(upgrades, LevelMajor)
	if len(plan.InlineUpdates) != 0 {
		t.Errorf("none-type should be dropped, got %+v", plan.InlineUpdates)
	}
}

func TestPlan_MixedRefAndInline(t *testing.T) {
	upgrades := []Upgrade{
		{Alias: "ktor", VersionRef: "ktor", Current: "3.0.0", Latest: "3.1.0", Type: "minor"},
		{Alias: "gson", Section: "libraries", InlineVersion: true, Current: "2.11.0", Latest: "2.12.0", Type: "minor"},
	}
	plan := Plan(upgrades, LevelMinor)
	if len(plan.VersionRefUpdates) != 1 {
		t.Errorf("want 1 ref update, got %d", len(plan.VersionRefUpdates))
	}
	if len(plan.InlineUpdates) != 1 {
		t.Errorf("want 1 inline update, got %d", len(plan.InlineUpdates))
	}
}
