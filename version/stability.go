package version

import "regexp"

type Stability string

const (
	Stable    Stability = "stable"
	RC        Stability = "rc"
	Beta      Stability = "beta"
	Alpha     Stability = "alpha"
	Milestone Stability = "milestone"
	Snapshot  Stability = "snapshot"
)

type StabilityFilter string

const (
	StableOnly   StabilityFilter = "STABLE_ONLY"
	PreferStable StabilityFilter = "PREFER_STABLE"
	All          StabilityFilter = "ALL"
)

var stabilityPatterns = []struct {
	re        *regexp.Regexp
	stability Stability
}{
	{regexp.MustCompile(`(?i)[-.]?SNAPSHOT$`), Snapshot},
	{regexp.MustCompile(`(?i)[-.](?:alpha|a(?:\d|[-.]|$))[-.]?\d*`), Alpha},
	{regexp.MustCompile(`(?i)[-.](?:beta|b(?:\d|[-.]|$))[-.]?\d*`), Beta},
	{regexp.MustCompile(`(?i)[-.](?:M|milestone)[-.]?\d*`), Milestone},
	{regexp.MustCompile(`(?i)[-.](?:RC|CR)[-.]?\d*`), RC},
}

func Classify(version string) Stability {
	for _, p := range stabilityPatterns {
		if p.re.MatchString(version) {
			return p.stability
		}
	}
	return Stable
}
