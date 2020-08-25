package main

import (
    "log"
    "fmt"
    "github.com/nspin/kgo/pkg"
)

func main() {
    iomem, err := kgo.ParseIOMem()
    if err != nil {
        log.Fatal(err)
    }
    iomem.DebugShow()

    m, err := kgo.AvailableMemory()
    if err != nil {
        log.Fatal(err)
    }

    for _, r := range m.Ranges {
        fmt.Printf("%x-%x\n", r.Start, r.End)
    }
}
