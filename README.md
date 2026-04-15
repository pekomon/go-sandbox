# go-sandbox

Multi-project Go repository (Go baseline: 1.25). Each subproject lives in its own module for clean boundaries, reproducible builds, and focused CI.

## Subprojects
- **todo-cli** — Local TODO manager CLI with JSON persistence and an interactive menu UI.
- **guessr** — Number guessing CLI with hints, stats tracking, and deterministic seeds.
- **filesort** — Directory sorter that buckets files by type with a dry-run preview.
- **snake** — Ebiten-based arcade snake clone with instant restarts and score overlay.
- **brickbreaker** — Arkanoid-style Ebiten game with deterministic gameplay tests, three built-in levels, and a level-by-level speed ramp.
- **speedpads** — Speden Spelit inspired memory-speed game scaffold with deterministic sequence-state tests and an Ebiten placeholder shell.
- **memesweeper** — Ebiten puzzler inspired by Minesweeper with meme tiles, flags, and difficulty presets.
- **weathertape** — Terminal weather dashboard that renders ASCII “tape” forecasts from JSON data sources.
- **thumbforge** — CLI for batch thumbnail generation: resize images offline and export fixed-size assets.
- **dungeondice** — Dice-driven CLI roguelike simulator with classes, seeded runs, and combat summaries.
- **triviagoblin** — Trivia quiz CLI with deterministic shuffling, category filters, and a ready-to-run sample question set.

## Principles
- Standard library first. Any external dependency must be justified in PR notes.
- High test coverage; tests run offline by default.
- Commit flow per subproject: tests first (failing), then feature making tests pass.

## Build & Test
```bash
make list
make build-all
make test-all
```

## For AI agents & contributors
See [agents.md](agents.md) for repository conventions, PR workflow (tests-first), Go/toolchain rules, CI patterns, and subproject templates.
