package search

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"dimagor555.pro/maven-deps/httputil"
)

const searchAPI = "https://search.maven.org/solrsearch/select"

type Artifact struct {
	GroupID      string `json:"groupId"`
	ArtifactID   string `json:"artifactId"`
	LatestVersion string `json:"latestVersion"`
	VersionCount int    `json:"versionCount"`
}

type solrResponse struct {
	Response struct {
		Docs []struct {
			G             string `json:"g"`
			A             string `json:"a"`
			LatestVersion string `json:"latestVersion"`
			VersionCount  int    `json:"versionCount"`
		} `json:"docs"`
	} `json:"response"`
}

func Search(ctx context.Context, query string, limit int) ([]Artifact, error) {
	if limit <= 0 {
		limit = 10
	}

	solrQuery := buildSolrQuery(query)
	reqURL := fmt.Sprintf("%s?q=%s&rows=%d&wt=json", searchAPI, url.QueryEscape(solrQuery), limit)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := httputil.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("search maven central: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search maven central: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	var data solrResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	artifacts := make([]Artifact, len(data.Response.Docs))
	for i, doc := range data.Response.Docs {
		artifacts[i] = Artifact{
			GroupID:      doc.G,
			ArtifactID:   doc.A,
			LatestVersion: doc.LatestVersion,
			VersionCount: doc.VersionCount,
		}
	}
	return artifacts, nil
}

func buildSolrQuery(query string) string {
	trimmed := strings.TrimSpace(query)

	if strings.Contains(trimmed, ":") {
		parts := strings.SplitN(trimmed, ":", 2)
		if parts[1] != "" {
			return fmt.Sprintf(`g:"%s" AND a:"%s"`, parts[0], parts[1])
		}
		return fmt.Sprintf(`g:"%s"`, parts[0])
	}

	if strings.Contains(trimmed, ".") && !strings.Contains(trimmed, " ") {
		return fmt.Sprintf(`g:"%s"`, trimmed)
	}

	return trimmed
}
