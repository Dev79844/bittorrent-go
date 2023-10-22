package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"unicode"
	"strings"
	bencode "github.com/jackpal/bencode-go"
)

type Torrent struct {
	Announce string `json:"announce"`
	Info     struct {
		Length       int    `json:"length"`
		Name         string `json:"name"`
		Piece_length int    `json:"piece_length"`
		Pieces       string `json:"pieces"`
	} `json:"info"`
}

func bencodeNums(bencodedString string) (interface{},error){
	eIndex := strings.Index(bencodedString,"e")
	if eIndex == -1 {
		return nil,fmt.Errorf("invalid string")
	}
	num,err := strconv.Atoi(bencodedString[1:eIndex])
	if err!=nil{
		return nil,err
	}
	return num,nil
}

func bencodeStrings(bencodedString string) (interface{},error){
	var firstColonIndex int

	for i := 0; i < len(bencodedString); i++ {
		if bencodedString[i] == ':' {
			firstColonIndex = i
			break
		}
	}

	lengthStr := bencodedString[:firstColonIndex]

	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", err
	}

	return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil
}

func decodeBencode(bencodedString string) (interface{}, error) {
	if bencodedString[0] == 'd'{
		reader:= strings.NewReader(bencodedString)
		val,err := bencode.Decode(reader)
		if err!=nil{
			return nil, err
		}
		return val,nil
	}
	if bencodedString[0] == 'l' {
		reader:= strings.NewReader(bencodedString)
		val,err := bencode.Decode(reader)
		if err!=nil{
			return nil, err
		}
		return val,nil
	}
	if bencodedString[0] == 'i'{
		value,err := bencodeNums(bencodedString)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		return value,nil
	}
	if unicode.IsDigit(rune(bencodedString[0])) {
		value,err := bencodeStrings(bencodedString)
		if err != nil {
			return nil, fmt.Errorf(err.Error())
		}
		return value,nil
	} else {
		return "", fmt.Errorf("Only strings are supported at the moment")
	}
}

func readTorrentFile(filename string) (interface{},error){
	f,err := os.ReadFile(filename)
	if err!=nil{
		return nil, fmt.Errorf(err.Error())
	}
	decoded,err := decodeBencode(string(f))

	jsonOutput, _ := json.Marshal(decoded)
	var torrent Torrent
	json.Unmarshal(jsonOutput,&torrent)

	return torrent,nil
}

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]
		
		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else if command == "info"{
		filename := os.Args[2]

		val,err := readTorrentFile(filename)
		if err!= nil{
			fmt.Println(err)
			return
		}
		jsonOutput, _ := json.Marshal(val)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
