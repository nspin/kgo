package kgo

import (
	"os"
	"sort"
)

var (
	PageSize = uint64(os.Getpagesize())
	PageMask = PageSize - 1
)

func AlignUp(addr, alignment uint64) (aligned uint64) {
	return addr + ((alignment - (addr % alignment)) % alignment)
}

type Range struct {
	Start uint64
	End   uint64
}

func (r *Range) Size() uint64 {
	return r.End - r.Start
}

type Memory struct {
	Ranges []*Range
}

func NewMemory() (m *Memory) {
	return &Memory{}
}

func (m *Memory) Normalize() {
	m.Ranges = joinAdjacentRanges(filterEmptyRanges(m.Ranges))
}

func filterEmptyRanges(rs0 []*Range) (rs1 []*Range) {
	for _, r := range rs0 {
		if r.Size() != 0 {
			rs1 = append(rs1, r)
		}
	}
	return
}

func joinAdjacentRanges(rs0 []*Range) (rs1 []*Range) {
	for i := 0; i < len(rs0); {
		r := rs0[i]
		for {
			i++
			if !(i < len(rs0) && r.End == rs0[i].Start) {
				rs1 = append(rs1, r)
				break
			}
		}
	}
	return
}

type Segment struct {
	Buf []byte
	Mem *Range
}

type Segments struct {
	Memory   *Memory
	Segments []*Segment
}

func NewSegments(m *Memory) (s *Segments) {
	return &Segments{
		Memory: m,
	}
}

func (s *Segments) Insert(seg *Segment) {
	// TODO improve
	s.Segments = append(s.Segments, seg)
	sort.Sort(s)
}

func (s *Segments) Len() int {
	return len(s.Segments)
}
func (s *Segments) Swap(i, j int) {
	s.Segments[i], s.Segments[j] = s.Segments[j], s.Segments[i]
}
func (s *Segments) Less(i, j int) bool {
	return s.Segments[i].Mem.Start < s.Segments[j].Mem.Start
}

func (m *Memory) Len() int {
	return len(m.Ranges)
}
func (m *Memory) Swap(i, j int) {
	m.Ranges[i], m.Ranges[j] = m.Ranges[j], m.Ranges[i]
}
func (m *Memory) Less(i, j int) bool {
	return m.Ranges[i].Start < m.Ranges[j].Start
}
