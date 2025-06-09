package main

import (
	"fmt"
	"os"
	"tracktorrent/internal/parser"
)

func main() {
	args := os.Args
	fmt.Printf("Torrent files given: %v\n", args[1:])

	torrFile := parser.ParseTorrentFile(args[1])

	fmt.Printf("%x\n", torrFile.CalcInfoHash())
}
