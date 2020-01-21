package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"

	"encoding/hex"
	"fmt"
	"hash/crc32"
	"os"

	"reflect"
	"strconv"
)

var Hasher = sha1.New()

func (Footer *Footer) ReadFooter(file *os.File, name string) {
	for i := 0; i < len(Footer.Siglen); i++ {
		err := binary.Read(file, binary.BigEndian, &Footer.Siglen[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("signature[%d] length = %d\n", i, Footer.Siglen[i])
		Footer.Sigret[i] = make([]byte, Footer.Siglen[i])
		f, err := os.Create(fmt.Sprintf("%s/%s_sig_%d", PathUnpacked, name, i))
		defer f.Close()
		err = binary.Read(file, binary.BigEndian, &Footer.Sigret[i])
		if err != nil {
			fmt.Println(err)
			return
		}
		err = binary.Write(f, binary.BigEndian, Footer.Sigret[i])
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err := binary.Read(file, binary.BigEndian, &Footer.Padding)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = binary.Read(file, binary.BigEndian, &Footer.Sha1)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func main() {
	filename := firmware
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	os.Mkdir(PathUnpacked, os.ModePerm)

	header := Rotek{}
	err = binary.Read(file, binary.BigEndian, &header.Header)
	if err != nil {
		fmt.Println(err)
		return
	}
	n := bytes.IndexByte(header.Header.VendorName[:], 0)
	fmt.Println("vendor:", string(header.Header.VendorName[:n]))
	n = bytes.IndexByte(header.Header.DeviceName[:], 0)
	fmt.Println("device", string(header.Header.DeviceName[:n]), "hw", header.Header.MaxSupportedHwRevision)
	fmt.Println("version ", header.Header.VersionMajor, header.Header.VersionMinor, header.Header.VersionBuild)
	header.Footer.ReadFooter(file, firmware)

	for i := uint16(0); i < header.Header.NumofBlocks; i++ {
		var block Block
		err := binary.Read(file, binary.BigEndian, &block.Header)

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("type = ", block.Header.Type, Type(block.Header.Type).String())

		block.Footer.ReadFooter(file, Type(block.Header.Type).String())

		PreLen := block.Header.BlockSize - uint32(block.Header.BlockHeaderSize)

		Len := PreLen / BufferSize

		f, err := os.Create(fmt.Sprintf("%s/%s", PathUnpacked, Type(block.Header.Type).Name()))

		defer f.Close()

		hasher := sha1.New()
		tablePolynomial := crc32.MakeTable(0xEDB88320)
		hash := crc32.New(tablePolynomial)

		for j := uint32(0); j < Len; j++ {
			var Buffer [BufferSize]byte
			err = binary.Read(file, binary.BigEndian, &Buffer)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = binary.Write(f, binary.BigEndian, Buffer)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = hasher.Write(Buffer[:])
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err := hash.Write(Buffer[:])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		sha := hasher.Sum(nil)[:]
		shaa := block.Footer.Sha1[:]

		hashInBytes := hash.Sum(nil)[:]
		s := hex.EncodeToString(hashInBytes)
		n, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			panic(err)
		}
		crc := uint32(n)

		if block.Header.CRC == crc {
			fmt.Println("right crc")
		} else {
			fmt.Println("error crc!!!!")
		}
		if reflect.DeepEqual(sha, shaa) {
			fmt.Println("right sha1")
		} else {
			fmt.Println("error sha1!!!!")
		}
	}
}
