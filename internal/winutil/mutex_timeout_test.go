package winutil

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// The lock name is guessable: a process that opens the named mutex first
// and never releases it must produce ErrLockTimeout instead of hanging
// every caller for ever.
func TestNamedMutex_TimeoutInsteadOfInfiniteWait(t *testing.T) {
	old := lockTimeout
	lockTimeout = 150 * time.Millisecond
	defer func() { lockTimeout = old }()

	name := `Local\VinPN-Test-Timeout-` + time.Now().Format("1500405.000000")
	m1, err := NewNamedMutex(name)
	require.NoError(t, err)
	m2, err := NewNamedMutex(name)
	require.NoError(t, err)

	require.NoError(t, m1.Lock()) // foreign holder keeps it
	start := time.Now()
	require.ErrorIs(t, m2.Lock(), ErrLockTimeout)
	require.Less(t, time.Since(start), 5*time.Second, "the wait must be bounded")

	// Released, the same waiter acquires it.
	require.NoError(t, m1.Unlock())
	require.NoError(t, m2.Lock())
	require.NoError(t, m2.Unlock())
}
