package version

var stabilityRank = map[Stability]int{
	Stable:    0,
	RC:        1,
	Milestone: 2,
	Beta:      3,
	Alpha:     4,
	Snapshot:  5,
}

func FindLatest(versions []string, filter StabilityFilter) string {
	if len(versions) == 0 {
		return ""
	}
	if filter == All {
		return versions[len(versions)-1]
	}
	stable := lastWhere(versions, func(v string) bool {
		return Classify(v) == Stable
	})
	if filter == StableOnly {
		return stable
	}
	if stable != "" {
		return stable
	}
	return versions[len(versions)-1]
}

func FindLatestForCurrent(versions []string, current string) string {
	maxRank := stabilityRank[Classify(current)]
	return lastWhere(versions, func(v string) bool {
		return stabilityRank[Classify(v)] <= maxRank
	})
}

func lastWhere(versions []string, predicate func(string) bool) string {
	for i := len(versions) - 1; i >= 0; i-- {
		if predicate(versions[i]) {
			return versions[i]
		}
	}
	return ""
}
