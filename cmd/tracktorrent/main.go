package main

import (
	"fmt"
	"os"
	"tracktorrent/internal/parser"
)

func main() {
	args := os.Args
	fmt.Printf("Torrent files given: %v\n", args[1:])

	torrFile, err := parser.NewTorrentFile(args[1])
	if err != nil {
		panic(err.Error())
	}

	// shahash, _ := torrFile.CalcInfoHash()

	mapie, err := torrFile.GetPiecesAsHexArray()

	for idx, val := range mapie {
		fmt.Printf("%d: %s\n", idx, val)
	}
}
