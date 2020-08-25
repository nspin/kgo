package kgo

import (
    "fmt"
    "io/ioutil"

    "github.com/nspin/kgo/pkg/fdt"
    "github.com/nspin/kgo/pkg/fdt/bindings"
)

const (
    SysDtbPath = "/sys/firmware/fdt"
)

type Params struct {
    LinuxPath, LinuxParams string
    DtbPath string
    // Optional
    InitrdOk bool
    InitrdPath string
    XenOk bool
    XenPath, XenParams string
    StdoutPathOk bool
    StdoutPath string
}

func (p *Params) Run() (err error) {

    m, err := AvailableMemory()
    if err != nil {
        return
    }

    s := NewSegments(m)

    var ok bool
    var xenRange *Range
    if p.XenOk {
        var (
            xen []byte
            xenHeader *LinuxImageHeader
        )
        xen, err = ioutil.ReadFile(p.XenPath)
        if err != nil {
            return
        }
        xenHeader, err = ReadLinuxImageHeader(xen)
        if err != nil {
            return
        }
        xenRange, ok = s.AllocateAligned(0x200000, xenHeader.TextOffset, xenHeader.ImageSize, xen) // xenHeader.ImageSize)
        if !ok {
            fmt.Errorf("cannot allocate for xen")
        }
    }

    linux, err := ioutil.ReadFile(p.LinuxPath)
    if err != nil {
        return
    }
    linuxHeader, err := ReadLinuxImageHeader(linux)
    if err != nil {
        return
    }
    linuxRange, ok := s.AllocateAligned(0x200000, linuxHeader.TextOffset, linuxHeader.ImageSize, linux) // linuxHeader.ImageSize
    if !ok {
        fmt.Errorf("cannot allocate for linux")
    }

    initrd, err := ioutil.ReadFile(p.InitrdPath)
    if err != nil {
        return
    }
    initrdRange, ok := s.Allocate(initrd)
    if !ok {
        fmt.Errorf("cannot allocate for initrd")
    }

    dtb0, err := ioutil.ReadFile(p.DtbPath)
    if err != nil {
        return
    }
    dt, err := fdt.Read(dtb0)
    if err != nil {
        return
    }

    if p.XenOk {
        c := &bindings.ChosenXen{
            XenBootargs: p.XenParams,
            Dom0Bootargs: p.LinuxParams,
            KernelStart: linuxRange.Start,
            KernelSize: linuxRange.Size(),
        }
        c.SetRamdisk(initrdRange.Start, initrdRange.Size())
        if p.StdoutPathOk {
            c.SetStdoutPath(p.StdoutPath)
        } else {
            c.UseStdoutPath(dt)
        }
        c.Apply(dt)
    } else {
        c := &bindings.ChosenLinux{
            Bootargs: p.LinuxParams,
        }
        c.SetInitrd(initrdRange.Start, initrdRange.End)
        if p.StdoutPathOk {
            c.SetStdoutPath(p.StdoutPath)
        } else {
            c.UseStdoutPath(dt)
        }
        c.Apply(dt)
    }

    dtb, err := fdt.Write(dt)
    if err != nil {
        return
    }
    dtbRange, ok := s.Allocate(dtb)
    if !ok {
        fmt.Errorf("cannot allocate for dtb")
    }

    var purgatory []byte
    if p.XenOk {
        purgatory = Purgatory(dtbRange.Start, xenRange.Start)
    } else {
        purgatory = Purgatory(dtbRange.Start, linuxRange.Start)
    }
    purgatoryRange, ok := s.Allocate(purgatory)
    if !ok {
        fmt.Errorf("cannot allocate for purgatory")
    }

    entry := purgatoryRange.Start

    fmt.Println("chosen:")
    dt.Root.Children["chosen"].DebugShow()
    fmt.Println()
    fmt.Printf("entry: 0x%x\n", entry)
    fmt.Println()
    s.DebugShow(0)
    fmt.Println()

    err = Load(entry, s, 0)
    if err != nil {
        return
    }

    err = Reboot()
    if err != nil {
        return
    }

    // UNREACHABLE
    // avoid gc (still prone to relocation)
    fmt.Println(s)
    return
}

