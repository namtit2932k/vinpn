package main

import (
	goodbyedpi "github.com/sickyturtlez/vinpn/assets/goodbyedpi"
	zapret2 "github.com/sickyturtlez/vinpn/assets/zapret2"
	"github.com/sickyturtlez/vinpn/internal/dpi"
	"github.com/sickyturtlez/vinpn/internal/shell"
	"github.com/sickyturtlez/vinpn/internal/store"
)

func newDPIManager(paths store.Paths) *dpi.Manager { // headless only
	return shell.NewDPIManager(paths, goodbyedpi.FS, zapret2.FS, nil)
}

// stopDPI removes the WinDivert services left behind. An engine process
// owned by a dead VinPN is already gone: it lived in a kill-on-close job.
func stopDPI(paths store.Paths) func() error {
	return newDPIManager(paths).Stop
}
