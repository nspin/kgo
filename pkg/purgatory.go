package kgo

import (
	"encoding/binary"
)

func ldr(imm, reg uint32) uint32 {
	//          [imm19              ][Rt ]
	// 01011000 xxxxxxxx xxxxxxxx xxxyyyyy
	return 0x58000000 | (imm << 5) | reg
}

func Purgatory(dtb_base, kernel_entry uint64) []byte {
	dtb_base_lo := uint32(dtb_base)
	dtb_base_hi := uint32(dtb_base >> 32)
	kernel_entry_lo := uint32(kernel_entry)
	kernel_entry_hi := uint32(kernel_entry >> 32)

	instrs := []uint32{
		ldr(6, 4),  //     ldr x4, kernel_entry
		ldr(7, 0),  //     ldr x0, dtb_base
		0xaa1f03e1, //     mov x1, xzr
		0xaa1f03e2, //     mov x2, xzr
		0xaa1f03e3, //     mov x3, xzr
		0xd61f0080, //     br x4
		//
		// kernel_entry:
		kernel_entry_lo, //     .quad ${kernel_entry}
		kernel_entry_hi, //
		// dtb_base:
		dtb_base_lo, //     .quad ${dtb_base}
		dtb_base_hi, //
	}

	code := make([]byte, len(instrs)*4)
	for i, instr := range instrs {
		binary.LittleEndian.PutUint32(code[i*4:], instr)
	}
	return code
}
