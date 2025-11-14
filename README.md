# go-sandbox

Multi-project Go repository (Go baseline: 1.25). Each subproject lives in its own module for clean boundaries, reproducible builds, and focused CI.

## Subprojects
- **todo-cli** — Local TODO manager CLI with JSON persistence, menu UI, and stdlib-only deps.
- **guessr** — Number guessing CLI with hints, stats tracking, and deterministic seeds.
- **filesort** — Directory sorter that buckets files by type with a dry-run preview.
- **weathertape** — Terminal weather dashboard that renders ASCII “tape” forecasts (module scaffolded; functionality queued up next).
- **snake** — Ebiten-based arcade snake clone with instant restarts and score overlay.

## Upcoming subprojects
- **MemeSweeper** — Ebiten desktop puzzler inspired by Minesweeper where tiles hide reaction memes instead of mines; timed boards, flag counts, and seeded layouts.
- **ThumbForge** — CLI for batch thumbnail generation: resize/crop images, preserve EXIF-safe metadata, and export fixed-size assets offline.

## Principles
- Standard library first. Any external dependency must be justified in PR notes.
- High test coverage; tests run offline by default.
- Commit flow per subproject: tests first (failing), then feature making tests pass.

## Build & Test
```bash
make test-all
```

## For AI agents & contributors
See [agents.md](agents.md) for repository conventions, PR workflow (tests-first), Go/toolchain rules, CI patterns, and subproject templates.
