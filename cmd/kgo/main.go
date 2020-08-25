package main

import (
    "log"
    "flag"
    "github.com/nspin/kgo/pkg"
)

func main() {
    p := &kgo.Params{}

    flag.StringVar(&p.LinuxPath, "linux", "", "path of linux (in 'Image' format)")
    flag.StringVar(&p.LinuxParams, "linux-params", "", "linux command line (optional)")
    flag.StringVar(&p.InitrdPath, "initrd", "", "path of initrd (optional)")
    flag.StringVar(&p.XenPath, "xen", "", "path of xen (in 'Image' format) (optional)")
    flag.StringVar(&p.XenParams, "xen-params", "", "xen command line (optional)")
    flag.StringVar(&p.DtbPath, "dtb", kgo.SysDtbPath, "path of dtb (default: /sys/firmware/fdt)")
    flag.StringVar(&p.StdoutPath, "stdout-path", "", "value for /chosen/stdout-path device tree property (optional)")

    flag.Parse()

    if p.LinuxPath == "" {
        log.Fatal("missing required flag: -linux")
    }
    if p.InitrdPath != "" {
        p.InitrdOk = true
    }
    if p.XenPath != "" {
        p.XenOk = true
    }
    if p.StdoutPath != "" {
        p.StdoutPathOk = true
    }

    err := p.Run()
    if err != nil {
        log.Fatal(err)
    }
}
