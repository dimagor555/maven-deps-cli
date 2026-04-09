package gradle

import "testing"

func TestUpdateInlineVersions_MapNotation(t *testing.T) {
	content := `[libraries]
gson = { module = "com.google.code.gson:gson", version = "2.11.0" }
`
	got, err := UpdateInlineVersions(content, []InlineReplacement{
		{Alias: "gson", Section: "libraries", NewVersion: "2.12.0"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `[libraries]
gson = { module = "com.google.code.gson:gson", version = "2.12.0" }
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestUpdateInlineVersions_StringNotation(t *testing.T) {
	content := `[libraries]
gson = "com.google.code.gson:gson:2.11.0"
`
	got, err := UpdateInlineVersions(content, []InlineReplacement{
		{Alias: "gson", Section: "libraries", NewVersion: "2.12.0"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `[libraries]
gson = "com.google.code.gson:gson:2.12.0"
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestUpdateInlineVersions_Plugin(t *testing.T) {
	content := `[plugins]
foo = { id = "com.foo", version = "1.0.0" }
`
	got, err := UpdateInlineVersions(content, []InlineReplacement{
		{Alias: "foo", Section: "plugins", NewVersion: "1.1.0"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `[plugins]
foo = { id = "com.foo", version = "1.1.0" }
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestUpdateInlineVersions_DoesNotTouchVersionRef(t *testing.T) {
	content := `[libraries]
ktor = { module = "io.ktor:ktor-client-core", version.ref = "ktor" }
`
	_, err := UpdateInlineVersions(content, []InlineReplacement{
		{Alias: "ktor", Section: "libraries", NewVersion: "9.9.9"},
	})
	if err == nil {
		t.Errorf("expected error: cannot inline-update a version.ref entry")
	}
}

func TestUpdateInlineVersions_PreservesOtherFields(t *testing.T) {
	content := `[libraries]
gson = { group = "com.google.code.gson", name = "gson", version = "2.11.0" }
`
	got, err := UpdateInlineVersions(content, []InlineReplacement{
		{Alias: "gson", Section: "libraries", NewVersion: "2.12.0"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `[libraries]
gson = { group = "com.google.code.gson", name = "gson", version = "2.12.0" }
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestUpdateInlineVersions_IgnoresOtherSection(t *testing.T) {
	content := `[versions]
gson = "2.11.0"

[libraries]
gson = "com.google.code.gson:gson:2.11.0"
`
	got, err := UpdateInlineVersions(content, []InlineReplacement{
		{Alias: "gson", Section: "libraries", NewVersion: "2.12.0"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := `[versions]
gson = "2.11.0"

[libraries]
gson = "com.google.code.gson:gson:2.12.0"
`
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestUpdateInlineVersions_EmptyReplacements_ReturnsOriginal(t *testing.T) {
	content := `[libraries]
gson = "com.google.code.gson:gson:2.11.0"
`
	got, err := UpdateInlineVersions(content, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got != content {
		t.Errorf("got:\n%s", got)
	}
}

func TestUpdateInlineVersions_UnknownAlias_Errors(t *testing.T) {
	content := `[libraries]
gson = "com.google.code.gson:gson:2.11.0"
`
	_, err := UpdateInlineVersions(content, []InlineReplacement{
		{Alias: "nope", Section: "libraries", NewVersion: "1.0.0"},
	})
	if err == nil {
		t.Errorf("expected error for unknown alias")
	}
}
