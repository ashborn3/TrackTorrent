package parser

import (
	"bytes"
	"fmt"
	"os"
	"tracktorrent/internal/structs"

	"github.com/jackpal/bencode-go"
)

func ParseTorrentFile(pathtofile string) structs.TorrentFile {
	var torrFile structs.TorrentFile
	fData, err := os.ReadFile(pathtofile)
	if err != nil {
		panic(fmt.Sprintf("Error reading given file: %s", err.Error()))
	}
	err = bencode.Unmarshal(bytes.NewReader(fData), &torrFile)
	if err != nil {
		panic("Error unmarshalling bencode file into struct")
	}

	return torrFile
}
