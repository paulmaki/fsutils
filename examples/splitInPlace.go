package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

// MODE
// 1. number of times
// 2. by size

var filename string
var numberOfParts int

func init() {
	flag.StringVar(&filename, "f", "", "Filename to split")
	flag.IntVar(&numberOfParts, "n", 5, "Number of parts to split")
	flag.Parse()

	if filename == "" {
		log.Println("no filename provided")
		os.Exit(1)
	}
}

func main() {

	f, err := NewFile(filename)
	if err != nil {
		log.Panicln(err)
	}

	err = f.Split(numberOfParts)
	if err != nil {
		log.Panicln(err)
	}
}

func NewFile(path string) (f *File, err error) {
	log.Println("OPENING", path)
	f = new(File)
	f.path = path
	f.f, err = os.OpenFile(path, os.O_RDWR, 0666)
	return f, err
}

type File struct {
	f    *os.File
	path string
}

func (f File) Size() (int64, error) {
	st, err := os.Stat(f.path)
	if err != nil {
		return 0, err
	}
	return st.Size(), nil
}

func (f *File) Split(nParts int) error {
	st, err := os.Stat(f.path)
	if err != nil {
		return err
	}
	si := st.Size()

	chunkSize := si / int64(nParts)
	log.Println("TOTAL", si)
	log.Println("CHUNK", chunkSize)

	numParts := int(si / chunkSize)
	log.Println("NPARTS", numParts)

	// Seek to end of file
	pos, err := f.f.Seek(0, 2)
	if err != nil {
		return err
	}

	var nwritten int64

	for idx := range make([]int, numParts) {

		partNum := numParts - idx

		log.Println("PART NUM", partNum)
		pos, err = f.f.Seek(-chunkSize, 1)
		if err != nil {
			return err
		}

		segmentFileName := fmt.Sprintf("%s.%03d", f.path, partNum)

		// Last segment can be move in place
		if partNum == 1 {
			err = os.Rename(f.path, segmentFileName)
			if err != nil {
				return err
			}
			break
		}

		var outFile *os.File
		outFile, err = os.Create(segmentFileName)

		log.Println("NOW AT 1", pos, segmentFileName)

		// Copy to end of file
		nwritten, err = io.CopyN(outFile, f.f, chunkSize)
		if err != nil {
			return err
		}

		pos, err = f.f.Seek(0, 1)
		if err != nil {
			return err
		}

		log.Printf("Bytes written %dMB\n", nwritten>>20)
		log.Println("NOW AT 2", pos)

		err = f.f.Sync()
		if err != nil {
			return err
		}

		// Jump back again
		pos, err = f.f.Seek(-nwritten, 1)
		if err != nil {
			return err
		}

		log.Println("SNAPPED BACK TO", pos)

		// Truncate at current position
		err = f.f.Truncate(pos)
		if err != nil {
			return err
		}
	}

	return nil
}
