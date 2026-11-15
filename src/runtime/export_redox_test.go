package runtime

import "internal/runtime/atomic"

type RedoxNote struct{ n note }

func (n *RedoxNote) Clear()             { noteclear(&n.n) }
func (n *RedoxNote) Wake()              { notewakeup(&n.n) }
func (n *RedoxNote) Wait(ns int64) bool { return notetsleepg(&n.n, ns) }

// RedoxSemaAccounting tests timeouts and wake tokens on a single M. Restore
// its original count even with a broken implementation so reporting is safe.
func RedoxSemaAccounting() (result [8]int64) {
	systemstack(func() {
		mp := getg().m
		saved := atomic.Load(&mp.waitsema)
		result[0] = int64(saved)
		if saved != 0 {
			return
		}
		result[1] = int64(semasleep(0))
		result[2] = int64(int32(atomic.Load(&mp.waitsema)))
		result[3] = int64(semasleep(0))
		result[4] = int64(int32(atomic.Load(&mp.waitsema)))
		semawakeup(mp)
		result[5] = int64(semasleep(0))
		result[6] = int64(int32(atomic.Load(&mp.waitsema)))
		result[7] = int64(semasleep(0))
		atomic.Store(&mp.waitsema, saved)
	})
	return
}
