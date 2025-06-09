package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	fmt.Printf("Torrent files given: %v\n", args[1:])
}
