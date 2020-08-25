package fdt

import (
    "bytes"
    "encoding/binary"
)

func Write(dt *DeviceTree) (fdt []byte, err error) {

    memrsv, err := writeMemrsv(dt)
    if err != nil {
        return
    }

    tree, strings, err := writeBody(dt)
    if err != nil {
        return
    }

    header := Header{
        Magic: Magic,
        Version: version,
        LastCompVersion: lastCompVersion,
        BootCpuidPhys: dt.BootCpuidPhys,
    }

    buf := &bytes.Buffer{}
    pad(buf, HeaderSize)

    align(buf, 8)
    header.OffMemRsvmap = uint32(buf.Len())
    buf.Write(memrsv)

    // pad(buf, 1024)

    align(buf, 4)
    header.OffDtStruct = uint32(buf.Len())
    header.SizeDtStruct = uint32(len(tree))
    buf.Write(tree)

    // pad(buf, 1024)

    header.OffDtStrings = uint32(buf.Len())
    header.SizeDtStrings = uint32(len(strings))
    buf.Write(strings)

    // pad(buf, 1024)

    header.TotalSize = uint32(buf.Len())

    fdt = buf.Bytes()

    headerBytes, err := WriteHeader(&header)
    if err != nil {
        return
    }
    fdt = append(headerBytes, fdt[len(headerBytes):]...)
    return
}

func WriteHeader(header *Header) (b []byte, err error) {
    buf := &bytes.Buffer{}
    err = binary.Write(buf, binary.BigEndian, header)
    if err != nil {
        return
    }
    b = buf.Bytes()
    return
}

func writeMemrsv(dt *DeviceTree) (b []byte, err error) {
    buf := &bytes.Buffer{}
    for entry := range dt.MemReserveMap {
        err = binary.Write(buf, binary.BigEndian, entry)
        if err != nil {
            return
        }
    }
    err = binary.Write(buf, binary.BigEndian, ReserveEntry{0, 0})
    if err != nil {
        return
    }
    b = buf.Bytes()
    return
}

func writeBody(dt *DeviceTree) (tree []byte, strings []byte, err error) {
    treeBuf := &bytes.Buffer{}
    stringsBuf := &bytes.Buffer{}
    err = binary.Write(treeBuf, binary.BigEndian, tokBeginNode)
    if err != nil {
        return
    }
    err = binary.Write(treeBuf, binary.BigEndian, uint32(0))
    if err != nil {
        return
    }
    err = writeNode(treeBuf, stringsBuf, dt.Root)
    if err != nil {
        return
    }
    err = binary.Write(treeBuf, binary.BigEndian, tokEndNode)
    if err != nil {
        return
    }
    err = binary.Write(treeBuf, binary.BigEndian, tokEnd)
    if err != nil {
        return
    }
    tree = treeBuf.Bytes()
    strings = stringsBuf.Bytes()
    return
}

func pad(buf *bytes.Buffer, n int) {
    padding := make([]byte, n)
    buf.Write(padding)
}

func align(buf *bytes.Buffer, n int) {
    pad(buf, (n - (buf.Len() % n)) % n)
}

func writeNode(tree *bytes.Buffer, strings *bytes.Buffer, node *Node) (err error) {
    for k, v := range node.Properties {
        err = binary.Write(tree, binary.BigEndian, tokProp)
        if err != nil {
            return
        }
        err = binary.Write(tree, binary.BigEndian, propHeader {
            Len: uint32(len(v)),
            NameOff: uint32(strings.Len()),
        })
        if err != nil {
            return
        }
        tree.Write(v)
        align(tree, 4)
        strings.Write([]byte(k))
        strings.WriteByte(0)
    }
    for k, child := range node.Children {
        err = binary.Write(tree, binary.BigEndian, tokBeginNode)
        if err != nil {
            return
        }
        tree.Write([]byte(k))
        tree.WriteByte(0)
        align(tree, 4)
        err = writeNode(tree, strings, child)
        if err != nil {
            return
        }
        err = binary.Write(tree, binary.BigEndian, tokEndNode)
        if err != nil {
            return
        }
    }
    return nil
}
