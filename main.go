package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const PathUnpacked = "unpacked"
const firmware = "firmware.bin"

func (Footer *Footer) ReadFooter(file *os.File, name string) {
	for i := 0; i < len(Footer.SignatureSize); i++ {
		err := binary.Read(file, binary.BigEndian, &Footer.SignatureSize[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("signature[%d] length = %d\n", i, Footer.SignatureSize[i])
		Footer.Signature[i] = make([]byte, Footer.SignatureSize[i])
		f, err := os.Create(fmt.Sprintf("%s/%s_sig_%d", PathUnpacked, name, i))
		defer f.Close()
		err = binary.Read(file, binary.LittleEndian, &Footer.Signature[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = binary.Write(f, binary.LittleEndian, Footer.Signature[i])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err := binary.Read(file, binary.LittleEndian, &Footer.Unused)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = binary.Read(file, binary.LittleEndian, &Footer.Sha1)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	file, err := os.Open(firmware)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	os.Mkdir(PathUnpacked, os.ModePerm)

	header := Rotek{}
	err = binary.Read(file, binary.LittleEndian, &header.Header)
	if err != nil {
		fmt.Println(err)
		return
	}
	value := make([]byte, 4)
	binary.LittleEndian.PutUint32(value, header.Header.HwRev)
	n := bytes.IndexByte(header.Header.Vendor[:], 0)
	fmt.Println("vendor:", string(header.Header.Vendor[:n]))
	n = bytes.IndexByte(header.Header.Device[:], 0)
	fmt.Println("device", string(header.Header.Device[:n]), "hw", binary.BigEndian.Uint32(value))
	v1, v2, v3, v4 := make([]byte, 2), make([]byte, 2), make([]byte, 2), make([]byte, 2)
	binary.LittleEndian.PutUint16(v1, header.Header.V1)
	binary.LittleEndian.PutUint16(v2, header.Header.V2)
	binary.LittleEndian.PutUint16(v3, header.Header.V3)
	binary.LittleEndian.PutUint16(v4, header.Header.V4)
	fmt.Println("version ", binary.BigEndian.Uint16(v1), binary.BigEndian.Uint16(v2), binary.BigEndian.Uint16(v3), binary.BigEndian.Uint16(v4))
	header.Footer.ReadFooter(file, firmware)

	for i := byte(0); i < header.Header.FileCount; i++ {
		var Block Block
		err := binary.Read(file, binary.BigEndian, &Block.Size)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = binary.Read(file, binary.BigEndian, &Block.HeaderSize)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = binary.Read(file, binary.LittleEndian, &Block.Type)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("type = ", Block.Type, Type(Block.Type).String())

		err = binary.Read(file, binary.LittleEndian, &Block.CRC)
		if err != nil {
			fmt.Println(err)
			return
		}
		Block.Footer.ReadFooter(file, Type(Block.Type).Name())
		Len := (Block.Size - uint32(Block.HeaderSize)) / BufferSize
		f, err := os.Create(fmt.Sprintf("%s/%s", PathUnpacked, Type(Block.Type).Name()))
		defer f.Close()
		for j := uint32(0); j < Len; j++ {
			var Buffer [BufferSize]byte
			err = binary.Read(file, binary.LittleEndian, &Buffer)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = binary.Write(f, binary.LittleEndian, Buffer)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
