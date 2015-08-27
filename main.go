package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/geekbros/DiskAnalyzer/diskanalyzer"
)

var (
	path = flag.String("path", "", "Path to directory to work with.")

	modeDirStructure = flag.Bool("dir-struct", false, `File structure mode. Returns hierarchical
		representation of given directory with given depth, as well internal files sizes.`)
	depth = flag.Int("depth", 3, "Depth of walk for File structure mode.")

	modeLargeFiles = flag.Bool("large-dirs", false, `Large directories mode. Returns all subdirectories of given
		directory, which size is bigger than given max size.`)
	maxSize = flag.Int("max-size", 1000, "Files with size bigger than this value are considered as large files.")
)

func main() {
	flag.Parse()
	if *path == "" {
		panic("Path flag is required and can't be empty.")
	}
	var jsonbytes []byte
	if *modeDirStructure {
		dir, err := diskanalyzer.GetDirectoryStructureWithDepth(*path, *depth)
		if err != nil {
			panic(err)
		}
		jsonbytes, err = json.Marshal(&dir)
		if err != nil {
			panic(err)
		}
	} else if *modeLargeFiles {
		m, err := diskanalyzer.GetLargeDirectories(*path)
		if err != nil {
			panic(err)
		}
		jsonbytes, err = json.Marshal(&m)
		if err != nil {
			panic(err)
		}
	} else {
		panic("One of modes has to be choosen.")
	}
	fmt.Println(string(jsonbytes))
}
