package source

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/apty/dbmigrate/source/mocks"
)

// FileSourceInternalTestSuite Suite to test
type FileSourceTestSuite struct {
	suite.Suite
	mockFR *mocks.FileReader
}

func (suite *FileSourceTestSuite) SetupSuite() {
	suite.mockFR = new(mocks.FileReader)

}

func (suite *FileSourceTestSuite) TestValidFileSource() {
	const basePath = "/test"
	var dirs = []string{"sub1", "sub2"}
	var files = []string{"03-file1.UP.sql", "03-file1.DOWN.sql", "02-file2.UP.sql", "02-file2.DOWN.sql"}
	suite.mockFR.On("ReadDirs", basePath).Return(dirs, nil)
	suite.mockFR.On("ReadfilesWithExtension", path.Join(basePath, dirs[0]), ".sql").Return(files, nil)
	suite.mockFR.On("ReadfilesWithExtension", path.Join(basePath, dirs[1]), ".sql").Return(files, nil)
	suite.mockFR.On("ReadFileAsString", path.Join(basePath, dirs[0], files[0])).Return(files[0], nil)
	suite.mockFR.On("ReadFileAsString", path.Join(basePath, dirs[0], files[1])).Return(files[1], nil)
	suite.mockFR.On("ReadFileAsString", path.Join(basePath, dirs[0], files[2])).Return(files[2], nil)
	suite.mockFR.On("ReadFileAsString", path.Join(basePath, dirs[0], files[3])).Return(files[3], nil)
	suite.mockFR.On("ReadFileAsString", path.Join(basePath, dirs[1], files[0])).Return(files[0], nil)
	suite.mockFR.On("ReadFileAsString", path.Join(basePath, dirs[1], files[1])).Return(files[1], nil)
	suite.mockFR.On("ReadFileAsString", path.Join(basePath, dirs[1], files[2])).Return(files[2], nil)
	suite.mockFR.On("ReadFileAsString", path.Join(basePath, dirs[1], files[3])).Return(files[3], nil)
	fs, err := GetFileSource(basePath, suite.mockFR)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), fs)
	schema, _ := fs.GetSchemaList()
	assert.EqualValues(suite.T(), dirs, schema)
	versions, _ := fs.GetSortedVersions("sub1")
	assert.EqualValues(suite.T(), []int{2, 3}, versions)
	upName, up, _ := fs.GetMigrationUpFile("sub1", 3)
	assert.EqualValues(suite.T(), "03-file1.UP.sql", up)
	assert.EqualValues(suite.T(), "03-file1.UP.sql", upName)
	downName, down, _ := fs.GetMigrationDownFile("sub1", 2)
	assert.EqualValues(suite.T(), "02-file2.DOWN.sql", down)
	assert.EqualValues(suite.T(), "02-file2.DOWN.sql", downName)
}

func TestFileSourceTestSuite(t *testing.T) {
	suite.Run(t, new(FileSourceTestSuite))
}
