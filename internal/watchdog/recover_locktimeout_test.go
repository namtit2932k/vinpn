package watchdog_test

import (
	"testing"

	"github.com/sickyturtlez/vinpn/internal/store"
	"github.com/sickyturtlez/vinpn/internal/watchdog"
	"github.com/sickyturtlez/vinpn/internal/winutil"
	"github.com/stretchr/testify/require"
)

// stuckLock models a foreign process that opened Local\VinPN-State first
// and never releases it.
type stuckLock struct{}

func (stuckLock) Lock() error   { return winutil.ErrLockTimeout }
func (stuckLock) Unlock() error { return nil }

// A guessable lock name must not become a DNS outage: when the acquire
// times out, recovery runs without the cross-process lock instead of
// leaving every adapter on loopback.
func TestRestore_LockTimeoutStillRestores(t *testing.T) {
	d, fd, stops := setup(t, false, dirty())
	d.States = store.NewStateStore(d.States.Path(), stuckLock{})
	out, err := watchdog.RestoreIfOrphaned(d)
	require.NoError(t, err)
	require.Equal(t, watchdog.Restored, out)
	require.Len(t, fd.restored, 1)
	require.Equal(t, 1, *stops)
	st, err := d.States.Load()
	require.NoError(t, err)
	require.Equal(t, store.PhaseClean, st.Phase)
}
