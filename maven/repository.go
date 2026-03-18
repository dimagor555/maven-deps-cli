package maven

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"dimagor555.pro/maven-deps/httputil"
)

type Repository struct {
	Name string
	URL  string
}

func NewRepository(name, url string) Repository {
	return Repository{
		Name: name,
		URL:  strings.TrimRight(url, "/"),
	}
}

func (r Repository) MetadataURL(groupID, artifactID string) string {
	groupPath := strings.ReplaceAll(groupID, ".", "/")
	return fmt.Sprintf("%s/%s/%s/maven-metadata.xml", r.URL, groupPath, artifactID)
}

func (r Repository) FetchMetadata(ctx context.Context, groupID, artifactID string) (Metadata, error) {
	url := r.MetadataURL(groupID, artifactID)

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
		}
		m, err := r.doFetch(ctx, url, groupID, artifactID)
		if err == nil {
			return m, nil
		}
		if isNonRetryable(err) {
			return Metadata{}, err
		}
		lastErr = err
	}
	return Metadata{}, lastErr
}

func (r Repository) doFetch(ctx context.Context, url, groupID, artifactID string) (Metadata, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Metadata{}, fmt.Errorf("create request: %w", err)
	}

	resp, err := httputil.Client.Do(req)
	if err != nil {
		return Metadata{}, fmt.Errorf("fetch metadata from %s: %w", r.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Metadata{}, &notFoundError{repo: r.Name}
	}
	if resp.StatusCode != http.StatusOK {
		return Metadata{}, fmt.Errorf("fetch metadata from %s: %d %s", r.Name, resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Metadata{}, fmt.Errorf("read response from %s: %w", r.Name, err)
	}

	return ParseMetadata(body, groupID, artifactID)
}

type notFoundError struct {
	repo string
}

func (e *notFoundError) Error() string {
	return fmt.Sprintf("not found in %s", e.repo)
}

func isNonRetryable(err error) bool {
	if _, ok := err.(*notFoundError); ok {
		return true
	}
	if err, ok := err.(interface{ Timeout() bool }); ok && err.Timeout() {
		return true
	}
	return false
}
