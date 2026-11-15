package runtime_test

import (
	"runtime"
	"testing"
	"time"
)

func TestRedoxSemaTimeoutAccounting(t *testing.T) {
	want := [8]int64{0, -1, 0, -1, 0, 0, 0, -1}
	for i := 0; i < 32; i++ {
		if got := runtime.RedoxSemaAccounting(); got != want {
			t.Fatalf("semaphore accounting iteration %d: got %v, want %v", i, got, want)
		}
	}
}

func TestRedoxNoteTimeoutWake(t *testing.T) {
	var n runtime.RedoxNote
	for i := 0; i < 64; i++ {
		if n.Wait(0) {
			t.Fatalf("iteration %d: unwoken note reported wake", i)
		}
		done := make(chan struct{})
		go func() {
			time.Sleep(time.Duration(i%3) * time.Millisecond)
			n.Wake()
			close(done)
		}()
		// Either result is valid when the timeout races with the wake.
		n.Wait(int64(time.Millisecond))
		<-done
		if !n.Wait(0) {
			t.Fatalf("iteration %d: completed wake was lost", i)
		}
		n.Clear()
		if n.Wait(0) {
			t.Fatalf("iteration %d: stale wake after clear", i)
		}
		n.Wake()
		if !n.Wait(0) {
			t.Fatalf("iteration %d: wake before wait was lost", i)
		}
		n.Clear()
	}
}
