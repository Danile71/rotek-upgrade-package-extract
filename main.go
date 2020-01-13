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
	fmt.Println("Unused", Footer.Unused)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}

	file, err := os.Open(os.Args[1])
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

	fmt.Println("Unk:", header.Header.Unk)
	fmt.Println("Unk1:", header.Header.Unk1)
	header.File = make([]Block, header.Header.FileCount)

	for i := byte(0); i < header.Header.FileCount; i++ {
		err := binary.Read(file, binary.BigEndian, &header.File[i].Size)

		if err != nil {
			fmt.Println(err)
			return
		}

		err = binary.Read(file, binary.BigEndian, &header.File[i].HeaderSize)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = binary.Read(file, binary.LittleEndian, &header.File[i].Type)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = binary.Read(file, binary.BigEndian, &header.File[i].CRC)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("type = ", header.File[i].Type, Type(header.File[i].Type).String())

		header.File[i].Footer.ReadFooter(file, Type(header.File[i].Type).Name())

		Len := (header.File[i].Size - uint32(header.File[i].HeaderSize)) / BufferSize

		f, err := os.Create(fmt.Sprintf("%s/%s", PathUnpacked, Type(header.File[i].Type).Name()))

		defer f.Close()

		hasher := sha1.New()
		tablePolynomial := crc32.MakeTable(0xEDB88320)
		hash := crc32.New(tablePolynomial)

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
		shaa := header.File[i].Footer.Sha1[:]

		hashInBytes := hash.Sum(nil)[:]
		s := hex.EncodeToString(hashInBytes)
		n, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			panic(err)
		}
		crc := uint32(n)

		if header.File[i].CRC == crc {
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
