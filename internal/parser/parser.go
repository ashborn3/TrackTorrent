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

func NewTorrentFile(pathtofile string) (*TorrentFile, error) {
	var torrFile TorrentFile
	fData, err := os.ReadFile(pathtofile)
	if err != nil {
		return nil, fmt.Errorf("error reading torrent file: %v", err.Error())
	}
	err = bencode.Unmarshal(bytes.NewReader(fData), &torrFile)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling torrent file: %v", err.Error())
	}

	return &torrFile, nil
}

func (tf *TorrentFile) CalcInfoHash() ([]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, tf.Info)
	if err != nil {
		return nil, fmt.Errorf("error bencode marshalling Info struct, shouldn't have happened")
	}

	hasher := sha1.New()
	hasher.Write(buf.Bytes())

	return hasher.Sum(nil), nil
}

func (tf *TorrentFile) GetPiecesAsArray() ([][]byte, error) {
	numPieces := (len(tf.Info.Pieces) + 19) / 20 // round up division
	pieces := make([][]byte, numPieces)
	j := 0
	pieceLength := 20 // SHA-1 hash length
	for i := 0; i < len(tf.Info.Pieces); i += pieceLength {
		pieces[j] = []byte(tf.Info.Pieces[i : i+pieceLength])
		j = j + 1
	}

	return pieces, nil
}

type TrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    []byte `bencode:"peers"`
}

func DecodeTrackerResponse(body *bytes.Buffer) (*TrackerResponse, error) {
	var trkrRsp TrackerResponse
	err := bencode.Unmarshal(body, &trkrRsp)
	if err != nil {
		return nil, err
	}
	return &trkrRsp, err
}
