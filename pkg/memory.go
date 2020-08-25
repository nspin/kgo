package kgo

import (
	"sort"
)

// TODO validate? deal with out of order and overlapping? find spec for /proc/iomem
// TODO handle secmem node in dt, other sources of information
func AvailableMemory() (m *Memory, err error) {
	ranges, err := availableMemory()
	if err != nil {
		return
	}
	m = &Memory{
		Ranges: ranges,
	}
	m.Normalize()
	// TODO secmem from dt
	m.Ranges = m.Ranges[1:]
	return
}

func availableMemory() (ranges []*Range, err error) {
	iomem, err := ParseIOMem()
	if err != nil {
		return
	}
	for _, outer := range iomem {
		if outer.Type == "System RAM" {
			start := outer.Range.Start
			for _, inner := range outer.Children {
				if inner.Type == "reserved" {
					ranges = append(ranges, &Range{
						Start: start,
						End:   inner.Range.Start,
					})
					start = inner.Range.End
				}
			}
			ranges = append(ranges, &Range{
				Start: start,
				End:   outer.Range.End,
			})
		}
	}
	return
}

// Greedily allocate a page-aligned (on both ends) chunk of memory of at least 'size'
// The chunk will also be at 'offset' from 'alignment'
func (m *Memory) AllocateInternal(size uint64, extraAlignment bool, alignment, offset uint64) (r *Range, ok bool) {
	if extraAlignment {
		if (alignment|offset)&PageMask != 0 {
			panic("alignment and offset must be multiples of page size")
		}
	}
	for i, ri := range m.Ranges {
		r = &Range{}
		if extraAlignment {
			r.Start = AlignUp(ri.Start, alignment) + offset
		} else {
			r.Start = AlignUp(ri.Start, PageSize)
		}
		r.End = r.Start + AlignUp(size, PageSize)
		if r.End < ri.End {
			// TODO improve
			m.Ranges[i] = &Range{
				Start: ri.Start,
				End:   r.Start,
			}
			m.Ranges = append(m.Ranges, &Range{
				Start: r.End,
				End:   ri.End,
			})
			sort.Sort(m)
			m.Normalize()
			ok = true
			return
		}
	}
	return
}

func (s *Segments) Allocate(buf []byte) (r *Range, ok bool) {
	size := uint64(len(buf))
	rp, ok := s.Memory.AllocateInternal(size, false, 0, 0)
	if !ok {
		return
	}
	s.Insert(&Segment{
		Buf: buf,
		Mem: rp,
	})
	r = &Range{
		Start: rp.Start,
		End:   rp.Start + size,
	}
	return
}

func (s *Segments) AllocateAligned(alignment, offset, size uint64, buf []byte) (r *Range, ok bool) {
	rp, ok := s.Memory.AllocateInternal(size, true, alignment, offset)
	if !ok {
		return
	}
	s.Insert(&Segment{
		Buf: buf,
		Mem: rp,
	})
	r = &Range{
		Start: rp.Start,
		End:   rp.Start + size,
	}
	return
}
