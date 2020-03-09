package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

var startDir string

func main() {
	flag.StringVar(&startDir, "d", ".", "starting directory")
	flag.Parse()

	log.Printf("Walking tree based at %q\n", startDir)

	filepath.Walk(startDir, walker)
}

func walker(path string, info os.FileInfo, err error) error {
	log.Printf("Current path:%s isdir:%t\n", path, info.IsDir())
	return nil
}
