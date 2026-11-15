package runtime_test

import (
	"runtime"
	"testing"
)

func TestRedoxSemaTimeoutAccounting(t *testing.T) {
	want := [8]int64{0, -1, 0, -1, 0, 0, 0, -1}
	for i := 0; i < 32; i++ {
		if got := runtime.RedoxSemaAccounting(); got != want {
			t.Fatalf("semaphore accounting iteration %d: got %v, want %v", i, got, want)
		}
	}
}
