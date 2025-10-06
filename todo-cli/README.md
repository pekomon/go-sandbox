todo-cli

Local TODO manager (CLI).

Roadmap (next PRs)

- JSON persistence at ~/.todo-cli/tasks.json
- Commands: add, list, done <id>, rm <id>, clear
- Sorting: newest-first by default; flag to reverse
- Tests: JSON load/save, ID assignment, command parsing
- Error cases: corrupted JSON; concurrency note (simple lock or documented limitation)
- Standard library only (initially)
