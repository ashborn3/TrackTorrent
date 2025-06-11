package downloader

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"tracktorrent/internal/config"
	"tracktorrent/internal/parser"
	"tracktorrent/internal/random"
)

type TorrentDownloader struct {
	PeerId          string
	TorrentFile     *parser.TorrentFile
	TrackerResponse *parser.TrackerResponse
	Conns           []*Conn
	PieceHashes     [][]byte
	Pieces          [][]byte
	Peers           []string
	Uploaded        int
	Downloaded      int
	Left            int64
	Compact         int
}

type Conn struct {
	ProtocolLength int
	Protocol       string
	InfoHash       []byte
	PeerId         []byte
	Conn           *net.Conn
	Alive          bool
}

func NewDownloader(torrentfile *parser.TorrentFile) (*TorrentDownloader, error) {
	torrdwnldr := TorrentDownloader{
		TorrentFile: torrentfile,
		PeerId:      random.RandStringBytes(20),
		Uploaded:    0,
		Downloaded:  0,
		Left:        torrentfile.Info.Length,
	}
	var err error
	torrdwnldr.PieceHashes, err = torrdwnldr.TorrentFile.GetPiecesAsArray()
	if err != nil {
		return nil, err
	}

	err = torrdwnldr.requestPeerList()
	if err != nil {
		return nil, err
	}

	torrdwnldr.Conns = make([]*Conn, len(torrdwnldr.TrackerResponse.Peers))

	for idx := range len(torrdwnldr.TrackerResponse.Peers) {
		err = torrdwnldr.handShakeWithPeer(idx)
		if err != nil {
			torrdwnldr.Conns[idx] = nil
		}
	}

	return &torrdwnldr, nil
}

func (td *TorrentDownloader) requestPeerList() error {
	baseUrl := td.TorrentFile.Announce

	infoHash, err := td.TorrentFile.CalcInfoHash()
	if err != nil {
		return err
	}
	encodedInfoHash := ""
	for _, b := range infoHash {
		encodedInfoHash += fmt.Sprintf("%%%02x", b)
	}

	getParams := url.Values{}
	getParams.Add("peer_id", td.PeerId)
	getParams.Add("port", strconv.Itoa(config.PORT))
	getParams.Add("uploaded", strconv.Itoa(td.Uploaded))
	getParams.Add("downloaded", strconv.Itoa(td.Downloaded))
	getParams.Add("left", strconv.FormatInt(td.Left, 10))
	getParams.Add("compact", strconv.Itoa(td.Compact))

	finalUrl := fmt.Sprintf("%s?info_hash=%s&%s", baseUrl, encodedInfoHash, getParams.Encode())

	resp, err := http.Get(finalUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(bodyBytes)
	td.TrackerResponse, err = parser.DecodeTrackerResponse(buf)
	if err != nil {
		return err
	}

	return nil
}

func (td *TorrentDownloader) handShakeWithPeer(peeridx int) error {
	peer := td.TrackerResponse.Peers[peeridx]

	conn, err := net.Dial("tcp", net.JoinHostPort(peer.Ip, strconv.Itoa(peer.Port)))
	if err != nil {
		return err
	}

	infohash, err := td.TorrentFile.CalcInfoHash()
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	pstr := "BitTorrent protocol"
	buf.WriteByte(byte(len(pstr))) // 1 byte
	buf.WriteString(pstr)          // 19 bytes
	buf.Write(make([]byte, 8))     // 8 bytes
	buf.Write(infohash[:])         // 20 bytes
	buf.Write([]byte(td.PeerId))   // 20 bytes

	handshake := buf.Bytes()

	_, err = conn.Write([]byte(handshake))
	if err != nil {
		return err
	}

	resp := make([]byte, 68)
	_, err = io.ReadFull(conn, resp) // ensures all 68 bytes are read
	if err != nil {
		return err
	}

	connStruct := Conn{
		ProtocolLength: int(resp[0]),
		Protocol:       string(resp[1 : 1+int(resp[0])]),
		InfoHash:       resp[1+int(resp[0])+8 : 1+int(resp[0])+8+20],
		PeerId:         resp[1+int(resp[0])+8+20:],
		Conn:           &conn,
		Alive:          true,
	}

	if !bytes.Equal(connStruct.InfoHash, infohash[:]) {
		return fmt.Errorf("mismatched infohash in handshake")
	}

	td.Conns[peeridx] = &connStruct

	return nil
}
