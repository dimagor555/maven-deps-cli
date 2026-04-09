package gradle

import (
	"fmt"
	"regexp"
	"strings"
)

type InlineReplacement struct {
	Alias      string
	Section    string
	NewVersion string
}

var (
	aliasLineRe      = regexp.MustCompile(`^(\s*)([A-Za-z0-9_.-]+)(\s*=\s*)(.*)$`)
	stringCoordRe    = regexp.MustCompile(`^(["'])([^:"']+:[^:"']+):([^"']+)(["'])(.*)$`)
	inlineVersionRe  = regexp.MustCompile(`(version\s*=\s*)(["'])([^"']*)(["'])`)
	versionRefMarker = regexp.MustCompile(`version\.ref\s*=`)
)

func UpdateInlineVersions(content string, replacements []InlineReplacement) (string, error) {
	if len(replacements) == 0 {
		return content, nil
	}
	bySection := groupInlineByAlias(replacements)
	applied := make(map[string]bool)

	lines := splitLinesKeepEnds(content)
	section := ""
	for i, line := range lines {
		stripped := stripLineEnding(line)
		trimmed := strings.TrimSpace(stripped)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			section = strings.Trim(trimmed, "[]")
			continue
		}
		targets, ok := bySection[section]
		if !ok {
			continue
		}
		newLine, key, err := tryReplaceInlineLine(stripped, targets)
		if err != nil {
			return "", err
		}
		if key != "" {
			lines[i] = newLine + lineEnding(line)
			applied[section+":"+key] = true
		}
	}
	return finalizeInline(lines, replacements, applied)
}

func groupInlineByAlias(reps []InlineReplacement) map[string]map[string]string {
	out := make(map[string]map[string]string)
	for _, r := range reps {
		if out[r.Section] == nil {
			out[r.Section] = make(map[string]string)
		}
		out[r.Section][r.Alias] = r.NewVersion
	}
	return out
}

func tryReplaceInlineLine(stripped string, targets map[string]string) (string, string, error) {
	m := aliasLineRe.FindStringSubmatch(stripped)
	if m == nil {
		return "", "", nil
	}
	alias := m[2]
	newVer, ok := targets[alias]
	if !ok {
		return "", "", nil
	}
	rest := m[4]
	if versionRefMarker.MatchString(rest) {
		return "", "", fmt.Errorf("alias %q uses version.ref, cannot inline-update", alias)
	}
	newRest, ok := replaceInlineValue(rest, newVer)
	if !ok {
		return "", "", fmt.Errorf("alias %q: no inline version found", alias)
	}
	return m[1] + alias + m[3] + newRest, alias, nil
}

func replaceInlineValue(rest, newVer string) (string, bool) {
	if sm := stringCoordRe.FindStringSubmatch(rest); sm != nil {
		return sm[1] + sm[2] + ":" + newVer + sm[4] + sm[5], true
	}
	if inlineVersionRe.MatchString(rest) {
		return inlineVersionRe.ReplaceAllString(rest, `${1}${2}`+newVer+`${4}`), true
	}
	return "", false
}

func finalizeInline(lines []string, reps []InlineReplacement, applied map[string]bool) (string, error) {
	for _, r := range reps {
		if !applied[r.Section+":"+r.Alias] {
			return "", fmt.Errorf("alias %q not found in [%s]", r.Alias, r.Section)
		}
	}
	return strings.Join(lines, ""), nil
}
