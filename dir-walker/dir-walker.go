package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Simple app that walks through all folders in a directory, and lists the files in the order of decreasing size

// Defining constants to convert size units
const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

type filePair struct {
	path     string
	fileinfo os.FileInfo
}

func IsHiddenFile(filename string) bool {
	return filename[0] == '.'
}

func IsHiddenDir(path string) bool {
	return strings.Contains(path, "/.")
}

func printHumanize(path string, size float64) {
	if size < KB {
		fmt.Printf("%s: %.2f B\n", path, size)
	} else if size > KB && size < MB {
		fmt.Printf("%s: %.2f KB\n", path, size/float64(KB))
	} else if size >= MB && size < GB {
		fmt.Printf("%s: %.2f MB\n", path, size/float64(MB))
	} else if size >= GB && size < TB {
		fmt.Printf("%s: %.2f GB\n", path, size/float64(GB))
	} else {
		fmt.Printf("%s: %.2f TB\n", path, size/float64(TB))
	}

}
func main() {

	pathPtr := flag.String("p", "/", "file path to explore")
	hiddenPtr := flag.Bool("s", false, "traverse hidden directories & files")
	limitPtr := flag.Int("c", 10, "line limit")

	flag.Parse()

	var files []filePair

	filepath.Walk(*pathPtr, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			//in case file gets "cleaned up" during walk
			return nil
		}
		if (IsHiddenFile(info.Name()) || IsHiddenDir(path)) && !(*hiddenPtr) {
			//skip hidden file/dir
			return nil
		}
		if info.IsDir() {
			//skip if we're looking at a dir
			return nil
		}
		files = append(files, filePair{path, info})
		return nil
	})
	sort.Slice(files, func(i, j int) bool {
		return files[i].fileinfo.Size() > files[j].fileinfo.Size()
	})

	count := 0
	for _, file := range files {
		printHumanize(file.path, float64(file.fileinfo.Size()))
		count += 1
		if count == *limitPtr {
			break
		}
	}
}
