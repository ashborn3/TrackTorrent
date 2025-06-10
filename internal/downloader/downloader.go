package downloader

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"tracktorrent/internal/config"
	"tracktorrent/internal/parser"
	"tracktorrent/internal/random"
)

type TorrentDownloader struct {
	PeerId      string
	TorrentFile *parser.TorrentFile
	PieceHashes [][]byte
	Pieces      [][]byte
	Peers       []string
	Uploaded    int
	Downloaded  int
	Left        int
	Compact     int
}

func NewDownloader(torrentfile *parser.TorrentFile) (*TorrentDownloader, error) {
	torrdwnldr := TorrentDownloader{
		TorrentFile: torrentfile,
		PeerId:      random.RandStringBytes(20),
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

	getParams := url.Values{}
	getParams.Add("info_hash", string(infoHash))
	getParams.Add("peer_id", td.PeerId)
	getParams.Add("port", strconv.Itoa(config.PORT))
	getParams.Add("uploaded", strconv.Itoa(td.Uploaded))
	getParams.Add("downloaded", strconv.Itoa(td.Downloaded))
	getParams.Add("left", strconv.Itoa(td.Left))
	getParams.Add("compact", strconv.Itoa(td.Compact))

	finalUrl := fmt.Sprintf("%s?%s", baseUrl, getParams.Encode())

	resp, err := http.Get(finalUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
