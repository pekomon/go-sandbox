# go-sandbox

Multi-project Go repository (Go baseline: 1.24.3). Each subproject lives in its own module for clean boundaries, reproducible builds, and focused CI.

## Subprojects
- **todo-cli** — Local TODO manager CLI with JSON persistence (to be implemented in upcoming PRs).
- _Planned_: **guessr**, **filesort** (to be added later, one at a time).

## Principles
- Standard library first. Any external dependency must be justified in PR notes.
- High test coverage; tests run offline by default.
- Commit flow per subproject: tests first (failing), then feature making tests pass.

## Build & Test
```bash
make test-all
```
