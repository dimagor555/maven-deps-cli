package main

import (
	"context"
	"sync"

	"dimagor555.pro/maven-deps/gradle"
	"dimagor555.pro/maven-deps/update"
	"dimagor555.pro/maven-deps/version"
)

func resolveUpgrades(ctx context.Context, entries map[string]gradle.CatalogEntry, resolver dependencyResolver) []update.Upgrade {
	var targets []gradle.CatalogEntry
	for _, e := range entries {
		if e.Version != "" {
			targets = append(targets, e)
		}
	}
	results := make([]update.Upgrade, len(targets))
	sem := make(chan struct{}, maxConcurrency)
	var wg sync.WaitGroup
	wg.Add(len(targets))
	for i, e := range targets {
		go func(i int, e gradle.CatalogEntry) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i] = buildUpgrade(ctx, e, resolver)
		}(i, e)
	}
	wg.Wait()
	var out []update.Upgrade
	for _, r := range results {
		if r.Alias != "" {
			out = append(out, r)
		}
	}
	return out
}

func buildUpgrade(ctx context.Context, e gradle.CatalogEntry, resolver dependencyResolver) update.Upgrade {
	meta, err := resolver.Resolve(ctx, e.GroupID, e.ArtifactID)
	if err != nil {
		return update.Upgrade{}
	}
	latest := version.FindLatestForCurrent(meta.Versions, e.Version)
	if latest == "" || latest == e.Version {
		return update.Upgrade{}
	}
	ut := version.GetUpgradeType(e.Version, latest)
	if ut == version.None {
		return update.Upgrade{}
	}
	return update.Upgrade{
		Alias:         e.SourceAlias,
		Section:       e.Section,
		VersionRef:    e.VersionRef,
		InlineVersion: e.InlineVersion,
		Current:       e.Version,
		Latest:        latest,
		Type:          string(ut),
	}
}

func applyPlan(content string, plan update.Result) (string, error) {
	refReps := make([]gradle.Replacement, 0, len(plan.VersionRefUpdates))
	for _, ru := range plan.VersionRefUpdates {
		refReps = append(refReps, gradle.Replacement{Alias: ru.Ref, NewVersion: ru.NewVersion})
	}
	out, err := gradle.UpdateVersionAliases(content, refReps)
	if err != nil {
		return "", err
	}
	inlineReps := make([]gradle.InlineReplacement, 0, len(plan.InlineUpdates))
	for _, iu := range plan.InlineUpdates {
		inlineReps = append(inlineReps, gradle.InlineReplacement{
			Alias:      iu.Alias,
			Section:    iu.Section,
			NewVersion: iu.NewVersion,
		})
	}
	return gradle.UpdateInlineVersions(out, inlineReps)
}
