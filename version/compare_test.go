package version

import "testing"

func TestGetUpgradeType_WhenMajor_ReturnsMajor(t *testing.T) {
	if got := GetUpgradeType("1.0.0", "2.0.0"); got != Major {
		t.Errorf("got %q, want %q", got, Major)
	}
}

func TestGetUpgradeType_WhenMinor_ReturnsMinor(t *testing.T) {
	if got := GetUpgradeType("1.0.0", "1.1.0"); got != Minor {
		t.Errorf("got %q, want %q", got, Minor)
	}
}

func TestGetUpgradeType_WhenPatch_ReturnsPatch(t *testing.T) {
	if got := GetUpgradeType("1.0.0", "1.0.1"); got != Patch {
		t.Errorf("got %q, want %q", got, Patch)
	}
}

func TestGetUpgradeType_WhenSame_ReturnsNone(t *testing.T) {
	if got := GetUpgradeType("1.0.0", "1.0.0"); got != None {
		t.Errorf("got %q, want %q", got, None)
	}
}

func TestGetUpgradeType_WhenDowngrade_ReturnsNone(t *testing.T) {
	cases := []struct{ current, latest string }{
		{"2.0.0", "1.5.0"},
		{"1.5.0", "1.3.9"},
	}
	for _, tc := range cases {
		if got := GetUpgradeType(tc.current, tc.latest); got != None {
			t.Errorf("GetUpgradeType(%q, %q) = %q, want %q", tc.current, tc.latest, got, None)
		}
	}
}

func TestGetUpgradeType_WhenTwoSegments_Works(t *testing.T) {
	if got := GetUpgradeType("1.0", "2.0"); got != Major {
		t.Errorf("got %q, want %q", got, Major)
	}
	if got := GetUpgradeType("1.0", "1.1"); got != Minor {
		t.Errorf("got %q, want %q", got, Minor)
	}
}
