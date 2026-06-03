// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package testing

import (
	"fmt"
	"os"
)

func captureStdout(stdout *os.File) func() string {
	f, err := os.CreateTemp("", "go-example-stdout-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Stdout = f
	return func() string {
		name := f.Name()
		f.Close()
		os.Stdout = stdout
		out, err := os.ReadFile(name)
		os.Remove(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: reading stdout capture: %v\n", err)
			os.Exit(1)
		}
		return string(out)
	}
}
