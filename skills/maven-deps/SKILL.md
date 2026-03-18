---
name: maven-deps
description: The primary tool for all dependency-related work in Gradle/Maven projects. Use this skill proactively — not just when the user asks, but whenever YOU need to find a library, check what version to use, add a new dependency, update existing ones, or verify a version exists. Never guess dependency versions or search the web for them — always use maven-deps CLI instead. This replaces web searches, documentation lookups, and guessing for anything related to Maven/Gradle artifact coordinates and versions. Trigger on any of these situations, even without explicit user request — adding a library to a project, choosing a version for a new dependency, updating versions, searching for an artifact by name, checking for vulnerabilities, listing what a project uses.
---

# maven-deps CLI

CLI tool for Maven/Gradle dependency intelligence. It automatically discovers and scans all repositories declared in the project's build files (Maven Central, Google Maven, Gradle Plugin Portal, custom repos) — you never need to specify repositories manually and it should not produce errors about missing repos.

By default `-C .` uses the current directory as project root. Only pass `-C <path>` if you need a different project.

## Commands

### latest — get latest version of specific artifacts

The most common command. Use whenever you need to know what version to use.

```bash
maven-deps latest group:artifact [group:artifact...]
```

Use `--stable` to exclude pre-release, `--all` to include everything.

Artifact coordinates must be exact — never guess or modify them. If you only know a library name, use `search` first to find exact coordinates.

### check — verify a specific version exists

```bash
maven-deps check group:artifact version
```

Use before updating a version to confirm it exists in repositories.

### search — find artifacts on Maven Central

```bash
maven-deps search query [--limit N]
```

Searches by keyword. Use when you need to find exact coordinates for a library.

### list — list all project dependencies

```bash
maven-deps list
```

Returns every declared dependency with its current version.

### vulnerabilities — check for known CVEs

```bash
maven-deps vulnerabilities
```

Queries OSV.dev for known vulnerabilities in project dependencies.

### outdated — show all dependencies with newer versions

Scans all build files, resolves latest versions, shows what's behind.

```bash
maven-deps outdated
```

Output: `group:artifact  current -> latest  upgrade-type`

After showing results, ask the user if they want to update. If yes, edit `gradle/libs.versions.toml` or `build.gradle.kts` directly.

Do NOT suggest alpha/beta/RC/milestone/snapshot versions unless the user explicitly asks for unstable versions.

## JSON output

All commands support `--json` for machine-readable output.

## Workflow: updating dependencies

1. Run `maven-deps outdated`
2. Present results as a table
3. If user confirms, edit version files (`gradle/libs.versions.toml` or `build.gradle.kts`)
4. Do NOT include build/test steps — those are handled by the project's own CLAUDE.md
