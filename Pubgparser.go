package PubgParser

import (
	"log"
	"os"
	"fmt"
	"encoding/binary"
	"bytes"
	"encoding/json"
)

// Kill Item struct
type Kill struct {
	KillerNetId string
	KillerName string
	VictimNetId string
	VictimName string
}

// Read a file in the "kill" format, parse and return a friendly struct from JSON
func ParseKillFile (filename string) Kill {
	data := ReadPubgFile(filename)

	var killRecord Kill
	err := json.Unmarshal(data, &killRecord)
	JsonErrorHandler(err, data)

	return killRecord

}

// Read any file in the standard format
// The format is defined as a 32-bit unsigned int (dword) defining the length of the string
// and then the string terminated by null (\x00)
// The file is un-obfuscated by processing each byte as: b = (b + 1) & 255
func ReadPubgFile(filename string) []byte {
	f, err := os.Open(filename)
	ErrorHandler(err)

	lengthData := make([]byte, 4)
	read, err := f.Read(lengthData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Length Read %v, result: %v\n", read, lengthData[:4])

	length := binary.LittleEndian.Uint32(lengthData[0:4])
	fmt.Printf("Total file length length: %v bytes\n", length)

	data := make([]byte, length)
	read, err = f.Read(data)
	ErrorHandler(err)

	data = bytes.Trim(data, "\x00")

	var result []byte
	for _, b := range data {
		b = (b + 1) & 255
		result = append(result, b)
	}

	s := string(result[:read])
	fmt.Printf("String output: %v\n", s)

	return result
}

// Same as above but with a longer (6-byte) length field??
func ReadPubgFileLong(filename string) []byte {
	f, err := os.Open(filename)
	ErrorHandler(err)

	lengthData := make([]byte, 6)
	read, err := f.Read(lengthData)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Length Read %v, result: %v\n", read, lengthData[:6])

	length := binary.LittleEndian.Uint32(lengthData[0:6])
	fmt.Printf("Total file length length: %v bytes\n", length)

	data := make([]byte, length)
	read, err = f.Read(data)
	ErrorHandler(err)

	data = bytes.Trim(data, "\x00")

	var result []byte
	for _, b := range data {
		b = (b + 1) & 255
		result = append(result, b)
	}

	s := string(result[:read])
	fmt.Printf("String output: %v\n", s)

	return result
}

func ErrorHandler (err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func JsonErrorHandler (err error, j []byte) {
	if err != nil {
		good := json.Valid(j)
		if good == false {
			fmt.Print("Shitty fucking json brah")
		}
		log.Fatal(err)
	}

}
