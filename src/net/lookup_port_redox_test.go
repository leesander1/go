// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox && cgo && !netgo

package net

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestRedoxCgoLookupPortValues(t *testing.T) {
	ctx := context.Background()
	for _, tt := range []struct {
		network string
		service string
		port    int
	}{
		{"tcp", "smtp", 25},
		{"tcp", "http", 80},
		{"udp", "domain", 53},
	} {
		port, err := cgoLookupPort(ctx, tt.network, tt.service)
		if err != nil || port != tt.port {
			t.Errorf("cgoLookupPort(%q, %q) = %d, %v; want %d, nil", tt.network, tt.service, port, err, tt.port)
		}
	}

	_, err := cgoLookupPort(ctx, "udp", "_-unknown-service-")
	var dnsErr *DNSError
	if !errors.As(err, &dnsErr) || !dnsErr.IsNotFound {
		t.Fatalf("unknown service error = %v; want not-found DNSError", err)
	}
}

func TestRedoxCgoLookupPortConcurrent(t *testing.T) {
	ctx := context.Background()
	const workers = 16
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			network, service, want := "tcp", "http", 80
			if i%2 != 0 {
				network, service, want = "udp", "domain", 53
			}
			port, err := cgoLookupPort(ctx, network, service)
			if err != nil {
				errs <- err
			} else if port != want {
				errs <- &AddrError{Err: "unexpected port", Addr: service}
			}
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}
