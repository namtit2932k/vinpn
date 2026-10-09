package startup_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sickyturtlez/vinpn/internal/startup"
	"github.com/stretchr/testify/require"
)

// A path under an admin-controlled root is used as it is: no copy, no stat
// (the path may not exist yet in a unit test).
func TestSafeExe_KeepsProgramFilesImage(t *testing.T) {
	if os.Getenv("ProgramFiles") == "" {
		t.Skip("ProgramFiles not set")
	}
	machineDir := t.TempDir()
	exe := filepath.Join(os.Getenv("ProgramFiles"), "VinPN", "vinpn.exe")
	got, err := startup.SafeExe(exe, machineDir)
	require.NoError(t, err)
	require.Equal(t, exe, got)
}

// The staged copy itself lives under the machine dir, so a run launched
// from it (autostart task) passes the same check on its next connect.
func TestSafeExe_AcceptsItsOwnStagedCopy(t *testing.T) {
	machineDir := t.TempDir()
	staged := filepath.Join(machineDir, "vinpn.exe")
	require.NoError(t, os.WriteFile(staged, []byte("img"), 0o600))
	got, err := startup.SafeExe(staged, machineDir)
	require.NoError(t, err)
	require.Equal(t, staged, got)
}

// An image a normal user can replace (portable run, repo checkout) is
// copied into the admin-only machine directory before any task or the
// watchdog runs it with elevation, and the copy follows rebuilds.
func TestSafeExe_StagesUserWritableImage(t *testing.T) {
	srcDir := t.TempDir()
	src := filepath.Join(srcDir, "vinpn.exe")
	require.NoError(t, os.WriteFile(src, []byte("build-one"), 0o644))
	machineDir := filepath.Join(t.TempDir(), "m")

	dst, err := startup.SafeExe(src, machineDir)
	require.NoError(t, err)
	require.Equal(t, filepath.Join(machineDir, "vinpn.exe"), dst)
	b, err := os.ReadFile(dst)
	require.NoError(t, err)
	require.Equal(t, "build-one", string(b))

	// Second call is idempotent while the source is unchanged.
	got, err := startup.SafeExe(src, machineDir)
	require.NoError(t, err)
	require.Equal(t, dst, got)

	// A different build refreshes the staged copy.
	require.NoError(t, os.WriteFile(src, []byte("build-two-different-size"), 0o644))
	require.NoError(t, os.Chtimes(src, time.Now().Add(time.Minute), time.Now().Add(time.Minute)))
	got, err = startup.SafeExe(src, machineDir)
	require.NoError(t, err)
	require.Equal(t, dst, got)
	b, err = os.ReadFile(dst)
	require.NoError(t, err)
	require.Equal(t, "build-two-different-size", string(b))
}

// A missing source is an error, never a silent fallback to the original
// (unstaged) path.
func TestSafeExe_MissingSourceFails(t *testing.T) {
	_, err := startup.SafeExe(filepath.Join(t.TempDir(), "ghost", "vinpn.exe"), t.TempDir())
	require.Error(t, err)
}
