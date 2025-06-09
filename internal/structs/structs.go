package structs

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
