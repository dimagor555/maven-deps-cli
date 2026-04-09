# maven-deps-cli

Fast CLI tool for Maven/Gradle dependency intelligence. Single binary, zero overhead.

Scans Gradle build files and version catalogs, queries Maven Central, Google Maven, Gradle Plugin Portal, and any custom repositories declared in your project.

## Install

Requires Go 1.22+.

### Build from source

```bash
git clone https://github.com/dimagor555/maven-deps-cli.git
cd maven-deps-cli
go install ./cmd/maven-deps
```

The binary is installed to `$GOPATH/bin` (usually `~/go/bin`), make sure it's in your `PATH`.

### Update

```bash
cd maven-deps-cli
git pull
go install ./cmd/maven-deps
```

## Commands

### `latest` — get latest versions

```bash
maven-deps latest io.insert-koin:koin-core org.jetbrains.kotlinx:kotlinx-serialization-json
```

### `check` — verify a version exists

```bash
maven-deps check io.insert-koin:koin-core 4.1.1
```

### `search` — find artifacts by keyword

```bash
maven-deps search flowmvi
```

### `list` — list all project dependencies

```bash
maven-deps list
```

### `outdated` — show dependencies with newer versions

```bash
maven-deps outdated
```

### `update` — bump catalog dependencies

```bash
maven-deps update patch   # only patch upgrades
maven-deps update minor   # patch + minor
maven-deps update major   # patch + minor + major
```

Scans `gradle/libs.versions.toml`, resolves latest versions, shows a preview
and asks to confirm before writing. Edits are line-based and preserve
formatting/comments. Shared `version.ref` entries are bumped to the minimum
of their latest versions so all consumers stay compatible.

Flags:
- `--yes` — skip confirmation
- `--force-dirty` — proceed even if the catalog has uncommitted changes

### `vulnerabilities` — check for known CVEs

```bash
maven-deps vulnerabilities
```

## Flags

| Flag | Description |
|------|-------------|
| `-C <path>` | Project path (default: `.`) |
| `--json` | JSON output |
| `--stable` | Only stable versions (for `latest`) |
| `--all` | Include pre-release (for `latest`) |

## Claude Code Plugin

This repo includes a Claude Code plugin with a `maven-deps` skill. To install:

```bash
claude plugin marketplace add dimagor555/maven-deps-cli
claude plugin install maven-deps@dimagor555
```

The skill teaches Claude to use `maven-deps` automatically for all dependency-related tasks — finding libraries, checking versions, updating dependencies, and scanning for vulnerabilities.

## License

MIT
