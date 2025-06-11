package main

import (
	"encoding/json"
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

	b, err := json.MarshalIndent(torrdwnldr, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
