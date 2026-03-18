package version

import (
	"strconv"
	"strings"
)

type UpgradeType string

const (
	Major UpgradeType = "major"
	Minor UpgradeType = "minor"
	Patch UpgradeType = "patch"
	None  UpgradeType = "none"
)

func GetUpgradeType(current, latest string) UpgradeType {
	cur := parseSegments(current)
	lat := parseSegments(latest)

	for len(cur) < len(lat) {
		cur = append(cur, 0)
	}
	for len(lat) < len(cur) {
		lat = append(lat, 0)
	}

	if lat[0] != cur[0] {
		if lat[0] > cur[0] {
			return Major
		}
		return None
	}
	if len(cur) > 1 && lat[1] != cur[1] {
		if lat[1] > cur[1] {
			return Minor
		}
		return None
	}
	if len(cur) > 2 && lat[2] != cur[2] {
		if lat[2] > cur[2] {
			return Patch
		}
		return None
	}
	return None
}

func parseSegments(version string) []int {
	cleaned := strings.SplitN(version, "-", 2)[0]
	cleaned = strings.SplitN(cleaned, "+", 2)[0]
	parts := strings.Split(cleaned, ".")
	segments := make([]int, 0, len(parts))
	for _, p := range parts {
		n, _ := strconv.Atoi(p)
		segments = append(segments, n)
	}
	return segments
}
