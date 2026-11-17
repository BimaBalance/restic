//go:build windows

package restorer

import (
	"os"
	"path/filepath"
	"testing"

	rtest "github.com/restic/restic/internal/test"
)

func TestWindowsCreateFileDoesNotFollowExternalSymlink(t *testing.T) {
	root := t.TempDir()

	outside := filepath.Join(root, "outside.txt")
	restoreDir := filepath.Join(root, "restore")
	link := filepath.Join(restoreDir, "target.txt")

	rtest.OK(t, os.Mkdir(restoreDir, 0700))

	original := []byte("DO NOT MODIFY THIS FILE")
	rtest.OK(t, os.WriteFile(outside, original, 0600))

	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("cannot create symlink: %v", err)
	}

	f, err := createFile(link, 5, false, false)
	rtest.OK(t, err)
	rtest.OK(t, f.Close())

	got, err := os.ReadFile(outside)
	rtest.OK(t, err)
	rtest.Equals(t, original, got, "external file was modified through symlink")
}
