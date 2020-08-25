package fdt

import (
    "bytes"
    "encoding/binary"
    "fmt"
)

func ValAsString(b []byte) (v string, err error) {
    if len(b) == 0 || b[len(b)-1] != 0 {
        err = fmt.Errorf("invalid null-terminated string: %s", b)
        return
    }
    v = string(b[:len(b)-1])
    return
}

func ValFromString(v string) []byte {
    return append([]byte(v), 0)
}

func ValFromStrings(v []string) []byte {
    b := []byte{}
    for _, s := range v {
        b = append(b, ValFromString(s)...)
    }
    return b
}

func ValFromVals(v [][]byte) []byte {
    b := []byte{}
    for _, x := range v {
        b = append(b, x...)
    }
    return b
}

func ValBigEndian(v interface{}) []byte {
    buf := &bytes.Buffer{}
    err := binary.Write(buf, binary.BigEndian, v)
    if err != nil {
        panic("how did this happen")
    }
    return buf.Bytes()
}
