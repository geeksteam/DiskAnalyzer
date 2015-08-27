package diskanalyzer

import (
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
)

const (
	// Value of max depth of filesystem walk function.
	// Used in module's GetDirectoryStructure function.
	maxDepth = 2

	// Value of directory's size in megabytes, above which it's considered as a large directory.
	maxSize = 20
)

// Directory type is a os's directory representation for GetDirectoryStructure function.
type Directory struct {
	fullpath string
	size     int64
	subdirs  []Directory
}

// GetFullDirSize - Get Directory size include subdirs
func GetFullDirSize(fullpath string) (result int64, err error) {
	err = filepath.Walk(fullpath, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			result += info.Size()
		}
		return nil
	})
	result = int64(math.Ceil(float64(result) / float64(1000000)))
	return
}

func combineMaps(from, to *map[string]string) {
	for k := range *from {
		(*to)[k] = (*from)[k]
	}
}

// GetLargeDirectories returns a map of directories, which are bigger than maxSize.
func GetLargeDirectories(root string) (map[string]string, error) {
	return getDirectoriesBiggerThan(maxSize, root)
}

func getDirectoriesBiggerThan(size int64, root string) (result map[string]string, err error) {
	var dir *os.File
	if dir, err = os.Open(root); err != nil {
		return
	}
	defer dir.Close()

	result = make(map[string]string)

	var dirInfo os.FileInfo
	if dirInfo, err = dir.Stat(); err == nil && dirInfo.IsDir() {
		var dirSize int64
		if dirSize, err = GetFullDirSize(root); err == nil && dirSize >= size {
			result[root] = fmt.Sprintf("%v", dirSize)
			var files []os.FileInfo
			if files, err = dir.Readdir(0); err == nil {
				for _, v := range files {
					var subMap map[string]string
					if subMap, err = getDirectoriesBiggerThan(size, path.Join(root, v.Name())); err == nil {
						combineMaps(&subMap, &result)
					}
				}
			}
		}
	}
	return
}

// GetDirectoryStructure returns hierarchical representation of os's directories,
// as well as their sizes. Depth of walk is set by maxDepth constant.
func GetDirectoryStructure(fullpath string) (result Directory, err error) {
	return getDirectoryStructureWithDepth(fullpath, maxDepth)
}

func getDirectoryStructureWithDepth(fullpath string, depth int) (result Directory, err error) {
	var fileSize int64
	if fileSize, err = GetFullDirSize(fullpath); err != nil {
		return
	}

	var currentDir *os.File
	if currentDir, err = os.Open(fullpath); err != nil {
		return
	}
	defer currentDir.Close()

	var files []os.FileInfo
	if files, err = currentDir.Readdir(0); err != nil {
		return
	}

	var directories []Directory

	for _, v := range files {
		if v.IsDir() && depth > 0 {
			var subDir Directory
			subDirPath := path.Join(fullpath, v.Name())
			if subDir, err = getDirectoryStructureWithDepth(subDirPath, depth-1); err == nil {
				directories = append(directories, subDir)
			} else {
				return
			}
		}
	}

	result = Directory{fullpath, fileSize, directories}
	return result, nil
}
