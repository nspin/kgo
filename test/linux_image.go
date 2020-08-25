package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "github.com/nspin/kgo/pkg"
)

func main() {

    linuxFile := "../linux/arch/arm64/boot/Image"

    linux, err := ioutil.ReadFile(linuxFile)
    if err != nil {
        log.Fatal(err)
    }

    header, err := kgo.ReadLinuxImageHeader(linux)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", header)
}
