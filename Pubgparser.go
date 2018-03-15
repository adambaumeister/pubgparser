package PubgParser

import (
	"log"
	"os"
	"fmt"
	"encoding/binary"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
)

// Entire replay struct
type Replay struct {
	Kills []Kill
	Knockdowns []Groggy
	Summary []ReplaySummary
}

// Kill Item struct
type Kill struct {
	KillerNetId string
	KillerName  string
	VictimNetId string
	VictimName  string
}
//ReplaySummary
type ReplaySummary struct {
	MatchId            string
	BIsServerRecording bool
	RecordUserId       string
	RecordUserNickName string
	MapName            string
	WeatherName        string
	RegionName         string
	NumPlayers         int
	NumTeams           int

	PlayerStateSummaries []PlayerStateSummary
}
// Player state summary, contained within replay summary files
type PlayerStateSummary struct {
	PlayerName              string
	TeamNumber              int
	Ranking                 int
	HeadShots               int
	NumKills                int
	TotalGivenDamages       float64
	LongestDistanceKill     float64
	TotalMovedDistanceMeter float64
}
// Groggy, aka, knockdown struct
type Groggy struct {
	InstigatorNetId			string
	InstigatorName			string
	VictimNetId				string
	VictimName				string
}

// Read a file in the "kill" format, parse and return a friendly struct from JSON
func ParseKillFile (filename string) Kill {
	data := ReadPubgFile(filename)

	var killRecord Kill
	err := json.Unmarshal(data, &killRecord)
	JsonErrorHandler(err, data)

	return killRecord

}

// Read a file in the "replay summary" format and return a friendly structerino
func ParseSummaryFile (filename string) ReplaySummary {
	var result ReplaySummary
	data := ReadPubgFile(filename)
	err := json.Unmarshal(data, &result)
	JsonErrorHandler(err, data)

	return result
}

func ParseGroggyFile (filename string) Groggy {
	var result Groggy
	data := ReadPubgFile(filename)
	err := json.Unmarshal(data, &result)
	JsonErrorHandler(err, data)

	return result
}

func ReadPubgDirectory (dir string) Replay {
	dirList, err := ioutil.ReadDir(dir)
	ErrorHandler(err)
	var r Replay
	for _, file := range dirList {
		if strings.Contains(file.Name(), "kill") {
			path := fmt.Sprintf("%v\\%v", dir, file.Name())
			r.Kills = append(r.Kills, ParseKillFile(path))
		} else if strings.Contains(file.Name(), "groggy") {
			path := fmt.Sprintf("%v\\%v", dir, file.Name())
			fmt.Printf("%v\n", path)
			r.Knockdowns = append(r.Knockdowns, ParseGroggyFile(path))
		} else if strings.Contains(file.Name(), "ReplaySummary") {
			path := fmt.Sprintf("%v\\%v", dir, file.Name())
			r.Summary = append(r.Summary, ParseSummaryFile(path))
		}
	}
	return r
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
	//fmt.Printf("Total file length length: %v bytes\n", length)

	data := make([]byte, length)
	read, err = f.Read(data)
	ErrorHandler(err)

	data = bytes.Trim(data, "\x00")

	var result []byte
	for _, b := range data {
		b = (b + 1) & 255
		result = append(result, b)
	}

	//s := string(result[:read])
	//fmt.Printf("String output: %v\n", s)

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
