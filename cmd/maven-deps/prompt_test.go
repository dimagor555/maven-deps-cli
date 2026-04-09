package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfirm_Yes_ReturnsTrue(t *testing.T) {
	for _, in := range []string{"y\n", "Y\n", "yes\n", "YES\n"} {
		var out bytes.Buffer
		got, err := Confirm(&out, strings.NewReader(in), "proceed?")
		if err != nil {
			t.Fatal(err)
		}
		if !got {
			t.Errorf("input %q: got false", in)
		}
	}
}

func TestConfirm_No_ReturnsFalse(t *testing.T) {
	for _, in := range []string{"n\n", "N\n", "no\n", "\n"} {
		var out bytes.Buffer
		got, err := Confirm(&out, strings.NewReader(in), "proceed?")
		if err != nil {
			t.Fatal(err)
		}
		if got {
			t.Errorf("input %q: got true", in)
		}
	}
}

func TestConfirm_EOF_ReturnsFalse(t *testing.T) {
	var out bytes.Buffer
	got, err := Confirm(&out, strings.NewReader(""), "proceed?")
	if err != nil {
		t.Fatal(err)
	}
	if got {
		t.Errorf("EOF: got true, want false")
	}
}

func TestConfirm_PromptWritten(t *testing.T) {
	var out bytes.Buffer
	_, _ = Confirm(&out, strings.NewReader("n\n"), "proceed?")
	if !strings.Contains(out.String(), "proceed?") {
		t.Errorf("prompt not written: %q", out.String())
	}
}
