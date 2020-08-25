package fdt

const (
    Magic = 0xd00dfeed
    HeaderSize = 10 * 4

    version = 17
    lastCompVersion = 16

    tokBeginNode = uint32(1)
    tokEndNode = uint32(2)
    tokProp = uint32(3)
    tokNop = uint32(4)
    tokEnd = uint32(9)
)

type Header struct {
    Magic uint32
    TotalSize uint32
    OffDtStruct uint32
    OffDtStrings uint32
    OffMemRsvmap uint32
    Version uint32
    LastCompVersion uint32
    BootCpuidPhys uint32
    SizeDtStrings uint32
    SizeDtStruct uint32
}

type DeviceTree struct {
    MemReserveMap []*ReserveEntry
    Root *Node
    BootCpuidPhys uint32
}

type ReserveEntry struct {
    Address uint64
    Size uint64
}

type Value []byte

type Node struct {
    Properties map[string]Value
    Children map[string]*Node
}

type propHeader struct {
    Len uint32
    NameOff uint32
}

func NewNode() *Node {
    return &Node{
        Properties: make(map[string]Value),
        Children: make(map[string]*Node),
    }
}
