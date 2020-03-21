package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var sc = bufio.NewScanner(os.Stdin)

	var filesOnly bool
	var dirsOnly bool
	flag.BoolVar(&filesOnly, "files", false, "files only")
	flag.BoolVar(&dirsOnly, "dirs", false, "dirs only")
	flag.Parse()

	var err error
	var abspath string
	var st os.FileInfo

	for sc.Scan() {
		st, err = os.Stat(sc.Text())
		if err != nil {
			continue
		}

		if filesOnly && st.IsDir() {
			continue
		}

		if dirsOnly && !st.IsDir() {
			continue
		}

		abspath, err = filepath.Abs(sc.Text())
		if err != nil {
			continue
		}
		fmt.Println(abspath)
	}
}
