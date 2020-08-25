package main

import (
    "github.com/nspin/kgo/pkg/fdt"
    "log"
    "io/ioutil"
)

func main() {
    fdt0, err := ioutil.ReadFile("./test/rk3399-gru-kevin-r5.dtb")
    if err != nil {
        log.Fatal(err)
    }
    dt0, err := fdt.Read(fdt0)
    if err != nil {
        log.Fatal(err)
    }
    dt0.DebugShow()
    // fdt1, err := fdt.Write(dt0)
    // if err != nil {
    //     log.Fatal(err)
    // }
    // dt1, err := fdt.Read(fdt1)
    // if err != nil {
    //     log.Fatal(err)
    // }
    // fdt_debug.DebugShowDeviceTree(dt1)
}
