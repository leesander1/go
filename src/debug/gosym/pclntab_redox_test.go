// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build redox

package gosym

import (
	"debug/elf"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestRedoxExternalPclnTable(t *testing.T) {
	f, err := elf.Open(os.Args[0])
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	syms, err := f.Symbols()
	if err != nil {
		t.Fatal(err)
	}
	byName := make(map[string]elf.Symbol)
	for _, sym := range syms {
		byName[sym.Name] = sym
	}
	for _, name := range []string{"runtime.pclntab", "runtime.epclntab", "runtime.text"} {
		if _, ok := byName[name]; !ok {
			t.Fatalf("missing ELF symbol %s", name)
		}
	}
	start, end := byName["runtime.pclntab"], byName["runtime.epclntab"]
	if start.Section != end.Section || int(start.Section) >= len(f.Sections) {
		t.Fatal("invalid line-table section bounds")
	}
	section := f.Sections[start.Section]
	data, err := section.Data()
	if err != nil {
		t.Fatal(err)
	}
	if start.Value < section.Addr || end.Value <= start.Value || end.Value-section.Addr > uint64(len(data)) {
		t.Fatal("line table lies outside its ELF section")
	}
	// External linking places the table inside a larger ELF section.
	pcln := data[start.Value-section.Addr : end.Value-section.Addr]
	table, err := NewTable(nil, NewLineTable(pcln, byName["runtime.text"].Value))
	if err != nil {
		t.Fatal(err)
	}
	pc := uint64(reflect.ValueOf(TestRedoxExternalPclnTable).Pointer())
	file, line, fn := table.PCToLine(pc)
	if fn == nil || fn.Name != "debug/gosym.TestRedoxExternalPclnTable" {
		t.Fatalf("PCToLine(%#x): function = %+v", pc, fn)
	}
	if !strings.HasSuffix(file, "/pclntab_redox_test.go") || line <= 0 {
		t.Fatalf("PCToLine(%#x): location = %s:%d", pc, file, line)
	}
	if table.LookupFunc(fn.Name) != fn {
		t.Fatal("LookupFunc disagrees with PCToLine")
	}
}
