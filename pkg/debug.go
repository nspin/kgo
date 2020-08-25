package kgo

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

func (r *Range) DebugShow(level int) {
	indent(level)
	fmt.Printf("0x%x-0x%x (0x%x)\n", r.Start, r.End, r.Size())
}

func (m *Memory) DebugShow(level int) {
	for _, r := range m.Ranges {
		r.DebugShow(level)
	}
}

func (s *Segments) DebugShow(level int) {
	indent(level)
	fmt.Printf("Memory:\n")
	s.Memory.DebugShow(level + 1)
	indent(level)
	fmt.Printf("Segments:\n")
	for _, seg := range s.Segments {
		indent(level + 1)
		fmt.Printf("Mem:\n")
		seg.Mem.DebugShow(level + 2)
		indent(level + 1)
		fmt.Printf("Buf:\n")
		DebugDump(seg.Buf, level+2, 64)
	}
}

func DebugDump(b []byte, level int, max int) {
	if 0 < max && max < len(b) {
		b = b[:max]
	}
	re := regexp.MustCompile(`(?m)^[0-9a-f]+  `)
	dump := hex.Dump(b)
	fmt.Print(re.ReplaceAllString(dump, strings.Repeat("  ", level)))
}

func indent(level int) {
	fmt.Print(strings.Repeat("  ", level))
}
