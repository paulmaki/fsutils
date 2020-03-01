package main

import (
	"flag"
	"log"
	"os"

	"github.com/paulmaki/fsutils"
)

func main() {

	var config fsutils.Config

	flag.IntVar(&config.MinFilesBeforeQuestion, "min", 100, "Minimum observed files before deleting")
	flag.IntVar(&config.MaxExtensionLength, "maxlen", 0, "max extension length [0 to ignore]")
	flag.BoolVar(&config.ClearNonExtensionFiles, "nonext", false, "Automatically clear files without extensions")
	flag.BoolVar(&config.TruncateFiles, "truncate", false, "Truncate the file before removing")
	flag.StringVar(&config.Basepath, "d", ".", "Base path")
	flag.Parse()

	err := fsutils.RemoveByExtension(config)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
