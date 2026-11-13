// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package net

import (
	"os"
	"strings"
)

const redoxNetcfg = "/scheme/netcfg/ifaces/eth0"

var redoxInterfaces = [...]Interface{
	{
		Index: 1,
		MTU:   1486,
		Name:  "loopback",
		Flags: FlagUp | FlagLoopback | FlagRunning,
	},
	{
		Index: 2,
		MTU:   1486,
		Name:  "eth0",
		Flags: FlagUp | FlagBroadcast | FlagRunning,
	},
}

// If the ifindex is zero, interfaceTable returns mappings of all
// network interfaces. Otherwise it returns a mapping of a specific
// interface.
func interfaceTable(ifindex int) ([]Interface, error) {
	if ifindex < 0 || ifindex > len(redoxInterfaces) {
		return nil, errNoSuchInterface
	}

	first, last := 0, len(redoxInterfaces)
	if ifindex != 0 {
		first, last = ifindex-1, ifindex
	}

	ift := make([]Interface, 0, last-first)
	for i := first; i < last; i++ {
		ifi := redoxInterfaces[i]
		if ifi.Name == "eth0" {
			if mac, err := os.ReadFile(redoxNetcfg + "/mac"); err == nil {
				if hw, err := ParseMAC(strings.TrimSpace(string(mac))); err == nil {
					ifi.HardwareAddr = hw
				}
			}
		}
		ift = append(ift, ifi)
	}
	return ift, nil
}

// If the ifi is nil, interfaceAddrTable returns addresses for all
// network interfaces. Otherwise it returns addresses for a specific
// interface.
func interfaceAddrTable(ifi *Interface) ([]Addr, error) {
	if ifi != nil && (ifi.Index < 1 || ifi.Index > len(redoxInterfaces)) {
		return nil, errNoSuchInterface
	}

	addrs := make([]Addr, 0, 2)
	if ifi == nil || ifi.Index == redoxInterfaces[0].Index {
		_, network, err := ParseCIDR("127.0.0.1/8")
		if err != nil {
			return nil, err
		}
		network.IP = ParseIP("127.0.0.1")
		addrs = append(addrs, network)
	}
	if ifi == nil || ifi.Index == redoxInterfaces[1].Index {
		data, err := os.ReadFile(redoxNetcfg + "/addr/list")
		if err != nil {
			return nil, err
		}
		cidr := strings.TrimSpace(string(data))
		if cidr != "" && cidr != "Not configured" {
			ip, network, err := ParseCIDR(cidr)
			if err != nil {
				return nil, err
			}
			network.IP = ip
			if !ip.IsUnspecified() {
				addrs = append(addrs, network)
			}
		}
	}
	return addrs, nil
}

// Redox does not currently expose joined multicast groups through netcfg.
func interfaceMulticastAddrTable(ifi *Interface) ([]Addr, error) {
	if ifi == nil || ifi.Index < 1 || ifi.Index > len(redoxInterfaces) {
		return nil, errNoSuchInterface
	}
	return nil, nil
}
