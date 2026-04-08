package main

import (
	"strings"
	"testing"
)

func TestUpgradeOrder_Major_ReturnsZero(t *testing.T) {
	if upgradeOrder("major") != 0 {
		t.Errorf("got %d, want 0", upgradeOrder("major"))
	}
}

func TestUpgradeOrder_Minor_ReturnsOne(t *testing.T) {
	if upgradeOrder("minor") != 1 {
		t.Errorf("got %d, want 1", upgradeOrder("minor"))
	}
}

func TestUpgradeOrder_Patch_ReturnsTwo(t *testing.T) {
	if upgradeOrder("patch") != 2 {
		t.Errorf("got %d, want 2", upgradeOrder("patch"))
	}
}

func TestSortResults_SortsByUpgradeTypeThenName(t *testing.T) {
	input := []outdatedResult{
		{GroupID: "com.z", ArtifactID: "lib", Upgrade: "patch"},
		{GroupID: "com.a", ArtifactID: "lib", Upgrade: "major"},
		{GroupID: "com.m", ArtifactID: "lib", Upgrade: "minor"},
		{GroupID: "com.b", ArtifactID: "lib", Upgrade: "major"},
	}
	sortResults(input)

	want := []string{
		"major:com.a:lib",
		"major:com.b:lib",
		"minor:com.m:lib",
		"patch:com.z:lib",
	}
	for i, r := range input {
		got := r.Upgrade + ":" + r.GroupID + ":" + r.ArtifactID
		if got != want[i] {
			t.Errorf("index %d: got %q, want %q", i, got, want[i])
		}
	}
}

func TestSortResults_WithinTypeAlphabeticalByArtifact(t *testing.T) {
	input := []outdatedResult{
		{GroupID: "org", ArtifactID: "z-lib", Upgrade: "minor"},
		{GroupID: "org", ArtifactID: "a-lib", Upgrade: "minor"},
	}
	sortResults(input)

	if input[0].ArtifactID != "a-lib" {
		t.Errorf("got %q, want a-lib first", input[0].ArtifactID)
	}
}

func TestSortResults_Empty_DoesNotPanic(t *testing.T) {
	sortResults(nil)
	sortResults([]outdatedResult{})
}

func TestProgressLine_ContainsDoneAndTotal(t *testing.T) {
	got := progressLine(3, 10)
	if !strings.Contains(got, "3") || !strings.Contains(got, "10") {
		t.Errorf("progressLine(3, 10) = %q, should contain 3 and 10", got)
	}
}

func TestProgressLine_Format(t *testing.T) {
	got := progressLine(1, 50)
	want := "checking 1/50..."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
