package formats_test

import (
	"os"
	"testing"

	"github.com/sickyturtlez/vinpn/internal/rules/formats"
	"github.com/stretchr/testify/require"
)

func TestDetect_VinPN(t *testing.T) {
	data, err := os.ReadFile("testdata/vinpn.txt")
	require.NoError(t, err)
	for _, name := range []string{"x.txt", "list.yaml", "x.json", ""} {
		f, err := formats.Detect(name, data)
		require.NoError(t, err, name)
		require.Equal(t, formats.VinPN, f, name)
	}
}

func TestParse_VinPN(t *testing.T) {
	data, err := os.ReadFile("testdata/vinpn.txt")
	require.NoError(t, err)
	r, err := formats.Parse(formats.VinPN, data)
	require.NoError(t, err)
	require.Len(t, r.Entries, 2)
	require.Equal(t, 1, r.Skipped)
	require.Equal(t, 2, r.Counts["domain"])
	e := r.Entries[0]
	require.Equal(t, "youtube.com", e.Pattern.Value)
	require.Equal(t, 3, e.Line)
	require.NotNil(t, e.Action)
	require.Equal(t, "www.google.com", e.Action.SNI)
	require.Equal(t, "www.google.com", e.Action.Connect)
	require.True(t, r.Entries[1].Action.Block)
}
