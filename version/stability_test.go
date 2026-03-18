package version

import "testing"

func TestClassify_WhenStableVersion_ReturnsStable(t *testing.T) {
	cases := []string{"3.5.11", "1.0", "2.0.0"}
	for _, v := range cases {
		if got := Classify(v); got != Stable {
			t.Errorf("Classify(%q) = %q, want %q", v, got, Stable)
		}
	}
}

func TestClassify_WhenSnapshot_ReturnsSnapshot(t *testing.T) {
	cases := []string{"1.0-SNAPSHOT", "2.0.0-SNAPSHOT"}
	for _, v := range cases {
		if got := Classify(v); got != Snapshot {
			t.Errorf("Classify(%q) = %q, want %q", v, got, Snapshot)
		}
	}
}

func TestClassify_WhenAlpha_ReturnsAlpha(t *testing.T) {
	cases := []string{"1.0-alpha-1", "1.0.0-alpha1", "1.0-a1"}
	for _, v := range cases {
		if got := Classify(v); got != Alpha {
			t.Errorf("Classify(%q) = %q, want %q", v, got, Alpha)
		}
	}
}

func TestClassify_WhenBeta_ReturnsBeta(t *testing.T) {
	cases := []string{"1.0-beta-1", "1.0.0-beta1", "1.0-b1"}
	for _, v := range cases {
		if got := Classify(v); got != Beta {
			t.Errorf("Classify(%q) = %q, want %q", v, got, Beta)
		}
	}
}

func TestClassify_WhenRC_ReturnsRC(t *testing.T) {
	cases := []string{"1.0-RC1", "1.0-rc-2", "1.0-CR1"}
	for _, v := range cases {
		if got := Classify(v); got != RC {
			t.Errorf("Classify(%q) = %q, want %q", v, got, RC)
		}
	}
}

func TestClassify_WhenMilestone_ReturnsMilestone(t *testing.T) {
	cases := []string{"1.0-M1", "1.0-milestone-2"}
	for _, v := range cases {
		if got := Classify(v); got != Milestone {
			t.Errorf("Classify(%q) = %q, want %q", v, got, Milestone)
		}
	}
}

func TestClassify_WhenNoFalsePositive_ReturnsStable(t *testing.T) {
	cases := []string{"1.0-bar", "1.0-ace"}
	for _, v := range cases {
		if got := Classify(v); got != Stable {
			t.Errorf("Classify(%q) = %q, want %q", v, got, Stable)
		}
	}
}
