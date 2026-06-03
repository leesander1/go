// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !js && !wasip1

// TODO(@musiol, @odeke-em): re-unify this entire file back into
// example.go when js/wasm gets an os.Pipe implementation
// and no longer needs this separation.

package testing

import (
	"fmt"
	"os"
	"time"
)

func runExample(eg InternalExample) (ok bool) {
	if chatty.on {
		fmt.Printf("%s=== RUN   %s\n", chatty.prefix(), eg.Name)
	}

	// Capture stdout.
	stdout := os.Stdout
	finishCapture := captureStdout(stdout)

	finished := false
	start := time.Now()

	// Clean up in a deferred call so we can recover if the example panics.
	defer func() {
		timeSpent := time.Since(start)

		// Close pipe, restore stdout, get output.
		out := finishCapture()

		err := recover()
		ok = eg.processRunResult(out, timeSpent, finished, err)
	}()

	// Run example.
	eg.F()
	finished = true
	return
}
