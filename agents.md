# Agents guide

Audience: code-generation agents and contributors (Codex, Copilot, ChatGPT Coder).
Language: English only for code, comments, READMEs, commit messages.

## Repository shape
- Multi-project Go repo; **each subproject is its own module** (has its own `go.mod`).
- Current subprojects: `todo-cli/`, `guessr/`, `filesort/`, `snake/`, `brickbreaker/`, `memesweeper/`, `weathertape/`, `thumbforge/`, `dungeondice/`, `triviagoblin/`.
- Root contains top-level docs plus a coordinating `Makefile` for list/build/test/cover/clean loops across subprojects.

## Go & toolchain
Baseline: Go 1.25.

CI uses actions/setup-go@v5 with go-version: "1.25.x".

Prefer not to use toolchain directives in modules to avoid redundant downloads in CI.

## Workflow rules
- **One logical unit per PR**. Keep PRs focused and small.
- Prefer **tests-first**: a failing tests PR, then a feature PR that makes tests pass.
- Use **one GitHub issue per app/subproject effort**. Keep that issue high-level: summary, v1 scope, non-goals, and major milestones.
- Track detailed sequencing and temporary implementation notes in the local `.codex-plan.md` file. Treat it as a local planning aid; **do not commit it**.
- When multiple PRs contribute to the same app issue:
  - intermediate PRs should use `Refs #<n>`
  - only the final PR that completes the planned app slice should use `Closes #<n>`
- Branch naming: `feat/...`, `test/...`, `fix/...`, `docs/...`, `infra/...`, `ci/...`.
- Commit style:
  - tests-only PR: `test(<area>): ...`
  - implementation PR: `feat(<area>): ...`
  - docs-only PR: `docs(<area>): ...`
  - infra/ci: `chore(...)` or `ci(...)`
- **English** for all commit messages and docs.

## Coding standards
- **Stdlib-first**. Add external dependencies only when clearly justified in the PR body.
- Tests run **offline**, deterministic, no network.
- Use table-driven tests and keep them hermetic (`t.TempDir()`, env overrides).
- Clear separation of CLI (under `cmd/<name>/`) and logic (under `internal/` or `pkg/`).
- Return errors (do not panic) and write failures to **stderr**; normal output to **stdout**.
- Default to **newest-first** sorting by ID where appropriate; use flags to change.
- JSON persistence: pretty-print, atomic write via temp+rename, create parent dirs.

## Known environment variables & paths
- **todo-cli**
  - Storage file: `~/.todo-cli/tasks.json`
  - `TODO_CLI_PATH` — override the JSON file path
  - `TODO_CLI_MENU=1` — start in interactive menu mode if no args
- **guessr**
  - Stats file: `~/.guessr/stats.json`
  - `GUESSR_STATS_PATH` — override the JSON file path
  - Flags: `--max`, `--attempts`, `--seed` (0 = non-deterministic)

## Subproject template (each module should have)
- `README.md` — install, usage, tests/coverage, env vars, exit codes
- `go.mod` — `go 1.25`
- `Makefile` — targets: `deps` (tidy), `build` (only main pkg), `test`, `cover`, `clean`
- `.gitignore` — standard Go ignores
- `cmd/<binary>/` — CLI entrypoint (`package main`)
- `internal/` — implementation packages (no circular deps)

## CI
- **One workflow per subproject** under `.github/workflows/<name>.yml`.
- Use `on.push.paths` / `on.pull_request.paths` filters so only the touched subproject runs.
- Steps: checkout, setup-go (1.25.x), `go test ./... -coverprofile=cover.out`, upload artifacts.

## Adding a new subproject (checklist)
1) Create folder with module files: `README.md`, `.gitignore`, `Makefile`, `go.mod`, `internal/.gitkeep`, `cmd/<name>/.gitkeep`.
2) Update the root README subprojects list (and remove the project from the “Upcoming subprojects” section if present).
3) Add a dedicated workflow file for the new subproject.
4) Open one GitHub issue for the app with summary, v1 scope, non-goals, and planned milestones.
5) Break the work into focused PRs as needed, using tests-first sequencing where it helps.
6) Use the local `.codex-plan.md` file for detailed task breakdown, PR ordering, and active-session notes. Do not commit it.
7) Keep stdlib-only unless otherwise justified.

## Authoring prompts for Codex (guidelines)
- Always include: branch name, PR title, PR body (with `Refs #...` or `Closes #...` as appropriate), exact file paths and **full file contents** to write.
- When updating an existing PR: explicitly say **do NOT create a new branch/PR**, and specify which file(s) to replace.
- Keep prompts deterministic: avoid vague wording (“update as needed”).
- Keep each PR focused on one task; split follow-ups into later PRs under the same app issue when possible.

---
