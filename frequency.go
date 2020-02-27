package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	filepath.Walk(".", walker)
}

var counts = make(map[string]int)

func walker(path string, info os.FileInfo, err error) error {
	// SPlit allparts

	pathparts := strings.Split(path, "/")

	for idx, part := range pathparts {
		counts[part] += 1

		if counts[part] > 1000 {
			log.Println(part, counts[part], pathparts[:idx+1])
		}
	}
	return nil
}
