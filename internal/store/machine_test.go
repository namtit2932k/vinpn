package store_test

import (
	"path/filepath"
	"testing"

	"github.com/sickyturtlez/vinpn/internal/store"
	"github.com/stretchr/testify/require"
)

func TestWithMachineDir(t *testing.T) {
	p := store.ResolvePaths(`C:\Program Files\VinPN\vinpn.exe`, `C:\Users\u\AppData\Roaming`)
	m := store.WithMachineDir(p, `C:\ProgramData\VinPN`)
	require.Equal(t, filepath.Join(`C:\ProgramData\VinPN`, "lan-ca.crt"), m.LANCACert)
	require.Equal(t, filepath.Join(`C:\ProgramData\VinPN`, "lan-ca.key"), m.LANCAKey)
	require.Equal(t, `C:\ProgramData\VinPN`, m.MachineDir)
	// state.json and the engine binaries must leave the user-writable data
	// directory: the elevated restore task and engine processes act on them.
	require.Equal(t, filepath.Join(`C:\ProgramData\VinPN`, "state.json"), m.State)
	require.Equal(t, filepath.Join(`C:\ProgramData\VinPN`, "bin"), m.BinDir)
	require.NotEqual(t, p.State, m.State)
}

func TestWithDataDir(t *testing.T) {
	p := store.ResolvePaths(`C:\Program Files\VinPN\vinpn.exe`, `C:\Users\u\AppData\Roaming`)
	require.False(t, p.Portable)
	require.Equal(t, filepath.Join(`C:\Users\u\AppData\Roaming`, "VinPN", "state.json"), p.State)
	// The spawner re-points a staged copy at its data directory.
	q := store.WithDataDir(p, `D:\data`)
	require.Equal(t, `D:\data`, q.DataDir)
	require.Equal(t, filepath.Join(`D:\data`, "settings.json"), q.Settings)
	require.Equal(t, filepath.Join(`D:\data`, "bin"), q.BinDir)
	require.Equal(t, store.WithMachineDir(q, `C:\ProgramData\VinPN`).State,
		filepath.Join(`C:\ProgramData\VinPN`, "state.json"))
	// The calls commute: moving the data dir never drags state.json or the
	// engine binaries back into the user-writable tree.
	m := store.WithMachineDir(q, `C:\ProgramData\VinPN`)
	require.Equal(t, filepath.Join(`C:\ProgramData\VinPN`, "state.json"), store.WithDataDir(m, `E:\other`).State)
	require.Equal(t, filepath.Join(`C:\ProgramData\VinPN`, "bin"), store.WithDataDir(m, `E:\other`).BinDir)
}
