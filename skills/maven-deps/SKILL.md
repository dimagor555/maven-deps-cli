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

After showing results, ask the user if they want to update. If yes, prefer `maven-deps update` (see below) — it edits `gradle/libs.versions.toml` for you. Only edit files manually when the user explicitly asks or when the change is not expressible through the catalog.

Do NOT suggest alpha/beta/RC/milestone/snapshot versions unless the user explicitly asks for unstable versions.

### update — bump catalog dependencies automatically

Scans `gradle/libs.versions.toml`, resolves latest versions, previews all changes, asks for confirmation, then applies line-based edits that preserve formatting and comments.

```bash
maven-deps update patch   # only patch upgrades
maven-deps update minor   # patch + minor
maven-deps update major   # patch + minor + major
```

The upgrade level is positional and nested: `minor` includes patch, `major` includes minor and patch.

Shared `version.ref` entries (multiple aliases pointing at the same version key) are bumped to `min(latest)` across all consumers so nothing breaks. If any consumer has a latest that the chosen level cannot cover, the whole ref stays at its current value.

Blocks by default if `gradle/libs.versions.toml` has uncommitted changes in git. Pass `--force-dirty` to override.

Flags:
- `--yes` — skip the interactive confirmation prompt
- `--force-dirty` — proceed even if the catalog is dirty in git

Examples:

```bash
maven-deps update patch --yes              # safest: only patch bumps, no prompt
maven-deps update minor                    # interactive, asks before writing
maven-deps update major --force-dirty      # allow running on a dirty working tree
```

Only touches `gradle/libs.versions.toml`. For inline dependencies in `build.gradle.kts` (outside the catalog), edit manually.

## JSON output

All commands support `--json` for machine-readable output.

## Workflow: updating dependencies

1. Run `maven-deps outdated` to see what's behind
2. Present results as a table
3. If user confirms, run `maven-deps update <level>` (patch/minor/major) — it previews and writes the catalog
4. For inline deps outside `libs.versions.toml`, edit `build.gradle.kts` directly
5. Do NOT include build/test steps — those are handled by the project's own CLAUDE.md
