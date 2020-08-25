package fdt

import (
    "fmt"
    "io"
    "encoding/binary"
    "bytes"
)

func getString(fdt []byte, header *Header, offset uint32) ([]byte, error) {
    start := header.OffDtStrings + offset
    for end := start; end < header.OffDtStrings + header.SizeDtStrings && end < uint32(len(fdt)); end++ {
        if fdt[end] == 0 {
            return fdt[start : end], nil
        }
    }
    return nil, fmt.Errorf("invalid string:", offset)
}

func skipPadding(r io.Reader, off uint32) error {
    n := (4 - (off % 4)) % 4
    padding := make([]byte, n)
    _, err := io.ReadFull(r, padding)
    return err
}

func parseNode(fdt []byte, header *Header, r io.Reader) (node *Node, err error) {
    node = NewNode()
    var (
        tok uint32
        k []byte
        v []byte
        child *Node
    )
    for {
        err = binary.Read(r, binary.BigEndian, &tok)
        if err != nil {
            return
        }
        switch tok {
        case tokNop:
        case tokProp:
            ph := propHeader{}
            err = binary.Read(r, binary.BigEndian, &ph)
            if err != nil {
                return
            }
            k, err = getString(fdt, header, ph.NameOff)
            if err != nil {
                return
            }
            v = make([]byte, ph.Len)
            _, err = io.ReadFull(r, v)
            if err != nil {
                return
            }
            node.Properties[string(k)] = v
            skipPadding(r, ph.Len)
        case tokBeginNode:
            k = []byte{}
            var c byte
            var i uint32
            i = 0
            for {
                err = binary.Read(r, binary.BigEndian, &c)
                if err != nil {
                    return
                }
                i++
                if c == 0 {
                    break
                }
                k = append(k, c)
            }
            skipPadding(r, i)
            child, err = parseNode(fdt, header, r)
            if err != nil {
                return
            }
            node.Children[string(k)] = child
        case tokEndNode:
            return
        default:
            err = fmt.Errorf("unexpected token:", tok)
            return
        }
    }
    return nil, nil
}

func Read(fdt []byte) (dt *DeviceTree, err error) {
    dt = &DeviceTree{}

    header, err := ReadHeader(fdt)
    if err != nil {
        return
    }

    dt.BootCpuidPhys = header.BootCpuidPhys

    r := bytes.NewReader(fdt[header.OffMemRsvmap:])
    for {
        entry := &ReserveEntry{}
        err = binary.Read(r, binary.BigEndian, entry)
        if err != nil {
            return
        }
        if entry.Address == 0 && entry.Size == 0 {
            break
        }
        dt.MemReserveMap = append(dt.MemReserveMap, entry)
    }

    r = bytes.NewReader(fdt[header.OffDtStruct : header.OffDtStruct + header.SizeDtStruct])

    var tok uint32
    var emptyKey uint32
    err = binary.Read(r, binary.BigEndian, &tok)
    if err != nil {
        return
    }
    if tok != tokBeginNode {
        err = fmt.Errorf("unexpected token:", tok)
        return
    }
    err = binary.Read(r, binary.BigEndian, &emptyKey)
    if emptyKey != 0 {
        err = fmt.Errorf("unexpected root node key:", emptyKey)
        return
    }
    dt.Root, err = parseNode(fdt, header, r)
    if err != nil {
        return
    }
    err = binary.Read(r, binary.BigEndian, &tok)
    if err != nil {
        return
    }
    if tok != tokEnd {
        err = fmt.Errorf("unexpected token:", tok)
        return
    }

    return
}

func ReadHeader(fdt []byte) (header *Header, err error) {
    header = &Header{}
    r := bytes.NewReader(fdt)
    err = binary.Read(r, binary.BigEndian, header)
    if err != nil {
        err = fmt.Errorf("binary.Read failed:", err)
        return
    }
    if header.Magic != Magic {
        err = fmt.Errorf("invalid magic")
        return
    }
    if header.Version != version {
        err = fmt.Errorf("invalid version:", header.Version)
        return
    }
    // if header.LastCompVersion != lastCompVersion {
    //     err = fmt.Errorf("invalid last_comp_version:", header.LastCompVersion)
    //     return
    // }
    return
}
