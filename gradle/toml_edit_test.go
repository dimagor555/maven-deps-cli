package gradle

import "testing"

func TestUpdateVersionAliases_SingleReplacement(t *testing.T) {
	content := `[versions]
ktor = "3.1.1"
kotlin = "2.1.0"
`
	got, err := UpdateVersionAliases(content, []Replacement{{Alias: "ktor", NewVersion: "3.2.0"}})
	if err != nil {
		t.Fatal(err)
	}
	want := `[versions]
ktor = "3.2.0"
kotlin = "2.1.0"
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestUpdateVersionAliases_MultipleReplacements(t *testing.T) {
	content := `[versions]
ktor = "3.1.1"
kotlin = "2.1.0"
`
	got, err := UpdateVersionAliases(content, []Replacement{
		{Alias: "ktor", NewVersion: "3.2.0"},
		{Alias: "kotlin", NewVersion: "2.2.0"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `[versions]
ktor = "3.2.0"
kotlin = "2.2.0"
`
	if got != want {
		t.Errorf("got:\n%s", got)
	}
}

func TestUpdateVersionAliases_PreservesComments(t *testing.T) {
	content := `[versions]
# kotlin stack
ktor = "3.1.1"  # comment
kotlin = "2.1.0"
`
	got, err := UpdateVersionAliases(content, []Replacement{{Alias: "ktor", NewVersion: "3.2.0"}})
	if err != nil {
		t.Fatal(err)
	}
	want := `[versions]
# kotlin stack
ktor = "3.2.0"  # comment
kotlin = "2.1.0"
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestUpdateVersionAliases_PreservesBlankLines(t *testing.T) {
	content := `[versions]

ktor = "3.1.1"

kotlin = "2.1.0"
`
	got, err := UpdateVersionAliases(content, []Replacement{{Alias: "ktor", NewVersion: "3.2.0"}})
	if err != nil {
		t.Fatal(err)
	}
	want := `[versions]

ktor = "3.2.0"

kotlin = "2.1.0"
`
	if got != want {
		t.Errorf("got:\n%s", got)
	}
}

func TestUpdateVersionAliases_SingleQuotes(t *testing.T) {
	content := `[versions]
ktor = '3.1.1'
`
	got, err := UpdateVersionAliases(content, []Replacement{{Alias: "ktor", NewVersion: "3.2.0"}})
	if err != nil {
		t.Fatal(err)
	}
	want := `[versions]
ktor = '3.2.0'
`
	if got != want {
		t.Errorf("got:\n%s", got)
	}
}

func TestUpdateVersionAliases_IgnoresOtherSections(t *testing.T) {
	content := `[versions]
ktor = "3.1.1"

[libraries]
ktor = "io.ktor:ktor-client-core:3.1.1"
`
	got, err := UpdateVersionAliases(content, []Replacement{{Alias: "ktor", NewVersion: "3.2.0"}})
	if err != nil {
		t.Fatal(err)
	}
	want := `[versions]
ktor = "3.2.0"

[libraries]
ktor = "io.ktor:ktor-client-core:3.1.1"
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestUpdateVersionAliases_UnknownAlias_Errors(t *testing.T) {
	content := `[versions]
ktor = "3.1.1"
`
	_, err := UpdateVersionAliases(content, []Replacement{{Alias: "unknown", NewVersion: "1.0.0"}})
	if err == nil {
		t.Errorf("expected error for unknown alias")
	}
}

func TestUpdateVersionAliases_EmptyReplacements_ReturnsOriginal(t *testing.T) {
	content := `[versions]
ktor = "3.1.1"
`
	got, err := UpdateVersionAliases(content, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != content {
		t.Errorf("got:\n%s", got)
	}
}
