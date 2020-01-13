package main

const BufferSize = 128

const usage = "Usage: rotek-upgrade-package-extract [firmware-file-name]"
const PathUnpacked = "unpacked"
const firmware = "firmware.bin"

type Footer struct {
	SignatureSize [3]uint16
	Signature     [3][]byte
	Unused        uint16
	Sha1          [20]byte
}

type Block struct {
	Size       uint32
	HeaderSize uint16
	Type       uint16
	CRC        uint32
	Footer     Footer
}

type RotekHeader struct {
	Vendor    [32]byte
	Device    [32]byte
	V1        uint16
	V2        uint16
	V3        uint16
	V4        uint16
	Unk       [12]byte
	HwRev     uint32
	FileCount byte
	Unk1      [3]byte
}

type Rotek struct {
	Header RotekHeader
	Footer Footer
	File   []Block
}

type Type int

const (
	UBoot = iota
	Bootloader
	Kernel
	Rootfs
	Unk4 //misc?
	Persistent
	BackupKernel
)

func (t Type) String() string {
	return [...]string{"u-boot Image", "bootloader Image", "Linux Kernel Image", "Root FS Image", "unk4", "Branding Image", "Backup Linux Kernel Image"}[t]
}
func (t Type) Name() string {
	return [...]string{"u-boot.bin", "bootloader.bin", "boot.img", "rootfs.img", "unk4", "persist.img", "backup_kernel.img"}[t]
}
