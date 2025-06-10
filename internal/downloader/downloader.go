package downloader

import "tracktorrent/internal/parser"

type TorrentDownloader struct {
	TorrentFile *parser.TorrentFile
	Pieces      [][]byte
}

func NewDownloader(torrentfile *parser.TorrentFile) (*TorrentDownloader, error) {
	torrdwnldr := TorrentDownloader{
		TorrentFile: torrentfile,
	}
	var err error
	torrdwnldr.Pieces, err = torrdwnldr.TorrentFile.GetPiecesAsArray()
	if err != nil {
		return nil, err
	}
	return &torrdwnldr, nil
}
