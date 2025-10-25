# filesort

CLI tool to classify files in a directory into subfolders by type (by file extension).  
Standard library only (initially). Tests will operate on a temporary directory and support a `--dry-run` mode.

**Planned features (future PRs):**
- Detect by extension using pure stdlib
- `--dry-run` shows the planned moves without touching the filesystem
- Clear reports of what moved where
- Offline tests using temp dirs
