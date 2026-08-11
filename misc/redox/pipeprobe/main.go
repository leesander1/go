// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"syscall"
	"time"
)

func main() {
	fmt.Printf("REDOXPIPE_PROBE_BEGIN gomaxprocs=%d\n", runtime.GOMAXPROCS(0))
	if err := rawSyscallPipeBlockedEOF(); err != nil {
		fmt.Printf("REDOXPIPE_PROBE_FAIL err=%v\n", err)
		os.Exit(1)
	}
	if err := osPipeBlockedEOF(); err != nil {
		fmt.Printf("REDOXPIPE_PROBE_FAIL err=%v\n", err)
		os.Exit(1)
	}
	fmt.Println("REDOXPIPE_PROBE_SUCCESS")
}

func rawSyscallPipeBlockedEOF() error {
	fmt.Println("REDOXPIPE_CASE_BEGIN name=raw-syscall-pipe-blocked-eof")
	var fds [2]int
	if err := syscall.Pipe2(fds[:], syscall.O_CLOEXEC); err != nil {
		return fmt.Errorf("pipe2: %w", err)
	}
	readFD, writeFD := fds[0], fds[1]
	defer syscall.Close(readFD)

	done := make(chan error, 1)
	go func() {
		var buf [1]byte
		fmt.Printf("REDOXPIPE_RAW_READ_BEGIN fd=%d\n", readFD)
		n, err := syscall.Read(readFD, buf[:])
		fmt.Printf("REDOXPIPE_RAW_READ_DONE n=%d err=%v\n", n, err)
		if n != 0 || err != nil {
			done <- fmt.Errorf("read got n=%d err=%v, want EOF", n, err)
			return
		}
		done <- nil
	}()

	time.Sleep(50 * time.Millisecond)
	fmt.Printf("REDOXPIPE_RAW_CLOSE_WRITE_BEGIN fd=%d\n", writeFD)
	if err := syscall.Close(writeFD); err != nil {
		return fmt.Errorf("close write fd: %w", err)
	}
	fmt.Println("REDOXPIPE_RAW_CLOSE_WRITE_DONE")

	select {
	case err := <-done:
		if err != nil {
			return err
		}
	case <-time.After(3 * time.Second):
		return fmt.Errorf("raw syscall pipe read did not observe EOF")
	}
	fmt.Println("REDOXPIPE_CASE_PASS name=raw-syscall-pipe-blocked-eof")
	return nil
}

func osPipeBlockedEOF() error {
	fmt.Println("REDOXPIPE_CASE_BEGIN name=os-pipe-blocked-eof")
	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("os pipe: %w", err)
	}
	defer r.Close()

	done := make(chan error, 1)
	go func() {
		var buf [1]byte
		fmt.Printf("REDOXPIPE_OS_READ_BEGIN fd=%d\n", r.Fd())
		n, err := r.Read(buf[:])
		fmt.Printf("REDOXPIPE_OS_READ_DONE n=%d err=%v\n", n, err)
		if n != 0 || err != io.EOF {
			done <- fmt.Errorf("read got n=%d err=%v, want EOF", n, err)
			return
		}
		done <- nil
	}()

	time.Sleep(50 * time.Millisecond)
	fmt.Printf("REDOXPIPE_OS_CLOSE_WRITE_BEGIN fd=%d\n", w.Fd())
	if err := w.Close(); err != nil {
		return fmt.Errorf("close write file: %w", err)
	}
	fmt.Println("REDOXPIPE_OS_CLOSE_WRITE_DONE")

	select {
	case err := <-done:
		if err != nil {
			return err
		}
	case <-time.After(3 * time.Second):
		return fmt.Errorf("os pipe read did not observe EOF")
	}
	fmt.Println("REDOXPIPE_CASE_PASS name=os-pipe-blocked-eof")
	return nil
}
