package fsutils

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var truncateFiles bool
var numRemoves = make(map[string]int)
var c Counts
var removeDirParts = make(map[string]bool)
var config Config

type Counts struct {
	remove int
	dirs   int
	total  int
}

func (c Counts) String() string {
	return fmt.Sprintf("Counts:Dir:%d [File %d / %d]", c.dirs, c.remove, c.total)
}

type Config struct {
	minFilesBeforeQuestion int
	maxExtensionLength     int
	clearNonExtensionFiles bool
	truncateFiles          bool
}

func init() {
	config = NewConfig()
}

func NewConfig() (config Config) {
	flag.IntVar(&config.minFilesBeforeQuestion, "min", 100, "Minimum observed files before deleting")
	flag.IntVar(&config.maxExtensionLength, "maxlen", 0, "max extension length [0 to ignore]")
	flag.BoolVar(&config.clearNonExtensionFiles, "nonext", false, "Automatically clear files without extensions")
	flag.BoolVar(&config.truncateFiles, "truncate", false, "Truncate the file before removing")
	flag.Parse()

	return
}

func RemoveByExtension(basepath string) {

	if len(basepath) == 0 {
		print("Failed to get default path")
		flag.Usage()
		os.Exit(1)
	}

	deleteMap := make(map[string]bool)

	var ext string

	scanner := bufio.NewScanner(os.Stdin)

	extensionFiles := make(map[string][]string)

	filepath.Walk(basepath, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			c.dirs++
			fmt.Fprintf(os.Stderr, "Dir %v %q\n", c, path)

			return nil
		}

		c.total++

		ext = filepath.Ext(path)

		if config.clearNonExtensionFiles && len(ext) <= 1 {
			err := removePath(path)
			println("Removing non ext", path)
			if err != nil {
				log.Println(err)
			}

			return err
		}

		if config.maxExtensionLength > 0 && len(ext) > config.maxExtensionLength+1 {
			err := removePath(path)
			println("Removing long extension ext", path)
			return err
		}

		extensionFiles[ext] = append(extensionFiles[ext], path)

		if len(extensionFiles[ext]) < config.minFilesBeforeQuestion {
			return nil
		}

		var removefileByDirPart bool

		// Check all dir parts
		for _, dpart := range strings.Split(filepath.Dir(path), "/") {
			if _, oktoremovePart := removeDirParts[dpart]; oktoremovePart {
				log.Println("Remove by dir part", dpart, path)
				err := removePath(path)
				if err != nil {
					log.Println(err)
				}
				removefileByDirPart = true
				continue
			}
		}

		if removefileByDirPart {
			return nil
		}

		if _, ok := deleteMap[ext]; !ok {
			// println("new ext", ext)

			fmt.Fprintf(os.Stderr, "\n\nCurrent file: %q\n", path)
			fmt.Fprintf(os.Stderr, "[[%q]]\n", ext)
			fmt.Fprintf(os.Stderr, "[de] file extension\n[dp] directory parts\n")
			scanner.Scan()

			switch scanner.Text() {
			case "de":
				log.Printf("deleting %s\n", ext)
				deleteMap[ext] = true

				for _, filepath := range extensionFiles[ext] {
					if err := removePath(filepath); err != nil {
						log.Println("Couldn't remove file", path)
					}
					println("Removed built up file", filepath)
				}
			case "dp":
				log.Println("Directory part")

				parts := strings.Split(filepath.Dir(path), "/")

				var dirparts = make(map[string]string)

				for idx, part := range parts {
					dirparts[fmt.Sprintf("%d", idx)] = part
					log.Println(idx, part)
				}

				log.Println("<ENTER>")
				scanner.Scan()

				if dpart, ok := dirparts[scanner.Text()]; ok {
					log.Println("Want to skip dir part", dpart)
					removeDirParts[dpart] = true
				}

			default:
				fmt.Fprintf(os.Stderr, "keeping files with extension:%q\n", ext)
				deleteMap[ext] = false
			}
		}

		if deleteMap[ext] {
			if err := removePath(path); err != nil {
				log.Println(err)
			}
		}

		return nil
	})

	return
}

func removePath(path string) error {

	var err error

	if truncateFiles {
		err := os.Truncate(path, 0)
		if err != nil {
			return err
		}
	}

	err = os.Remove(path)
	if err != nil {
		return err
	}

	ext := filepath.Ext(path)
	numRemoves[ext]++
	c.remove++

	fmt.Fprintf(os.Stderr, "D | %q | %d %s\n", ext, numRemoves[ext], c)
	return nil
}
