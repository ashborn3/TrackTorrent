package main

import (
	"fmt"
	"os"
	"tracktorrent/internal/downloader"
	"tracktorrent/internal/parser"
)

func main() {
	args := os.Args
	fmt.Printf("Torrent files given: %v\n", args[1:])

	torrFile, err := parser.NewTorrentFile(args[1])
	if err != nil {
		panic(err.Error())
	}

	torrdwnldr, err := downloader.NewDownloader(torrFile)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%s\n", torrdwnldr.TrackerResponse.Peers)
}
