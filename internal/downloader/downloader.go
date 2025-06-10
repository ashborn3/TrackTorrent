package downloader

import (
	"bytes"
	"fmt"
	"io"
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
	PieceHashes     [][]byte
	Pieces          [][]byte
	Peers           []string
	Uploaded        int
	Downloaded      int
	Left            int64
	Compact         int
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
