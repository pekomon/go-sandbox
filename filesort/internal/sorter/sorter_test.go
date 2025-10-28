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

// Expected classes (to be implemented in #11):
// images: .jpg .jpeg .png .gif
// docs:   .pdf .doc .docx .txt .md
// videos: .mp4 .mov .avi
// other:  everything else
// Destination layout: <root>/<class>/<filename>

func TestBuildPlan_DryRunAndApply_BasicLayout(t *testing.T) {
	root := t.TempDir()

	// Seed files in the root
	img1 := touch(t, root, "photo.JPG")
	img2 := touch(t, root, "pic.png")
	doc1 := touch(t, root, "notes.md")
	vid1 := touch(t, root, "clip.MP4")
	oth1 := touch(t, root, "archive.tar.gz")

	// Build plan (should enumerate moves but not touch FS when dryRun=true)
	p, err := sorter.BuildPlan(root, true)
	if err == nil {
		// Because stubs return ErrNotImplemented, we expect err != nil until implementation PR.
		t.Fatalf("expected failure (not implemented), got nil")
	}
	_ = p

	// The following assertions describe the *intended* behavior and will be enabled in #11.
	// They are kept here so the test body is self-documenting.
	_ = img1
	_ = img2
	_ = doc1
	_ = vid1
	_ = oth1

	// After BuildPlan with dryRun=true, on-disk layout must be unchanged:
	gotRoot := listDir(t, root)
	// Expect only the original files present (order sorted for stability)
	wantRoot := []string{"archive.tar.gz", "clip.MP4", "notes.md", "photo.JPG", "pic.png"}
	if strings.Join(gotRoot, ",") != strings.Join(wantRoot, ",") {
		t.Logf("GOT:  %v", gotRoot)
		t.Logf("WANT: %v", wantRoot)
		// Uncomment once implemented:
		// t.Fatalf("unexpected root layout after dry-run")
	}
}

func TestApply_ExecutesMoves(t *testing.T) {
	root := t.TempDir()

	_ = touch(t, root, "photo.jpg")
	_ = touch(t, root, "doc.txt")
	_ = touch(t, root, "movie.mp4")
	_ = touch(t, root, "README") // other

	// Build a plan without dry-run, then apply it.
	p, err := sorter.BuildPlan(root, false)
	if err == nil {
		// Expect failure until #11 implements it.
		t.Fatalf("expected failure (not implemented), got nil")
	}

	// When implemented, Apply should create folders and move files:
	// err = sorter.Apply(p)
	// if err != nil { t.Fatalf("apply: %v", err) }

	// wantDirs := []string{"docs", "images", "other", "videos"}
	// gotDirs := listDir(t, root)
	// slices.Sort(wantDirs)
	// if strings.Join(gotDirs, ",") != strings.Join(wantDirs, ",") {
	// t.Fatalf("unexpected dirs in root: %v", gotDirs)
	// }
}
