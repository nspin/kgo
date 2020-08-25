package main

import (
    "log"
    "io/ioutil"
    "github.com/nspin/kgo/pkg"
)

func main() {
    purgatory := kgo.Purgatory(0x1337, 0x1338)
    err := ioutil.WriteFile("purgatory.bin", purgatory, 0644)
    if err != nil {
        log.Fatal(err)
    }
}
