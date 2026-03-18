# maven-deps

Go CLI tool for Maven/Gradle dependency intelligence. Single binary, zero overhead.

## Commands

```bash
go build -o maven-deps ./cmd/maven-deps   # Build
go test ./...                               # Run all tests
```

## Architecture

```
cmd/maven-deps/
  main.go              # cobra setup, global flags (-C, --json)
  latest.go            # latest command orchestration
  check.go             # check command orchestration
  search.go            # search command orchestration
  list.go              # list command orchestration
  outdated.go          # outdated command orchestration
  vulnerabilities.go   # vulnerabilities command orchestration
maven/
  repository.go        # Repository struct: HTTP fetch + XML parse
  wellknown.go         # Maven Central, Google, Gradle Plugin Portal instances
  resolver.go          # Parallel resolution across repos, dedup, proxy deprioritization
  metadata.go          # MavenMetadata struct
gradle/
  repo.go              # Parse repository declarations from Gradle DSL
  dep.go               # Parse dependency declarations from Gradle DSL
  toml.go              # Parse libs.versions.toml via BurntSushi/toml
project/
  root.go              # Find project root (walk up to settings.gradle)
  files.go             # Find build files (walk down, skip build/.gradle/.git)
version/
  stability.go         # Classify version stability
  filter.go            # Pick latest version by stability criteria
  compare.go           # Determine upgrade type (major/minor/patch)
search/
  central.go           # Maven Central Solr Search API client
vulnerability/
  osv.go               # OSV.dev batch API client
```

## Dependency direction

cmd/* → all packages (wiring layer)
All other packages → nothing from this project (leaf packages)

## Go code conventions (from welftrans reference)

- No comments in code
- Files: snake_case. Packages: single lowercase word
- No stuttering (maven.Repository not maven.MavenRepository)
- Errors: return error, wrap with fmt.Errorf("operation: %w", err)
- Interfaces: defined in consumer package, named by behavior
- Functions under 40 lines, files under 200 lines
- context.Context as first param for cancellable operations
- Logging: log/slog from stdlib
- Table-driven tests, fakes for deps
- Test naming: TestTypeName_WhenCondition_ExpectedResult
- TDD: test first → RED → implement → GREEN → refactor

## Go libraries

- `encoding/xml` — Maven metadata XML parsing
- `github.com/BurntSushi/toml` — libs.versions.toml
- `github.com/Masterminds/semver/v3` — version comparison
- `github.com/spf13/cobra` — CLI framework
- stdlib `net/http`, `regexp`
