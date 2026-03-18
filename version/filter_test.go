package version

import "testing"

func TestFindLatest_WhenAll_ReturnsLast(t *testing.T) {
	versions := []string{"1.0.0", "2.0.0-alpha1", "2.0.0"}
	if got := FindLatest(versions, All); got != "2.0.0" {
		t.Errorf("got %q, want %q", got, "2.0.0")
	}
}

func TestFindLatest_WhenStableOnly_ReturnsStable(t *testing.T) {
	versions := []string{"1.0.0", "2.0.0-alpha1"}
	if got := FindLatest(versions, StableOnly); got != "1.0.0" {
		t.Errorf("got %q, want %q", got, "1.0.0")
	}
}

func TestFindLatest_WhenPreferStable_FallsBackToPrerelease(t *testing.T) {
	versions := []string{"1.0.0-alpha1", "1.0.0-beta1"}
	if got := FindLatest(versions, PreferStable); got != "1.0.0-beta1" {
		t.Errorf("got %q, want %q", got, "1.0.0-beta1")
	}
}

func TestFindLatest_WhenEmpty_ReturnsEmpty(t *testing.T) {
	if got := FindLatest(nil, PreferStable); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestFindLatestForCurrent_WhenStable_ReturnsOnlyStable(t *testing.T) {
	versions := []string{"1.0.0", "2.0.0-alpha1", "2.0.0-beta1", "2.0.0-RC1", "2.0.0"}
	if got := FindLatestForCurrent(versions, "1.0.0"); got != "2.0.0" {
		t.Errorf("got %q, want %q", got, "2.0.0")
	}
}

func TestFindLatestForCurrent_WhenBeta_ReturnsBetaOrHigher(t *testing.T) {
	versions := []string{"1.0.0", "2.0.0-alpha1", "2.0.0-beta2"}
	if got := FindLatestForCurrent(versions, "1.0.0-beta1"); got != "2.0.0-beta2" {
		t.Errorf("got %q, want %q", got, "2.0.0-beta2")
	}
}

func TestFindLatestForCurrent_WhenAlpha_ReturnsAlphaOrHigher(t *testing.T) {
	versions := []string{"1.0.0-alpha1", "1.0.0-alpha2"}
	if got := FindLatestForCurrent(versions, "1.0.0-alpha1"); got != "1.0.0-alpha2" {
		t.Errorf("got %q, want %q", got, "1.0.0-alpha2")
	}
}

func TestFindLatestForCurrent_WhenNoMatch_ReturnsEmpty(t *testing.T) {
	versions := []string{"1.0.0-SNAPSHOT"}
	if got := FindLatestForCurrent(versions, "1.0.0"); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}
