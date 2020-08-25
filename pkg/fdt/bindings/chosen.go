package bindings

import (
	"fmt"
	"github.com/nspin/kgo/pkg/fdt"
)

type ChosenLinux struct {
	Bootargs string
	// Optional
	StdoutPathOk bool
	StdoutPath   string
	InitrdOk     bool
	InitrdStart  uint64
	InitrdEnd    uint64
}

func (c *ChosenLinux) SetStdoutPath(stdoutPath string) {
	c.StdoutPathOk = true
	c.StdoutPath = stdoutPath
}

func (c *ChosenLinux) UseStdoutPath(dt *fdt.DeviceTree) (ok bool) {
	chosen, ok := dt.Root.Children["chosen"]
	if ok {
		stdoutPathRaw, ok := chosen.Properties["stdout-path"]
		stdoutPath, err := fdt.ValAsString(stdoutPathRaw)
		if err != nil {
			panic("invalid /chosen/stdout-path")
		}
		if ok {
			c.SetStdoutPath(stdoutPath)
		}
	}
	return
}

func (c *ChosenLinux) SetInitrd(start, end uint64) {
	c.InitrdOk = true
	c.InitrdStart = start
	c.InitrdEnd = end
}

func (c *ChosenLinux) Node() (chosen *fdt.Node) {
	chosen = fdt.NewNode()
	chosen.Properties["bootargs"] = fdt.ValFromString(c.Bootargs)
	if c.StdoutPathOk {
		chosen.Properties["stdout-path"] = fdt.ValFromString(c.StdoutPath)
	}
	if c.InitrdOk {
		chosen.Properties["linux,initrd-start"] = fdt.ValBigEndian(uint32(c.InitrdStart))
		chosen.Properties["linux,initrd-end"] = fdt.ValBigEndian(uint32(c.InitrdEnd))
	}
	return
}

func (c *ChosenLinux) Apply(dt *fdt.DeviceTree) {
	dt.Root.Children["chosen"] = c.Node()
}

type ChosenXen struct {
	XenBootargs  string
	Dom0Bootargs string
	KernelStart  uint64
	KernelSize   uint64
	// Optional
	StdoutPathOk bool
	StdoutPath   string
	RamdiskOk    bool
	RamdiskStart uint64
	RamdiskSize  uint64
}

func (c *ChosenXen) SetStdoutPath(stdoutPath string) {
	c.StdoutPathOk = true
	c.StdoutPath = stdoutPath
}

func (c *ChosenXen) UseStdoutPath(dt *fdt.DeviceTree) (ok bool) {
	chosen, ok := dt.Root.Children["chosen"]
	if ok {
		stdoutPathRaw, ok := chosen.Properties["stdout-path"]
		stdoutPath, err := fdt.ValAsString(stdoutPathRaw)
		if err != nil {
			panic("invalid /chosen/stdout-path")
		}
		if ok {
			c.SetStdoutPath(stdoutPath)
		}
	}
	return
}

func (c *ChosenXen) SetRamdisk(start, size uint64) {
	c.RamdiskOk = true
	c.RamdiskStart = start
	c.RamdiskSize = size
}

func (c *ChosenXen) Node() (chosen *fdt.Node) {
	chosen = fdt.NewNode()
	chosen.Properties["xen,xen-bootargs"] = fdt.ValFromString(c.XenBootargs)
	chosen.Properties["xen,dom0-bootargs"] = fdt.ValFromString(c.Dom0Bootargs)
	if c.StdoutPathOk {
		chosen.Properties["stdout-path"] = fdt.ValFromString(c.StdoutPath)
	}

	chosen.Children["modules"] = fdt.NewNode()
	chosen.Children["modules"].Properties["#address-cells"] = fdt.ValBigEndian(uint32(2))
	chosen.Children["modules"].Properties["#size-cells"] = fdt.ValBigEndian(uint32(2))

	kernel := fdt.NewNode()
	kernel.Properties["compatible"] = fdt.ValFromStrings([]string{"multiboot,kernel", "multiboot,module"})
	kernel.Properties["reg"] = fdt.ValFromVals([][]byte{
		fdt.ValBigEndian(c.KernelStart),
		fdt.ValBigEndian(c.KernelSize),
	})
	chosen.Children["modules"].Children[fmt.Sprintf("module@%x", c.KernelStart)] = kernel

	if c.RamdiskOk {
		ramdisk := fdt.NewNode()
		ramdisk.Properties["compatible"] = fdt.ValFromStrings([]string{"multiboot,ramdisk", "multiboot,module"})
		ramdisk.Properties["reg"] = fdt.ValFromVals([][]byte{
			fdt.ValBigEndian(c.RamdiskStart),
			fdt.ValBigEndian(c.RamdiskSize),
		})
		chosen.Children["modules"].Children[fmt.Sprintf("module@%x", c.RamdiskStart)] = ramdisk
	}
	return
}

func (c *ChosenXen) Apply(dt *fdt.DeviceTree) {
	dt.Root.Children["chosen"] = c.Node()
}
