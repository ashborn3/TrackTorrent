package parser

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce     string     `bencode:"announce"`
	AnnounceList [][]string `bencode:"announce-list,omitempty"`
	Info         InfoDict   `bencode:"info"`
	Comment      string     `bencode:"comment,omitempty"`
	CreatedBy    string     `bencode:"created by,omitempty"`
	CreationDate int64      `bencode:"creation date,omitempty"`
	Encoding     string     `bencode:"encoding,omitempty"`
}

type InfoDict struct {
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"` // raw string of SHA1 hashes
	Private     int    `bencode:"private,omitempty"`

	Name   string `bencode:"name"`
	Length int64  `bencode:"length,omitempty"` // For single file
	Files  []File `bencode:"files,omitempty"`  // For multi-file torrents
}

type File struct {
	Length int64    `bencode:"length"`
	Path   []string `bencode:"path"` // Path is an array of directories + filename
}

func ParseTorrentFile(pathtofile string) TorrentFile {
	var torrFile TorrentFile
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

func (tf TorrentFile) CalcInfoHash() []byte {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, tf.Info)
	if err != nil {
		panic("Error calculating info hash, shouldn't have happened!")
	}

	hasher := sha1.New()
	hasher.Write(buf.Bytes())

	return hasher.Sum(nil)
}
