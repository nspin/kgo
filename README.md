# kgo

`kexec` implemented in pure Go, with some extra flexibility.
Designed to boot a type-1 hypervisor such as Xen or Hafnium, with [LinuxBoot](https://www.linuxboot.org/) in mind.

### Usage

```
Usage of kgo
  -dtb string
        path of dtb (default "/sys/firmware/fdt")
  -initrd string
        path of initrd (optional)
  -linux string
        path of linux (in 'Image' format)
  -linux-params string
        linux command line (optional)
  -stdout-path string
        value for /chosen/stdout-path device tree property (optional)
  -xen string
        path of xen (in 'Image' format) (optional)
  -xen-params string
        xen command line (optional)
```
