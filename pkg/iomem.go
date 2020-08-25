package kgo

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reIOMemLine *regexp.Regexp
)

func init() {
	reIOMemLine = regexp.MustCompile("(?P<indent>(?:  )*)(?P<start>[0-9a-f]+)-(?P<end>[0-9a-f]+) : (?P<type>.*)")
}

type IOMem []*IOMemNode

type IOMemNode struct {
	Range    *Range
	Type     string
	Children IOMem
}

// TODO validate (ensure ordered, non-overlapping, etc.)
func ParseIOMem() (iomem IOMem, err error) {
	f, err := os.Open("/proc/iomem")
	if err != nil {
		return
	}
	s := newUnScanner(f)
	return parseIOMem(s, 0)
}

func parseIOMem(s *unScanner, indent int) (iomem IOMem, err error) {
	var (
		line     string
		ok       bool
		start    uint64
		end      uint64
		children IOMem
	)

	for {
		line, ok, err = s.getLine()
		if err != nil {
			return
		}
		if !ok {
			return
		}

		match := reIOMemLine.FindStringSubmatch(line)
		if match == nil {
			err = fmt.Errorf("invalid iomem line: %s", line)
			return
		}
		matchIndent := match[1]
		matchStart := match[2]
		matchEnd := match[3]
		matchType := match[4]

		observedIndent := len(matchIndent) / 2
		if observedIndent > indent {
			err = fmt.Errorf("unexpected indent")
			return
		}
		if observedIndent < indent {
			s.unGetLine(line)
			return
		}

		start, err = strconv.ParseUint(matchStart, 16, 64)
		if err != nil {
			return
		}
		end, err = strconv.ParseUint(matchEnd, 16, 64)
		if err != nil {
			return
		}

		children, err = parseIOMem(s, indent+1)
		if err != nil {
			return
		}

		iomem = append(iomem, &IOMemNode{
			Range: &Range{
				Start: uint64(start),
				End:   uint64(end),
			},
			Type:     matchType,
			Children: children,
		})
	}
}

func (iomem *IOMem) DebugShow() {
	iomem.debugShow(0)
}

func (iomem IOMem) debugShow(indent int) {
	for _, node := range iomem {
		fmt.Printf("%s%x-%x : %s\n", strings.Repeat("  ", indent), node.Range.Start, node.Range.End, node.Type)
		node.Children.debugShow(indent + 1)
	}
}
