package source

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//FileReader utlilty helps you read file system
type FileReader interface {
	ReadDirs(root string) ([]string, error)
	ReadfilesWithExtension(root string, extn string) ([]string, error)
	ReadFileAsString(path string) (string, error)
}

// fileReaderImpl  implmentation of FileReader
type fileReaderImpl struct {
}

//ReadDirs reads all sub directories in a directory
func (f fileReaderImpl) ReadDirs(root string) ([]string, error) {
	var dirList []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		// check if it is a regular file (not dir)
		if info.Mode().IsDir() {
			dirList = append(dirList, info.Name())
		}
		return nil
	})
	return dirList, err
}

// ReadfilesWithExtension just for testing
func (f fileReaderImpl) ReadfilesWithExtension(root string, extn string) ([]string, error) {
	return filepath.Glob(filepath.Join(root , "*" + extn))
}

func (f fileReaderImpl) ReadFileAsString(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil

}
