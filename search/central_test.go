package search

import "testing"

func TestBuildSolrQuery_WhenGroupArtifact_ReturnsGAQuery(t *testing.T) {
	got := buildSolrQuery("io.ktor:ktor-client-core")
	want := `g:"io.ktor" AND a:"ktor-client-core"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildSolrQuery_WhenGroupOnly_ReturnsGroupQuery(t *testing.T) {
	got := buildSolrQuery("io.ktor:")
	want := `g:"io.ktor"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildSolrQuery_WhenDottedName_ReturnsGroupQuery(t *testing.T) {
	got := buildSolrQuery("io.ktor")
	want := `g:"io.ktor"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBuildSolrQuery_WhenKeyword_ReturnsAsIs(t *testing.T) {
	got := buildSolrQuery("kotlinx serialization")
	want := "kotlinx serialization"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
