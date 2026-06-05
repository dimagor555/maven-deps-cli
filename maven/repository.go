package maven

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"dimagor555.pro/maven-deps/httputil"
)

const (
	maxAttempts   = 4
	baseBackoff   = 400 * time.Millisecond
	maxBackoff    = 8 * time.Second
	rateLimitBase = 1500 * time.Millisecond
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
	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			if err := sleepBackoff(ctx, attempt, lastErr); err != nil {
				return Metadata{}, err
			}
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
	return Metadata{}, fmt.Errorf("after %d attempts: %w", maxAttempts, lastErr)
}

func sleepBackoff(ctx context.Context, attempt int, lastErr error) error {
	delay := backoffDuration(attempt, lastErr)
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func backoffDuration(attempt int, lastErr error) time.Duration {
	base := baseBackoff
	var rl *rateLimitError
	if errors.As(lastErr, &rl) {
		if rl.retryAfter > 0 {
			return rl.retryAfter + time.Duration(rand.Int63n(int64(rateLimitBase)))
		}
		base = rateLimitBase
	}
	backoff := base << (attempt - 1)
	if backoff > maxBackoff {
		backoff = maxBackoff
	}
	jitter := time.Duration(rand.Int63n(int64(base)))
	return backoff + jitter
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
	if resp.StatusCode == http.StatusTooManyRequests {
		return Metadata{}, &rateLimitError{
			repo:       r.Name,
			retryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		}
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

func NewNotFoundError(repo string) error {
	return &notFoundError{repo: repo}
}

type rateLimitError struct {
	repo       string
	retryAfter time.Duration
}

func (e *rateLimitError) Error() string {
	return fmt.Sprintf("rate limited by %s (429)", e.repo)
}

func parseRetryAfter(v string) time.Duration {
	if v == "" {
		return 0
	}
	if secs, err := strconv.Atoi(strings.TrimSpace(v)); err == nil && secs >= 0 {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		d := time.Until(t)
		if d > 0 {
			return d
		}
	}
	return 0
}

func (e *notFoundError) Error() string {
	return fmt.Sprintf("not found in %s", e.repo)
}

func IsNotFound(err error) bool {
	var nf *notFoundError
	return errors.As(err, &nf)
}

func isNonRetryable(err error) bool {
	if _, ok := err.(*notFoundError); ok {
		return true
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	return false
}
