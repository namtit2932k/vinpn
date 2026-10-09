package startup_test

import (
	"encoding/xml"
	"io"
	"strings"
	"testing"

	"github.com/sickyturtlez/vinpn/internal/brand"
	"github.com/sickyturtlez/vinpn/internal/startup"
	"github.com/stretchr/testify/require"
)

func TestTaskXML(t *testing.T) {
	x := startup.TaskXML(startup.Task{Name: "VinPN", Exe: `C:\Program Files\VinPN\vinpn.exe`,
		Args: "--autostart", UserID: `PC\Đức`, TimeLimit: "PT0S"})
	require.Contains(t, x, "<RunLevel>HighestAvailable</RunLevel>")
	require.Contains(t, x, "<LogonTrigger>")
	require.Contains(t, x, `<Command>C:\Program Files\VinPN\vinpn.exe</Command>`)
	require.Contains(t, x, "<Arguments>--autostart</Arguments>")
	require.Contains(t, x, `<UserId>PC\Đức</UserId>`)
	require.Contains(t, x, "<ExecutionTimeLimit>PT0S</ExecutionTimeLimit>")
	require.Contains(t, x, "<DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>")
	require.Contains(t, x, "<MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>")
	var v any
	dec := xml.NewDecoder(strings.NewReader(x))
	dec.CharsetReader = func(_ string, r io.Reader) (io.Reader, error) { return r, nil } // in-memory string is UTF-8
	require.NoError(t, dec.Decode(&v))
}

func TestTaskXML_EscapesSpecialChars(t *testing.T) {
	x := startup.TaskXML(startup.Task{Name: "x", Exe: `C:\A&B\<g>.exe`, UserID: "u", TimeLimit: "PT0S"})
	require.Contains(t, x, `C:\A&amp;B\&lt;g&gt;.exe`)
}

func TestRecoveryAndAutostartTasks(t *testing.T) {
	r := startup.RecoveryTask(`C:\g.exe`, "")
	require.Equal(t, brand.TaskRecovery, r.Name)
	require.Equal(t, "--restore", r.Args)
	require.Equal(t, "PT5M", r.TimeLimit)
	a := startup.AutostartTask(`C:\g.exe`, "")
	require.Equal(t, brand.TaskAutostart, a.Name)
	require.Equal(t, "--autostart", a.Args)
	require.Equal(t, "PT0S", a.TimeLimit)
	require.NotEmpty(t, a.UserID)
}

// A staged copy of the executable lives outside the profile, so the task
// carries the original data directory (quoted: Task Scheduler splits the
// argument string on spaces).
func TestTasks_CarryDataDir(t *testing.T) {
	s := startup.RecoveryTask(`C:\ProgramData\VinPN\vinpn.exe`, `C:\Users\Đức\AppData\Roaming\VinPN`)
	require.Equal(t, `--restore --data-dir "C:\Users\Đức\AppData\Roaming\VinPN"`, s.Args)
	u := startup.AutostartTask(`C:\ProgramData\VinPN\vinpn.exe`, `D:\data`)
	require.Equal(t, `--autostart --data-dir "D:\data"`, u.Args)
}

func TestUTF16WithBOM(t *testing.T) {
	b := startup.EncodeUTF16LE("<a/>")
	require.Equal(t, []byte{0xFF, 0xFE, '<', 0, 'a', 0, '/', 0, '>', 0}, b)
}
