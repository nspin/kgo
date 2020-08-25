package kgo

import (
	"fmt"
	"syscall"
	"unsafe"
    "reflect"

	"golang.org/x/sys/unix"
)

func Reboot() error {
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_KEXEC)
    if err != nil {
		return fmt.Errorf("sys_reboot(..., kexec) = %v", err)
	}
	return nil
}

type RawSegment struct {
    Buf uintptr
    Bufsz uint
    Mem uintptr
    Memsz uint
}

func (s *Segments) toRaw() []RawSegment {
    raw := make([]RawSegment, len(s.Segments))
    for i, seg := range s.Segments {
        raw[i] = RawSegment{
			Buf: uintptr((unsafe.Pointer(&seg.Buf[0]))),
			Bufsz: uint(len(seg.Buf)),
            Mem: uintptr(seg.Mem.Start),
            Memsz: uint(seg.Mem.End - seg.Mem.Start),
        }
    }
    return raw
}

// TODO better hope the reaper doesn't reap or move stuff around
func Load(entry uint64, segments *Segments, flags uint64) error {
    raw := segments.toRaw()
	_, _, errno := unix.Syscall6(
		unix.SYS_KEXEC_LOAD,
		uintptr(entry),
		uintptr(len(raw)),
		uintptr(unsafe.Pointer(&raw[0])),
		uintptr(flags),
		0,
        0,
    )
    if errno != 0 {
        return fmt.Errorf("SYS_KEXEC_LOAD errno = %d", errno)
	}
    raw1 := segments.toRaw()
    if !reflect.DeepEqual(raw, raw1) {
        panic("raw segments changed")
    }
    return nil
}
