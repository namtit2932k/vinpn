package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/sickyturtlez/vinpn/internal/app"
	"github.com/sickyturtlez/vinpn/internal/brand"
	"github.com/sickyturtlez/vinpn/internal/certs"
	"github.com/sickyturtlez/vinpn/internal/certstore"
	"github.com/sickyturtlez/vinpn/internal/cli"
	"github.com/sickyturtlez/vinpn/internal/logx"
	"github.com/sickyturtlez/vinpn/internal/store"
	"github.com/sickyturtlez/vinpn/internal/sysdns"
	"github.com/sickyturtlez/vinpn/internal/sysproxy"
	"github.com/sickyturtlez/vinpn/internal/watchdog"
	"github.com/sickyturtlez/vinpn/internal/winutil"
)

// runHeadless handles --watchdog, --restore, --remove-certs and --export. It never touches Wails.
func runHeadless(mode cli.Mode) int {
	exe, err := os.Executable()
	if err != nil {
		return 1
	}
	paths := store.WithMachineDir(store.ResolvePaths(exe, os.Getenv("APPDATA")), store.DefaultMachineDir())
	if mode.DataDir != "" {
		paths = store.WithDataDir(paths, mode.DataDir)
	}
	// state.json and the LAN CA key live here; without the admin-only DACL
	// a normal user could hand --restore a forged DNS snapshot. Refuse to
	// run recovery rather than trust an unsecurable directory.
	if err := winutil.SecureDir(paths.MachineDir); err != nil {
		fmt.Fprintln(os.Stderr, brand.AppName+": securing "+paths.MachineDir+": "+err.Error())
		return 1
	}
	logger := slog.Default()
	if w, err := logx.NewRotating(paths.LogDir, "vinpn", 5<<20, 3); err == nil {
		defer w.Close()
		logger = slog.New(slog.NewTextHandler(w, nil)).With("mode", modeName(mode.Kind))
	}
	if mode.Kind == cli.KindExport {
		// Read-only: no state lock, no recovery. A GUI exe has no console of
		// its own: borrow the caller's so the result can be read.
		winutil.AttachParentConsole()
		if err := app.ExportTo(paths, mode.ExportPath, brand.Version); err != nil {
			logger.Error("export failed", "err", err)
			fmt.Fprintln(os.Stderr, "VinPN: export failed:", err)
			return 1
		}
		_, _ = fmt.Fprintln(os.Stdout, "VinPN: settings exported to", mode.ExportPath)
		return 0
	}
	lock, err := winutil.NewNamedMutex(brand.StateMutex)
	if err != nil {
		logger.Error("state mutex", "err", err)
		return 1
	}
	roots := certstore.NewWindows(certstore.LocalMachine)
	d := watchdog.Deps{
		States:  store.NewStateStore(paths.State, lock),
		DNS:     sysdns.NewManager(sysdns.NewWindowsAPI(), time.Sleep),
		StopDPI: stopDPI(paths),
		Alive:   winutil.ProcessAlive,
		Log:     logger,

		RestoreSysProxy: sysproxy.Manager{API: sysproxy.NewWindowsAPI()}.RestoreIfOurs,
		DeleteRule:      winutil.DeleteNamedRule,
		RemoveCert: func(t string) error {
			// state.json is user-writable: remove only Fake SNI roots.
			return certstore.RemoveIfPrefix(roots, t, certs.SessionPrefix)
		},
		SweepSession: sweepSession(roots),
	}
	switch mode.Kind {
	case cli.KindWatchdog:
		err = watchdog.RunWatchdog(mode.ParentPID, mode.ParentStart, winutil.WaitForExit, d)
	case cli.KindRemoveCerts:
		// The uninstaller: undo whatever a run left, then remove every
		// VinPN root and the LAN CA files.
		err = errors.Join(watchdog.RunRestore(d), watchdog.RemoveAllCerts(roots, paths.LANCACert, paths.LANCAKey))
	default:
		err = watchdog.RunRestore(d)
	}
	if err != nil {
		logger.Error("recovery failed", "err", err)
		return 1
	}
	return 0
}

func modeName(k cli.Kind) string {
	switch k {
	case cli.KindWatchdog:
		return "watchdog"
	case cli.KindRestore:
		return "restore"
	case cli.KindRemoveCerts:
		return "remove-certs"
	case cli.KindAutostart:
		return "autostart"
	case cli.KindExport:
		return "export"
	}
	return "ui"
}

// sweepSession removes Fake SNI roots not in keep.
func sweepSession(s certstore.Store) func(keep []string) error {
	return func(keep []string) error {
		_, err := certstore.Sweep(s, certs.SessionPrefix, keep)
		return err
	}
}
