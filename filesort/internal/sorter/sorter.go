package sorter

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrNotImplemented = errors.New("not implemented")

// Class represents a logical destination folder.
type Class string

const (
	ClassImages Class = "images"
	ClassDocs   Class = "docs"
	ClassVideos Class = "videos"
	ClassOther  Class = "other"
)

// Plan is a set of moves from absolute source to absolute destination (both under Root).
type Plan struct {
	Root  string
	Moves map[string]string // srcAbs -> dstAbs
}

// BuildPlan analyzes files under root (non-recursive) and computes destination moves.
// When dryRun is true, the filesystem must remain untouched (this function only returns the plan).
func BuildPlan(root string, dryRun bool) (Plan, error) {
	if root == "" {
		return Plan{}, fmt.Errorf("root is required")
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return Plan{}, err
	}
	info, err := os.Stat(absRoot)
	if err != nil {
		return Plan{}, err
	}
	if !info.IsDir() {
		return Plan{}, fmt.Errorf("not a directory: %s", absRoot)
	}

	ents, err := os.ReadDir(absRoot)
	if err != nil {
		return Plan{}, err
	}

	moves := make(map[string]string)
	for _, e := range ents {
		if e.IsDir() {
			// Non-recursive by design (future PR could add recursion).
			continue
		}
		name := e.Name()
		src := filepath.Join(absRoot, name)

		cl := classifyByExt(name)
		dstDir := filepath.Join(absRoot, string(cl))
		dst := filepath.Join(dstDir, name)

		// Skip no-op moves (e.g., already in place, though this shouldn't happen for root files).
		if src == dst {
			continue
		}
		moves[src] = dst
	}

	return Plan{
		Root:  absRoot,
		Moves: moves,
	}, nil
}

// Apply executes the plan: create destination dirs and move files with os.Rename.
func Apply(p Plan) error {
	for src, dst := range p.Moves {
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		if err := os.Rename(src, dst); err != nil {
			return err
		}
	}
	return nil
}

// classifyByExt maps filename to a class by (last) extension, case-insensitively.
func classifyByExt(name string) Class {
	ext := strings.ToLower(filepath.Ext(name)) // uses only the last extension (e.g., .gz)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return ClassImages
	case ".pdf", ".doc", ".docx", ".txt", ".md":
		return ClassDocs
	case ".mp4", ".mov", ".avi":
		return ClassVideos
	default:
		return ClassOther
	}
}
