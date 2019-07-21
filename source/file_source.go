package source

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//FileSource file implementation for MigrationSource
type FileSource struct {
	schemas                 []string
	schemaUpMigrationsMap   map[string]map[int]string
	schemaDownMigrationsMap map[string]map[int]string
	sortedVersions          []int
	location                string
	reader                  FileReader
}

func getFileDetails(file string) (version int, isUp bool, isDown bool, err error) {
	parts := strings.Split(file, ".")
	if len(parts) != 3 {
		return -1, false, false, fmt.Errorf("file Name format Not correct, issue with number of separators in %s", file)
	}
	splits := strings.Split(parts[0], "-")
	if len(splits) != 2 {
		return -1, false, false, fmt.Errorf("file Name format Not correct, issue with number of separators in %s", file)
	}
	version, err2 := strconv.Atoi(splits[0])
	if err2 != nil {
		fmt.Println(err2)
		return -1, false, false, fmt.Errorf("file Name format Not correct, issue with version number of %s", file)
	}

	isUp = parts[1] == "UP"
	isDown = parts[1] == "DOWN"

	return version, isUp, isDown, nil
}

func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//GetFileSource returns as Filesource Object
func GetFileSource(baseLocation string, fs FileReader) (MigrationSource, error) {
	fmt.Println("Source Directory:",baseLocation)
	folders, err := fs.ReadDirs(baseLocation)
	if err != nil {
		return nil, err
	}
	var schemaUPFilesMap = make(map[string]map[int]string)
	var schemaDownFilesMap = make(map[string]map[int]string)
	var versions []int
	for _, folder := range folders {
		files, err := fs.ReadFilesWithExtension(filepath.Join(baseLocation, folder), ".sql")
		if err != nil {
			return nil, err
		}
		for _, file := range files {

			version, isUP, isDown, err1 := getFileDetails(filepath.Base(file))
			if err1 != nil {
				return nil, err1
			}
			if isUP {
				if schemaUPFilesMap[folder] == nil {
					schemaUPFilesMap[folder] = make(map[int]string)
				}
				schemaUPFilesMap[folder][version] = file
			}
			if isDown {
				if schemaDownFilesMap[folder] == nil {
					schemaDownFilesMap[folder] = make(map[int]string)
				}
				schemaDownFilesMap[folder][version] = file
			}
			if isUP || isDown {
				versions = append(versions, version)
			}
		}

	}
	versions = unique(versions)
	sort.Ints(versions)
	source := FileSource{
		schemas:                 folders,
		schemaUpMigrationsMap:   schemaUPFilesMap,
		schemaDownMigrationsMap: schemaDownFilesMap,
		sortedVersions:          versions,
		location:                baseLocation,
		reader:                  fs,
	}
	return &source, nil
}

// GetSchemaList returns list of folders which will be used ad shcema names
func (fs *FileSource) GetSchemaList() ([]string, error) {
	return fs.schemas, nil
}

// GetSortedVersions gets the list of verrsion numbers
func (fs *FileSource) GetSortedVersions(schema string) ([]int, error) {
	return fs.sortedVersions, nil
}

// GetMigrationUpFile returns the Migration Up Files of specifed version
func (fs *FileSource) GetMigrationUpFile(schema string, version int) (string, string, error) {
	filePath := filepath.Join(fs.location,fs.schemaUpMigrationsMap[schema][version])
	// TODO Bug here remove schema from join
	fmt.Println(filePath)
	contents, err := fs.reader.ReadFileAsString(filePath)
	return fs.schemaUpMigrationsMap[schema][version], contents, err
}

// GetMigrationDownFile returns the Migration Down Files of specifed version
func (fs *FileSource) GetMigrationDownFile(schema string, version int) (string, string, error) {
	filePath := filepath.Join(fs.location, fs.schemaDownMigrationsMap[schema][version])
	contents, err := fs.reader.ReadFileAsString(filePath)
	return fs.schemaDownMigrationsMap[schema][version], contents, err
}
