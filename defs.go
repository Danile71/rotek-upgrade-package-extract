package main

const BufferSize = 32

const usage = "Usage: rotek-upgrade-package-extract [firmware-file-name]"
const PathUnpacked = "unpacked"
const firmware = "firmware.bin"

type Footer struct {
	Siglen  [3]uint16
	Sigret  [3][]byte
	Padding uint16
	Sha1    [20]byte
}

type BlockHeader struct {
	BlockSize       uint32
	BlockHeaderSize uint16
	Type            uint8
	Reserved        uint8
	CRC             uint32
}

type Block struct {
	Header BlockHeader
	Footer Footer
}

type RotekHeader struct {
	VendorName             [32]byte
	DeviceName             [32]byte
	VersionMajor           uint16
	VersionMinor           uint16
	VersionBuild           uint32
	BuildTime              uint64
	NumofBlocks            uint16
	ExtraSize              uint16
	MaxSupportedHwRevision uint32
	SignatureSize          uint16
	Padding                uint16
}

type Rotek struct {
	Header RotekHeader
	Footer Footer
	File2  []Block
}

type Type int

const (
	unused = iota
	UBoot
	Kernel
	Rootfs
	InstallScript
	Persistent
	BackupKernel
	PostDownloadScript
	Logo
)

func (t Type) String() string {
	return [...]string{"unused", "u-boot Image", "Linux Kernel Image", "Root FS Image", "Install Script", "Branding Image", "Backup Linux Kernel Image", "Post Download Script", "Logo"}[t]
}
func (t Type) Name() string {
	return [...]string{"unused", "u-boot.bin", "boot.img", "rootfs.img", "install_script", "persist.img", "backup_kernel.img", "post_download_script", "logo.bin"}[t]
}
