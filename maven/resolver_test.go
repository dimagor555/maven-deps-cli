package maven

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

const metadataXML = `<metadata><groupId>g</groupId><artifactId>a</artifactId><versioning><latest>1.1.0</latest><release>1.1.0</release><versions><version>1.0.0</version><version>1.1.0</version></versions></versioning></metadata>`

func TestResolveAll_RetriesTransientErrorThenSucceeds(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(metadataXML))
	}))
	defer srv.Close()

	repos := []Repository{NewRepository("test", srv.URL)}
	meta, err := ResolveAll(context.Background(), repos, "g", "a")
	if err != nil {
		t.Fatalf("expected success after retry, got %v", err)
	}
	if len(meta.Versions) != 2 {
		t.Errorf("expected 2 versions, got %v", meta.Versions)
	}
	if atomic.LoadInt32(&calls) < 2 {
		t.Errorf("expected at least one retry, calls=%d", calls)
	}
}

func TestResolveAll_AllTransientFailures_ReturnsErrorNotNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	repos := []Repository{NewRepository("test", srv.URL)}
	_, err := ResolveAll(context.Background(), repos, "g", "a")
	if err == nil {
		t.Fatal("expected error when all repos fail transiently")
	}
	if IsNotFound(err) {
		t.Errorf("transient failure must not be classified as not-found: %v", err)
	}
}

func TestResolveAll_RateLimitedThenSucceeds(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Write([]byte(metadataXML))
	}))
	defer srv.Close()

	repos := []Repository{NewRepository("test", srv.URL)}
	meta, err := ResolveAll(context.Background(), repos, "g", "a")
	if err != nil {
		t.Fatalf("expected recovery after 429, got %v", err)
	}
	if len(meta.Versions) != 2 {
		t.Errorf("expected 2 versions, got %v", meta.Versions)
	}
}

func TestResolveAll_NotFound_ClassifiedAsNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	repos := []Repository{NewRepository("test", srv.URL)}
	_, err := ResolveAll(context.Background(), repos, "g", "a")
	if err == nil {
		t.Fatal("expected error for missing artifact")
	}
	if !IsNotFound(err) {
		t.Errorf("404 from all repos must be classified as not-found: %v", err)
	}
}
