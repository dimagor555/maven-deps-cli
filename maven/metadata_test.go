package maven

import "testing"

func TestParseMetadata_WhenValidXML_ReturnsMetadata(t *testing.T) {
	xml := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<metadata>
  <groupId>io.ktor</groupId>
  <artifactId>ktor-server-core</artifactId>
  <versioning>
    <latest>3.1.1</latest>
    <release>3.1.1</release>
    <versions>
      <version>2.0.0</version>
      <version>3.0.0</version>
      <version>3.1.1</version>
    </versions>
  </versioning>
</metadata>`)

	m, err := ParseMetadata(xml, "io.ktor", "ktor-server-core")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.GroupID != "io.ktor" {
		t.Errorf("GroupID = %q, want %q", m.GroupID, "io.ktor")
	}
	if m.ArtifactID != "ktor-server-core" {
		t.Errorf("ArtifactID = %q, want %q", m.ArtifactID, "ktor-server-core")
	}
	if len(m.Versions) != 3 {
		t.Fatalf("len(Versions) = %d, want 3", len(m.Versions))
	}
	if m.Versions[2] != "3.1.1" {
		t.Errorf("Versions[2] = %q, want %q", m.Versions[2], "3.1.1")
	}
	if m.Latest != "3.1.1" {
		t.Errorf("Latest = %q, want %q", m.Latest, "3.1.1")
	}
	if m.Release != "3.1.1" {
		t.Errorf("Release = %q, want %q", m.Release, "3.1.1")
	}
}

func TestParseMetadata_WhenInvalidXML_ReturnsError(t *testing.T) {
	_, err := ParseMetadata([]byte("not xml"), "g", "a")
	if err == nil {
		t.Error("expected error for invalid XML")
	}
}
