package maven

import (
	"context"
	"fmt"
	"sync"
)

type Resolver struct {
	Repos []Repository

	mu    sync.Mutex
	cache map[string]cachedResult
}

type cachedResult struct {
	meta Metadata
	err  error
	done chan struct{}
}

func NewResolver(repos []Repository) *Resolver {
	return &Resolver{
		Repos: repos,
		cache: make(map[string]cachedResult),
	}
}

func (r *Resolver) Resolve(ctx context.Context, groupID, artifactID string) (Metadata, error) {
	key := groupID + ":" + artifactID

	r.mu.Lock()
	if cached, ok := r.cache[key]; ok {
		r.mu.Unlock()
		<-cached.done
		return cached.meta, cached.err
	}

	entry := cachedResult{done: make(chan struct{})}
	r.cache[key] = entry
	r.mu.Unlock()

	entry.meta, entry.err = resolveAll(ctx, r.Repos, groupID, artifactID)
	r.mu.Lock()
	r.cache[key] = entry
	r.mu.Unlock()
	close(entry.done)

	return entry.meta, entry.err
}

type repoResult struct {
	metadata Metadata
	repo     Repository
	err      error
}

func resolveAll(ctx context.Context, repos []Repository, groupID, artifactID string) (Metadata, error) {
	if len(repos) == 0 {
		return Metadata{}, fmt.Errorf("no repositories configured for %s:%s", groupID, artifactID)
	}

	results := make([]repoResult, len(repos))
	var wg sync.WaitGroup
	wg.Add(len(repos))

	for i, repo := range repos {
		go func(i int, repo Repository) {
			defer wg.Done()
			m, err := repo.FetchMetadata(ctx, groupID, artifactID)
			results[i] = repoResult{metadata: m, repo: repo, err: err}
		}(i, repo)
	}
	wg.Wait()

	hasCustom := false
	for _, r := range results {
		if r.err == nil && !ProxyTargetURLs[r.repo.URL] {
			hasCustom = true
			break
		}
	}

	seen := make(map[string]bool)
	var ordered []string
	var successful []Metadata

	for _, r := range results {
		if r.err != nil {
			continue
		}
		if hasCustom && ProxyTargetURLs[r.repo.URL] {
			continue
		}
		successful = append(successful, r.metadata)
		for _, v := range r.metadata.Versions {
			if !seen[v] {
				seen[v] = true
				ordered = append(ordered, v)
			}
		}
	}

	if len(successful) == 0 {
		return Metadata{}, fmt.Errorf("artifact %s:%s not found in any repository", groupID, artifactID)
	}

	last := ordered[len(ordered)-1]
	allLatest := collectField(successful, func(m Metadata) string { return m.Latest })
	allRelease := collectField(successful, func(m Metadata) string { return m.Release })

	return Metadata{
		GroupID:    groupID,
		ArtifactID: artifactID,
		Versions:   ordered,
		Latest:     pickBest(allLatest, last),
		Release:    pickBest(allRelease, last),
	}, nil
}

func ResolveAll(ctx context.Context, repos []Repository, groupID, artifactID string) (Metadata, error) {
	return resolveAll(ctx, repos, groupID, artifactID)
}

func collectField(metas []Metadata, getter func(Metadata) string) []string {
	var result []string
	for _, m := range metas {
		if v := getter(m); v != "" {
			result = append(result, v)
		}
	}
	return result
}

func pickBest(candidates []string, last string) string {
	for _, c := range candidates {
		if c == last {
			return last
		}
	}
	if len(candidates) > 0 {
		return candidates[len(candidates)-1]
	}
	return last
}
