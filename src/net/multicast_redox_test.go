// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package net

import (
	"testing"
	"time"
)

func TestIPv4MulticastLoopbackDelivery(t *testing.T) {
	ifi := loopbackInterface()
	if ifi == nil {
		t.Fatal("loopback interface not found")
	}

	group := &UDPAddr{IP: IPv4(224, 0, 0, 254), Port: 12346}
	listener, err := ListenMulticastUDP("udp4", ifi, group)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	if err := listener.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}

	sender, err := ListenUDP("udp4", &UDPAddr{IP: IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer sender.Close()

	want := []byte("redox multicast")
	if n, err := sender.WriteToUDP(want, group); err != nil {
		t.Fatal(err)
	} else if n != len(want) {
		t.Fatalf("WriteToUDP wrote %d bytes; want %d", n, len(want))
	}

	got := make([]byte, len(want))
	if n, _, err := listener.ReadFromUDP(got); err != nil {
		t.Fatal(err)
	} else if n != len(want) || string(got[:n]) != string(want) {
		t.Fatalf("ReadFromUDP = %q, %d bytes; want %q", got[:n], n, want)
	}
}

func TestIPv4MulticastEthernetEgress(t *testing.T) {
	ifi, err := InterfaceByName("eth0")
	if err != nil {
		t.Fatal(err)
	}
	addrs, err := ifi.Addrs()
	if err != nil {
		t.Fatal(err)
	}

	var source IP
	for _, addr := range addrs {
		if network, ok := addr.(*IPNet); ok {
			if ip := network.IP.To4(); ip != nil && !ip.IsLoopback() {
				source = ip
				break
			}
		}
	}
	if source == nil {
		t.Fatal("eth0 has no IPv4 address")
	}

	sender, err := ListenUDP("udp4", &UDPAddr{IP: source})
	if err != nil {
		t.Fatal(err)
	}
	defer sender.Close()

	want := []byte("redox ethernet multicast")
	group := &UDPAddr{IP: IPv4(239, 255, 0, 1), Port: 12347}
	if n, err := sender.WriteToUDP(want, group); err != nil {
		t.Fatal(err)
	} else if n != len(want) {
		t.Fatalf("WriteToUDP wrote %d bytes; want %d", n, len(want))
	}
}

func TestIPv4MulticastLoopbackFanout(t *testing.T) {
	ifi := loopbackInterface()
	if ifi == nil {
		t.Fatal("loopback interface not found")
	}

	group := &UDPAddr{IP: IPv4(224, 0, 0, 254), Port: 12348}
	listeners := make([]*UDPConn, 2)
	for i := range listeners {
		listener, err := ListenMulticastUDP("udp4", ifi, group)
		if err != nil {
			t.Fatal(err)
		}
		defer listener.Close()
		if err := listener.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
			t.Fatal(err)
		}
		listeners[i] = listener
	}

	sender, err := ListenUDP("udp4", &UDPAddr{IP: IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer sender.Close()

	want := []byte("redox multicast fanout")
	if n, err := sender.WriteToUDP(want, group); err != nil {
		t.Fatal(err)
	} else if n != len(want) {
		t.Fatalf("WriteToUDP wrote %d bytes; want %d", n, len(want))
	}

	for i, listener := range listeners {
		got := make([]byte, len(want))
		if n, _, err := listener.ReadFromUDP(got); err != nil {
			t.Fatalf("listener %d: %v", i, err)
		} else if n != len(want) || string(got[:n]) != string(want) {
			t.Fatalf("listener %d: ReadFromUDP = %q, %d bytes; want %q", i, got[:n], n, want)
		}
	}
}

func TestIPv4MulticastLoopbackOption(t *testing.T) {
	ifi := loopbackInterface()
	if ifi == nil {
		t.Fatal("loopback interface not found")
	}

	group := &UDPAddr{IP: IPv4(224, 0, 0, 254), Port: 12349}
	listener, err := ListenMulticastUDP("udp4", ifi, group)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	sender, err := ListenUDP("udp4", &UDPAddr{IP: IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer sender.Close()

	if err := setIPv4MulticastLoopback(sender.fd, false); err != nil {
		t.Fatal(err)
	}
	suppressed := []byte("redox multicast suppressed")
	if _, err := sender.WriteToUDP(suppressed, group); err != nil {
		t.Fatal(err)
	}
	if err := listener.SetReadDeadline(time.Now().Add(250 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, len(suppressed))
	if _, _, err := listener.ReadFromUDP(buf); err == nil {
		t.Fatal("ReadFromUDP received multicast with loopback disabled")
	} else if timeout, ok := err.(Error); !ok || !timeout.Timeout() {
		t.Fatalf("ReadFromUDP with loopback disabled: %v", err)
	}

	if err := setIPv4MulticastLoopback(sender.fd, true); err != nil {
		t.Fatal(err)
	}
	want := []byte("redox multicast restored")
	if _, err := sender.WriteToUDP(want, group); err != nil {
		t.Fatal(err)
	}
	if err := listener.SetReadDeadline(time.Now().Add(2 * time.Second)); err != nil {
		t.Fatal(err)
	}
	got := make([]byte, len(want))
	if n, _, err := listener.ReadFromUDP(got); err != nil {
		t.Fatal(err)
	} else if n != len(want) || string(got[:n]) != string(want) {
		t.Fatalf("ReadFromUDP = %q, %d bytes; want %q", got[:n], n, want)
	}
}
