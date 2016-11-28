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
	Fullpath string      `json:"path"`
	Size     int64       `json:"size"`
	Subdirs  []Directory `json:"sub_dirs"`
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
	return GetDirectoriesBiggerThan(maxSize, root)
}

// GetDirectoriesBiggerThan returns a map of directories, which are bigger than size.
func GetDirectoriesBiggerThan(size int64, root string) (result map[string]string, err error) {
	var dir *os.File
	if dir, err = os.Open(root); err != nil {
		return
	}
	defer dir.Close()

	result = make(map[string]string)

	var dirInfo os.FileInfo
	if dirInfo, err = dir.Stat(); err == nil && dirInfo.IsDir() {
		var dirSize int64
		if dirSize, err = GetFullDirSize(root); err == nil && dirSize >= size*1000 {
			// Skipping root directory
			if root != "/" {
				result[root] = fmt.Sprintf("%v", int64(math.Ceil(float64(dirSize)/1000)))
			}

			var files []os.FileInfo
			if files, err = dir.Readdir(0); err == nil {
				for _, v := range files {
					var subMap map[string]string
					if subMap, err = GetDirectoriesBiggerThan(size, path.Join(root, v.Name())); err == nil {
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
	return GetDirectoryStructureWithDepth(fullpath, maxDepth)
}

// GetDirectoryStructureWithDepth returns hierarchical representation of os's directories with given depth,
// as well as their sizes.
func GetDirectoryStructureWithDepth(fullpath string, depth int) (result Directory, err error) {
	var fileSize int64
	if fileSize, err = GetFullDirSize(fullpath); err != nil {
		return
	}

	var currentDir *os.File
	defer currentDir.Close()
	if currentDir, err = os.Open(fullpath); err != nil {
		return
	}

	var files []os.FileInfo
	if files, err = currentDir.Readdir(0); err != nil {
		return
	}

	var directories []Directory

	for _, v := range files {
		if v.IsDir() && depth > 0 {
			var subDir Directory
			subDirPath := path.Join(fullpath, v.Name())
			if subDir, err = GetDirectoryStructureWithDepth(subDirPath, depth-1); err == nil {
				directories = append(directories, subDir)
			} else {
				return
			}
		}
	}

	result = Directory{fullpath, fileSize, directories}
	return result, nil
}
