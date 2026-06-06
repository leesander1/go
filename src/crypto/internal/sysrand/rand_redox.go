// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sysrand

import (
	"io"
	"os"
	"sync"
)

var redoxRand struct {
	once sync.Once
	file *os.File
	err  error
}

func read(b []byte) error {
	redoxRand.once.Do(func() {
		redoxRand.file, redoxRand.err = os.Open("/scheme/rand")
	})
	if redoxRand.err != nil {
		return redoxRand.err
	}
	_, err := io.ReadFull(redoxRand.file, b)
	return err
}
