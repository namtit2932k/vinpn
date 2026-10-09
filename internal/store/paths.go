// Package store owns VinPN's on-disk data: paths, settings and state.
package store

import (
	"os"
	"path/filepath"

	"github.com/sickyturtlez/vinpn/internal/brand"
)

// Paths lists every file VinPN keeps in its data directory.
type Paths struct {
	DataDir            string
	Settings           string
	State              string
	Meta               string
	ScanCache          string
	CFScanCache        string
	ServersRemote      string
	ServersRemoteSig   string
	ServersDNSCrypt    string
	ServersDNSCryptSig string
	ServersCustom      string
	DPIBlacklist       string
	DPIStrategies      string
	DPIStrategiesSig   string
	DPIAutoHostlist    string
	LogDir             string
	BinDir             string
	Rules              string
	FragCache          string
	ListsDir           string
	// LAN CA files live in MachineDir (see WithMachineDir); until it is
	// set they default to DataDir, which tests use.
	MachineDir string
	LANCACert  string
	LANCAKey   string
	Portable   bool
}

// ResolvePaths picks the data directory: <exe dir>\data when a "portable"
// marker sits next to the executable, otherwise %APPDATA%\VinPN.
func ResolvePaths(exePath, appData string) Paths {
	exeDir := filepath.Dir(exePath)
	p := Paths{}
	if _, err := os.Stat(filepath.Join(exeDir, "portable")); err == nil {
		p.Portable = true
		return WithDataDir(p, filepath.Join(exeDir, "data"))
	}
	return WithDataDir(p, filepath.Join(appData, brand.AppName))
}

// DefaultMachineDir is the machine-wide directory (%ProgramData%\VinPN).
// winutil.SecureDir restricts it to SYSTEM and Administrators.
func DefaultMachineDir() string {
	return filepath.Join(os.Getenv("ProgramData"), brand.AppName)
}

// WithDataDir points every data-directory file at dir. A scheduled task or
// the watchdog runs a staged copy of the executable (startup.SafeExe), so
// the spawner passes the data directory it resolved on its own command
// line; without it the staged copy would read a different profile.
func WithDataDir(p Paths, dir string) Paths {
	p.DataDir = dir
	j := func(name string) string { return filepath.Join(dir, name) }
	p.Settings = j("settings.json")
	p.State = j("state.json")
	p.LANCACert = j("lan-ca.crt")
	p.LANCAKey = j("lan-ca.key")
	p.Meta = j("meta.json")
	p.ScanCache = j("scan-cache.json")
	p.CFScanCache = j("cfscan-cache.json")
	p.ServersRemote = j("servers-remote.json")
	p.ServersRemoteSig = j("servers-remote.json.sig")
	p.ServersDNSCrypt = j("servers-dnscrypt.md")
	p.ServersDNSCryptSig = j("servers-dnscrypt.md.minisig")
	p.ServersCustom = j("servers-custom.json")
	p.DPIBlacklist = j("dpi-blacklist.txt")
	p.DPIStrategies = j("dpi-strategies.json")
	p.DPIStrategiesSig = j("dpi-strategies.json.sig")
	p.DPIAutoHostlist = j("dpi-autohostlist.txt")
	p.LogDir = j("logs")
	p.BinDir = j("bin")
	p.Rules = j("rules.json")
	p.FragCache = j("frag-cache.json")
	p.ListsDir = j("lists")
	// Files WithMachineDir moved out of the user-writable tree stay there,
	// whatever the call order.
	machineFiles(&p)
	return p
}

// machineFiles re-applies the machine-dir files when MachineDir is set, so
// WithDataDir and WithMachineDir commute.
func machineFiles(p *Paths) {
	if p.MachineDir == "" {
		return
	}
	p.State = filepath.Join(p.MachineDir, "state.json")
	p.BinDir = filepath.Join(p.MachineDir, "bin")
	p.LANCACert = filepath.Join(p.MachineDir, "lan-ca.crt")
	p.LANCAKey = filepath.Join(p.MachineDir, "lan-ca.key")
}

// WithMachineDir moves the LAN CA files, state.json and the engine
// binaries to dir, a machine-wide directory only SYSTEM and Administrators
// can open (%ProgramData%\VinPN). Three payloads justify the move: the CA
// key ends up in the Root store, state.json drives what the elevated
// --restore task writes into system DNS and the proxy settings, and the
// engine binaries are executed with elevation — none of them may live in a
// directory the current user can rewrite.
func WithMachineDir(p Paths, dir string) Paths {
	p.MachineDir = dir
	machineFiles(&p)
	return p
}
