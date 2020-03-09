package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var startDir string
var linkDir string
var extensions = make(map[string]int)

func main() {
	flag.StringVar(&startDir, "sd", ".", "starting source directory")
	flag.StringVar(&linkDir, "ld", "CPIO", "linking directory")
	flag.Parse()

	log.Printf("Walking tree based at %q\n", startDir)

	filepath.Walk(startDir, walker)
}

func walker(path string, info os.FileInfo, err error) error {

	if strings.HasSuffix(path, "CPIO") {
		return filepath.SkipDir
	}

	if info.IsDir() {
		return nil
	}

	e := filepath.Ext(path)

	if len(e) < 2 {
		return nil
	}

	extensions[e]++

	if extensions[e] > 10 {
		log.Println("NUMEROUS")
	}

	newpath := filepath.Join(linkDir, e[1:], path)

	os.MkdirAll(filepath.Dir(newpath), 0700)

	err = os.Link(path, newpath)
	if err != nil {
		return err
	}
	log.Println("hard linked", path, newpath)

	return nil
}
