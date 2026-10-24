package sync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Three devices edit the same concept from the same base and sync in turn.
// The first sync wins at the concept's path. The second conflicts and saves
// its edit at the fixed fork path ConflictLocalPath. The third conflicts too,
// and ApplyConflictFailsafe writes its edit to that same fork path, replacing
// the second device's edit, which had already been committed to the hub.

// appendLine mirrors an agent appending a line to a concept file.
func appendLine(t *testing.T, dir, rel, line string) {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(rel))
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}
	if err := os.WriteFile(p, append(data, []byte(line+"\n")...), 0o644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
}

// Fails today: device B's committed edit is gone from the hub's current tree.
func TestHub_ThreeWritersKeepEveryCommittedEdit(t *testing.T) {
	f, dirA, dirB := twoDevicesFromSeed(t, "decisions/seed")
	dirC := t.TempDir()
	f.pull(dirC)

	edits := map[string]string{"A": "edit-from-A", "B": "edit-from-B", "C": "edit-from-C"}
	appendLine(t, dirA, "decisions/seed.md", edits["A"])
	appendLine(t, dirB, "decisions/seed.md", edits["B"])
	appendLine(t, dirC, "decisions/seed.md", edits["C"])
	f.sync(dirA, "A edits seed")
	f.sync(dirB, "B edits seed")
	f.sync(dirC, "C edits seed")

	// A fresh device sees exactly the hub's current tree.
	dirD := t.TempDir()
	f.pull(dirD)
	var current strings.Builder
	_ = filepath.Walk(dirD, func(p string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasPrefix(info.Name(), "seed") {
			data, _ := os.ReadFile(p)
			current.Write(data)
		}
		return nil
	})
	for name, edit := range edits {
		if !strings.Contains(current.String(), edit) {
			t.Errorf("device %s's committed edit %q is missing from the hub's current tree (files: %v)",
				name, edit, conflictFiles(t, dirD))
		}
	}
}

// Control: with two writers the fork path is written once, so both edits
// survive (one at the path, one in the fork), isolating the third writer as
// the cause above.
func TestHub_TwoWritersKeepEveryCommittedEdit(t *testing.T) {
	f, dirA, dirB := twoDevicesFromSeed(t, "decisions/seed")

	appendLine(t, dirA, "decisions/seed.md", "edit-from-A")
	appendLine(t, dirB, "decisions/seed.md", "edit-from-B")
	f.sync(dirA, "A edits seed")
	f.sync(dirB, "B edits seed")

	dirC := t.TempDir()
	f.pull(dirC)
	current := readFile(t, dirC, "decisions/seed.md") + readFile(t, dirC, "decisions/seed.conflict-local.md")
	for _, edit := range []string{"edit-from-A", "edit-from-B"} {
		if !strings.Contains(current, edit) {
			t.Errorf("committed edit %q is missing from the hub's current tree", edit)
		}
	}
}
