package sorter

import "errors"

var ErrNotImplemented = errors.New("not implemented")

// Class represents a logical destination folder, e.g. "images", "docs", "videos", "other".
type Class string

// Plan is a set of moves from absolute source to absolute destination (both must be under the root dir).
type Plan struct {
	Root  string
	Moves map[string]string // srcAbs -> dstAbs
}

// BuildPlan analyzes files under root (non-recursive for now) and computes destination moves.
// If dryRun is true, the plan should still be built fully; only Apply would refrain from changing the FS.
// STUB for now: return ErrNotImplemented so tests fail.
func BuildPlan(root string, dryRun bool) (Plan, error) {
	return Plan{}, ErrNotImplemented
}

// Apply executes the plan: create destination dirs as needed and move files.
// STUB for now: return ErrNotImplemented so tests fail.
func Apply(p Plan) error {
	return ErrNotImplemented
}
