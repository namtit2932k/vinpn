package startup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SafeExe returns the executable a scheduled task or the detached watchdog
// may run with elevation. An image under an admin-controlled program root
// (Program Files, Windows, the machine directory) is used as it is;
// anything else — a portable run from Downloads, a repo checkout — is first
// copied into machineDir, which SecureDir restricts to SYSTEM and
// Administrators. HighestAvailable tasks must never point at a file a
// normal user can replace, or local code swaps itself into admin at logon.
func SafeExe(exe, machineDir string) (string, error) {
	if exe == "" || machineDir == "" {
		return "", fmt.Errorf("startup: SafeExe: empty %q %q", exe, machineDir)
	}
	if underAdminRoot(exe, machineDir) {
		return exe, nil
	}
	dst := filepath.Join(machineDir, filepath.Base(exe))
	if err := stage(exe, dst); err != nil {
		return "", err
	}
	return dst, nil
}

// underAdminRoot reports whether path sits under a directory only
// administrators can rewrite: the machine directory itself, Program Files,
// Program Files (x86) or %windir%.
func underAdminRoot(path, machineDir string) bool {
	p := cleanPath(path)
	for _, r := range []string{
		machineDir,
		os.Getenv("ProgramFiles"),
		os.Getenv("ProgramFiles(x86)"),
		os.Getenv("windir"),
	} {
		if r == "" {
			continue
		}
		root := cleanPath(r)
		if p == root || strings.HasPrefix(p, root+`\`) {
			return true
		}
	}
	return false
}

func cleanPath(p string) string {
	return strings.TrimSuffix(strings.ToLower(filepath.Clean(p)), `\`)
}

// stage copies src to dst, refreshing the copy when src is a different
// build (size or modification time moved). A locked destination — the
// staged copy itself is running — keeps the previous image: it is still an
// admin-owned VinPN binary, and failing the connect over a stale watchdog
// would be worse.
func stage(src, dst string) error {
	si, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("startup: stage %s: %w", src, err)
	}
	if di, err := os.Stat(dst); err == nil && di.Size() == si.Size() && !di.ModTime().Before(si.ModTime()) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o700); err != nil {
		return err
	}
	b, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("startup: stage %s: %w", src, err)
	}
	tmp := dst + ".new"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	// The copy carries the source's mtime so the check above stays exact.
	_ = os.Chtimes(tmp, time.Now(), si.ModTime())
	if err := os.Rename(tmp, dst); err != nil {
		_ = os.Remove(tmp)
		if _, statErr := os.Stat(dst); statErr == nil {
			return nil // dst exists but is locked (running): keep it
		}
		return fmt.Errorf("startup: stage %s: %w", src, err)
	}
	return nil
}
