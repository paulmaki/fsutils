package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	log.Panicln("You realize this wipes directories?")
}

func WipeDirectory() {
	entries, err := ioutil.ReadDir(".")
	if err != nil {
		return
	}

	var ent Named
	var size int64

	var xbuf = bytes.Repeat([]byte{0x1}, 5<<20)

	for _, ent = range entries {
		size = ent.Size()

		if size > 5<<20 {
			log.Println(ent.Name(), "TOO BIG")
			continue
		}

		log.Println(29, "OPEN", ent.Name())
		fd, err := os.OpenFile(ent.Name(), os.O_WRONLY, 0)
		if err != nil {
			log.Println(31, err)
			log.Panicln()
		}
		nb, err := fd.Write(xbuf[:size])
		log.Println(35, nb, err, ent.Name())
		fd.Close()
	}
}

type Named interface {
	Name() string
	Size() int64
}
