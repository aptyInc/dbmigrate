package source

import (
	"github.com/spf13/afero"
	"path/filepath"
)

// ReaderImplementation  implementation of FileReader
type ReaderImplementation struct {
	fs afero.Fs
}

//ReadDirs reads all sub directories in a directory
func (f ReaderImplementation) ReadDirs(root string) ([]string, error) {
	var dirList []string
	files, err := afero.ReadDir(f.fs,root)
	if err != nil {
		return dirList,err
	}

	for _, f := range files {
		dirList = append(dirList, f.Name())
	}

	return dirList, err
}

// ReadFilesWithExtension just for testing
func (f ReaderImplementation) ReadFilesWithExtension(root string, extension string) ([]string, error) {
	return afero.Glob(f.fs,filepath.Join(root , "*" + extension))
}

func (f ReaderImplementation) ReadFileAsString(path string) (string, error) {
	b, err := afero.ReadFile(f.fs,path)
	if err != nil {
		return "", err
	}
	return string(b), nil

}