package update

import (
	"strconv"
	"strings"
)

type Level string

const (
	LevelPatch Level = "patch"
	LevelMinor Level = "minor"
	LevelMajor Level = "major"
)

type Upgrade struct {
	Alias         string
	Section       string
	VersionRef    string
	InlineVersion bool
	Current       string
	Latest        string
	Type          string
}

type RefUpdate struct {
	Ref        string
	OldVersion string
	NewVersion string
	Aliases    []string
}

type InlineUpdate struct {
	Alias      string
	Section    string
	OldVersion string
	NewVersion string
	Type       string
}

type Result struct {
	VersionRefUpdates []RefUpdate
	InlineUpdates     []InlineUpdate
}

func Plan(upgrades []Upgrade, level Level) Result {
	var out Result
	refGroups := make(map[string][]Upgrade)
	for _, u := range upgrades {
		if u.Type == "none" || u.Type == "" {
			continue
		}
		if u.VersionRef != "" {
			refGroups[u.VersionRef] = append(refGroups[u.VersionRef], u)
			continue
		}
		if u.InlineVersion && fitsLevel(u.Type, level) {
			out.InlineUpdates = append(out.InlineUpdates, InlineUpdate{
				Alias:      u.Alias,
				Section:    u.Section,
				OldVersion: u.Current,
				NewVersion: u.Latest,
				Type:       u.Type,
			})
		}
	}
	for ref, group := range refGroups {
		if ru, ok := planRefUpdate(ref, group, level); ok {
			out.VersionRefUpdates = append(out.VersionRefUpdates, ru)
		}
	}
	return out
}

func planRefUpdate(ref string, group []Upgrade, level Level) (RefUpdate, bool) {
	if len(group) == 0 {
		return RefUpdate{}, false
	}
	candidates := make([]string, 0, len(group))
	aliases := make([]string, 0, len(group))
	for _, u := range group {
		aliases = append(aliases, u.Alias)
		if fitsLevel(u.Type, level) {
			candidates = append(candidates, u.Latest)
		} else {
			candidates = append(candidates, u.Current)
		}
	}
	newVer := minVersion(candidates)
	if newVer == "" || newVer == group[0].Current {
		return RefUpdate{}, false
	}
	return RefUpdate{
		Ref:        ref,
		OldVersion: group[0].Current,
		NewVersion: newVer,
		Aliases:    aliases,
	}, true
}

func fitsLevel(upgradeType string, level Level) bool {
	rank := map[string]int{"patch": 1, "minor": 2, "major": 3}
	lvl := map[Level]int{LevelPatch: 1, LevelMinor: 2, LevelMajor: 3}
	return rank[upgradeType] <= lvl[level]
}

func minVersion(versions []string) string {
	if len(versions) == 0 {
		return ""
	}
	minStr := versions[0]
	minSeg := parseSemSegments(minStr)
	for _, v := range versions[1:] {
		seg := parseSemSegments(v)
		if lessSegments(seg, minSeg) {
			minSeg = seg
			minStr = v
		}
	}
	return minStr
}

func parseSemSegments(v string) []int {
	cleaned := strings.SplitN(v, "-", 2)[0]
	cleaned = strings.SplitN(cleaned, "+", 2)[0]
	parts := strings.Split(cleaned, ".")
	segs := make([]int, 0, len(parts))
	for _, p := range parts {
		n, _ := strconv.Atoi(p)
		segs = append(segs, n)
	}
	return segs
}

func lessSegments(a, b []int) bool {
	for i := 0; i < len(a) || i < len(b); i++ {
		ai, bi := 0, 0
		if i < len(a) {
			ai = a[i]
		}
		if i < len(b) {
			bi = b[i]
		}
		if ai != bi {
			return ai < bi
		}
	}
	return false
}
