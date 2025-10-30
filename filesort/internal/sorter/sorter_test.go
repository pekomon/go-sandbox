package sorter_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/pekomon/go-sandbox/filesort/internal/sorter"
)

func touch(t *testing.T, dir, name string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

func listDir(t *testing.T, dir string) []string {
	t.Helper()
	ents, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("readdir %s: %v", dir, err)
	}
	var names []string
	for _, e := range ents {
		names = append(names, e.Name())
	}
	slices.Sort(names)
	return names
}

func TestBuildPlan_DryRunAndApply_BasicLayout(t *testing.T) {
	root := t.TempDir()

	// Seed files in the root
	img1 := touch(t, root, "photo.JPG")
	img2 := touch(t, root, "pic.png")
	doc1 := touch(t, root, "notes.md")
	vid1 := touch(t, root, "clip.MP4")
	oth1 := touch(t, root, "archive.tar.gz")

	// Build plan (must not touch FS when dryRun=true)
	p, err := sorter.BuildPlan(root, true)
	if err != nil {
		t.Fatalf("build plan: %v", err)
	}

	// On-disk layout unchanged:
	gotRoot := listDir(t, root)
	wantRoot := []string{"archive.tar.gz", "clip.MP4", "notes.md", "photo.JPG", "pic.png"}
	if strings.Join(gotRoot, ",") != strings.Join(wantRoot, ",") {
		t.Fatalf("unexpected root layout after dry-run\nGOT:  %v\nWANT: %v", gotRoot, wantRoot)
	}

	// Sanity check: planned moves include the expected destinations
	if len(p.Moves) != 5 {
		t.Fatalf("expected 5 planned moves, got %d", len(p.Moves))
	}
	_ = img1
	_ = img2
	_ = doc1
	_ = vid1
	_ = oth1
}

func TestApply_ExecutesMoves(t *testing.T) {
	root := t.TempDir()

	_ = touch(t, root, "photo.jpg") // images
	_ = touch(t, root, "doc.txt")   // docs
	_ = touch(t, root, "movie.mp4") // videos
	_ = touch(t, root, "README")    // other

	p, err := sorter.BuildPlan(root, false)
	if err != nil {
		t.Fatalf("build plan: %v", err)
	}
	if err := sorter.Apply(p); err != nil {
		t.Fatalf("apply: %v", err)
	}

	gotDirs := listDir(t, root)
	wantDirs := []string{"docs", "images", "other", "videos"}
	slices.Sort(wantDirs)
	if strings.Join(gotDirs, ",") != strings.Join(wantDirs, ",") {
		t.Fatalf("unexpected dirs in root: %v", gotDirs)
	}

	// spot-check files exist under expected dirs
	if _, err := os.Stat(filepath.Join(root, "images", "photo.jpg")); err != nil {
		t.Fatalf("missing moved image: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "docs", "doc.txt")); err != nil {
		t.Fatalf("missing moved doc: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "videos", "movie.mp4")); err != nil {
		t.Fatalf("missing moved video: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "other", "README")); err != nil {
		t.Fatalf("missing moved other: %v", err)
	}
}
