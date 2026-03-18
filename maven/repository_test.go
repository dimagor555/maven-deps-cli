package maven

import "testing"

func TestRepository_MetadataURL_WhenStandardGroup_BuildsCorrectURL(t *testing.T) {
	r := NewRepository("test", "https://repo.example.com/maven2")
	got := r.MetadataURL("io.ktor", "ktor-server-core")
	want := "https://repo.example.com/maven2/io/ktor/ktor-server-core/maven-metadata.xml"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRepository_MetadataURL_WhenTrailingSlash_StripsIt(t *testing.T) {
	r := NewRepository("test", "https://repo.example.com/maven2/")
	got := r.MetadataURL("io.ktor", "ktor-server-core")
	want := "https://repo.example.com/maven2/io/ktor/ktor-server-core/maven-metadata.xml"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestNewRepository_WhenCreated_SetsFields(t *testing.T) {
	r := NewRepository("My Repo", "https://repo.example.com")
	if r.Name != "My Repo" {
		t.Errorf("Name = %q, want %q", r.Name, "My Repo")
	}
	if r.URL != "https://repo.example.com" {
		t.Errorf("URL = %q, want %q", r.URL, "https://repo.example.com")
	}
}
