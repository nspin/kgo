package kgo

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	LinuxImageHeaderSize         = 64
	LinuxImageHeaderMagic uint32 = 0x644d5241
)

type rawLinuxImageHeader struct {
	Code0      uint32
	Code1      uint32
	TextOffset uint64
	ImageSize  uint64
	Flags      uint64
	Res2       uint64
	Res3       uint64
	Res4       uint64
	Magic      uint32
	Res5       uint32
}

type LinuxImageHeader struct {
	TextOffset uint64
	ImageSize  uint64
	Flags      uint64
}

func ReadLinuxImageHeader(b []byte) (*LinuxImageHeader, error) {
	r := bytes.NewReader(b)
	raw := rawLinuxImageHeader{}
	err := binary.Read(r, binary.LittleEndian, &raw)
	if err != nil {
		return nil, err
	}
	if raw.Magic != LinuxImageHeaderMagic {
		return nil, fmt.Errorf("invalid Image header magic")
	}
	if raw.Res2 != 0 || raw.Res3 != 0 || raw.Res4 != 0 {
		return nil, fmt.Errorf("invalid Image reserved field")
	}
	return &LinuxImageHeader{
		TextOffset: raw.TextOffset,
		ImageSize:  raw.ImageSize,
		Flags:      raw.Flags,
	}, nil
}
