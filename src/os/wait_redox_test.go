// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_test

import (
	"errors"
	"os"
	"runtime"
	"syscall"
	"testing"
)

func TestRedoxWaitHelper(t *testing.T) {
	if os.Getenv("GO_REDOX_WAIT_HELPER") != "1" {
		return
	}
	var b [1]byte
	if n, err := os.Stdin.Read(b[:]); n != 1 || err != nil {
		os.Exit(38)
	}
	os.Exit(37)
}

func startRedoxWaitChild(t *testing.T) (*os.Process, *os.File) {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	p, err := os.StartProcess(exe, []string{exe, "-test.run=^TestRedoxWaitHelper$"}, &os.ProcAttr{
		Env:   append(os.Environ(), "GO_REDOX_WAIT_HELPER=1"),
		Files: []*os.File{r, os.Stdout, os.Stderr},
	})
	r.Close()
	if err != nil {
		w.Close()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		w.Close()
		p.Kill()
		p.Wait()
	})
	return p, w
}

func TestRedoxWaitBeforeReap(t *testing.T) {
	p, gate := startRedoxWaitChild(t)
	if _, err := gate.Write([]byte{1}); err != nil {
		t.Fatal(err)
	}
	gate.Close()
	for i := 0; i < 2; i++ {
		if ready, err := os.BlockUntilWaitable(p); !ready || err != nil {
			t.Fatalf("observation %d: ready=%v err=%v", i, ready, err)
		}
	}
	state, err := p.Wait()
	if err != nil || state.ExitCode() != 37 {
		t.Fatalf("Wait: state=%v err=%v", state, err)
	}
	var status syscall.WaitStatus
	if _, err := syscall.Wait4(p.Pid, &status, syscall.WNOHANG, nil); err != syscall.ECHILD {
		t.Fatalf("wait after reap: %v, want ECHILD", err)
	}
	if err := p.Signal(syscall.Signal(0)); !errors.Is(err, os.ErrProcessDone) {
		t.Fatalf("signal after Wait: %v, want ErrProcessDone", err)
	}
}

func TestRedoxWaitKilledChild(t *testing.T) {
	p, _ := startRedoxWaitChild(t)
	if err := p.Kill(); err != nil {
		t.Fatal(err)
	}
	if ready, err := os.BlockUntilWaitable(p); !ready || err != nil {
		t.Fatalf("observation: ready=%v err=%v", ready, err)
	}
	state, err := p.Wait()
	if err != nil {
		t.Fatal(err)
	}
	status := state.Sys().(syscall.WaitStatus)
	if !status.Signaled() || status.Signal() != syscall.SIGKILL {
		t.Fatalf("Wait status = %v, want SIGKILL", status)
	}
}

func TestRedoxWaitSignal(t *testing.T) {
	for i := 0; i < 20; i++ {
		p, gate := startRedoxWaitChild(t)
		done := make(chan error, 1)
		go func() {
			_, err := p.Wait()
			done <- err
		}()
		if _, err := gate.Write([]byte{1}); err != nil {
			t.Fatal(err)
		}
		gate.Close()
	loop:
		for {
			select {
			case err := <-done:
				if err != nil {
					t.Fatal(err)
				}
				break loop
			default:
			}
			if err := p.Signal(syscall.Signal(0)); err != nil && !errors.Is(err, os.ErrProcessDone) {
				t.Fatal(err)
			}
			runtime.Gosched()
		}
		if err := p.Signal(syscall.Signal(0)); !errors.Is(err, os.ErrProcessDone) {
			t.Fatalf("iteration %d: signal after Wait = %v", i, err)
		}
	}
}
