package fdt

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

func showNode(node *Node, level int) {
	for k, v := range node.Properties {
		indent(level)
		fmt.Printf("%v =\n", k)
		indentedDump(level+1, v)
	}
	for k, child := range node.Children {
		indent(level)
		fmt.Printf("%v {\n", k)
		showNode(child, level+1)
		indent(level)
		fmt.Printf("}\n")
	}
}

func indent(level int) {
	fmt.Printf(strings.Repeat("  ", level))
}

func indentedDump(level int, b []byte) {
	re := regexp.MustCompile(`(?m)^[0-9a-f]+  `)
	dump := hex.Dump(b)
	fmt.Print(re.ReplaceAllString(dump, strings.Repeat("  ", level)))
}

func (dt *DeviceTree) DebugShow() {
	for _, entry := range dt.MemReserveMap {
		fmt.Printf("/memreserve/ %x %x", entry.Address, entry.Size)
	}
	fmt.Printf("/ {\n")
	showNode(dt.Root, 1)
	fmt.Printf("}\n")
}

func (node *Node) DebugShow() {
	fmt.Printf("{\n")
	showNode(node, 1)
	fmt.Printf("}\n")
}
