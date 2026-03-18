package gradle

import (
	"strings"

	"github.com/BurntSushi/toml"
)

type CatalogEntry struct {
	GroupID    string
	ArtifactID string
	Version    string
}

type rawCatalog struct {
	Versions  map[string]interface{} `toml:"versions"`
	Libraries map[string]interface{} `toml:"libraries"`
	Plugins   map[string]interface{} `toml:"plugins"`
}

func ParseVersionCatalog(content string) map[string]CatalogEntry {
	var raw rawCatalog
	if _, err := toml.Decode(content, &raw); err != nil {
		return nil
	}

	versions := parseVersions(raw.Versions)
	entries := make(map[string]CatalogEntry)

	for alias, val := range raw.Libraries {
		if e, ok := parseLibrary(val, versions); ok {
			entries[alias] = e
		}
	}

	for alias, val := range raw.Plugins {
		if e, ok := parsePlugin(val, versions); ok {
			entries["plugin-"+alias] = e
		}
	}

	return entries
}

func parseVersions(raw map[string]interface{}) map[string]string {
	versions := make(map[string]string)
	for k, v := range raw {
		if s, ok := v.(string); ok {
			versions[k] = s
		}
	}
	return versions
}

func parseLibrary(val interface{}, versions map[string]string) (CatalogEntry, bool) {
	switch v := val.(type) {
	case string:
		return parseModuleString(v)
	case map[string]interface{}:
		return parseLibraryMap(v, versions)
	}
	return CatalogEntry{}, false
}

func parseModuleString(s string) (CatalogEntry, bool) {
	parts := strings.SplitN(s, ":", 3)
	if len(parts) < 2 {
		return CatalogEntry{}, false
	}
	entry := CatalogEntry{GroupID: parts[0], ArtifactID: parts[1]}
	if len(parts) == 3 {
		entry.Version = parts[2]
	}
	return entry, true
}

func parseLibraryMap(m map[string]interface{}, versions map[string]string) (CatalogEntry, bool) {
	var groupID, artifactID, ver string

	if module, ok := getString(m, "module"); ok {
		parts := strings.SplitN(module, ":", 2)
		if len(parts) == 2 {
			groupID, artifactID = parts[0], parts[1]
		}
	}

	if g, ok := getString(m, "group"); ok {
		groupID = g
	}
	if n, ok := getString(m, "name"); ok {
		artifactID = n
	}

	ver = resolveVersion(m, versions)

	if groupID == "" || artifactID == "" {
		return CatalogEntry{}, false
	}
	return CatalogEntry{GroupID: groupID, ArtifactID: artifactID, Version: ver}, true
}

func parsePlugin(val interface{}, versions map[string]string) (CatalogEntry, bool) {
	m, ok := val.(map[string]interface{})
	if !ok {
		return CatalogEntry{}, false
	}
	id, ok := getString(m, "id")
	if !ok {
		return CatalogEntry{}, false
	}
	ver := resolveVersion(m, versions)
	return CatalogEntry{
		GroupID:    id,
		ArtifactID: id + ".gradle.plugin",
		Version:    ver,
	}, true
}

func resolveVersion(m map[string]interface{}, versions map[string]string) string {
	if vObj, ok := m["version"]; ok {
		if s, ok := vObj.(string); ok {
			return s
		}
		if vm, ok := vObj.(map[string]interface{}); ok {
			if ref, ok := getString(vm, "ref"); ok {
				return versions[ref]
			}
		}
	}
	if vr, ok := getString(m, "version.ref"); ok {
		return versions[vr]
	}
	return ""
}

func getString(m map[string]interface{}, key string) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
