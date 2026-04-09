package gradle

import (
	"fmt"
	"regexp"
	"strings"
)

type Replacement struct {
	Alias      string
	NewVersion string
}

var versionAliasLine = regexp.MustCompile(`^(\s*)([A-Za-z0-9_.-]+)(\s*=\s*)(["'])([^"']*)(["'])(.*)$`)

func UpdateVersionAliases(content string, replacements []Replacement) (string, error) {
	if len(replacements) == 0 {
		return content, nil
	}
	byAlias := make(map[string]string, len(replacements))
	applied := make(map[string]bool, len(replacements))
	for _, r := range replacements {
		byAlias[r.Alias] = r.NewVersion
	}

	lines := splitLinesKeepEnds(content)
	section := ""
	for i, line := range lines {
		trimmed := strings.TrimSpace(stripLineEnding(line))
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			section = strings.Trim(trimmed, "[]")
			continue
		}
		if section != "versions" {
			continue
		}
		m := versionAliasLine.FindStringSubmatch(stripLineEnding(line))
		if m == nil {
			continue
		}
		alias := m[2]
		newVer, ok := byAlias[alias]
		if !ok {
			continue
		}
		ending := lineEnding(line)
		lines[i] = m[1] + alias + m[3] + m[4] + newVer + m[6] + m[7] + ending
		applied[alias] = true
	}

	for alias := range byAlias {
		if !applied[alias] {
			return "", fmt.Errorf("alias %q not found in [versions]", alias)
		}
	}
	return strings.Join(lines, ""), nil
}

func splitLinesKeepEnds(s string) []string {
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			out = append(out, s[start:i+1])
			start = i + 1
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func stripLineEnding(line string) string {
	return strings.TrimRight(line, "\r\n")
}

func lineEnding(line string) string {
	stripped := stripLineEnding(line)
	return line[len(stripped):]
}
