package source

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"path/filepath"
	"testing"
)

type FileReaderTestSuite struct {
	suite.Suite
	fs afero.Fs
	reader FileReader
	basePath string
}

func (suite *FileReaderTestSuite) SetupSuite() {
	suite.fs = afero.NewMemMapFs()
	suite.reader = &ReaderImplementation{suite.fs}
	suite.basePath =  filepath.Join("base")
}

func (suite *FileReaderTestSuite) SetupTest() {
	suite.fs.RemoveAll(suite.basePath)
	suite.fs.RemoveAll(filepath.Join("baseq"))
}

func (suite *FileReaderTestSuite) TestFileReader_ReadDirs_dirs(){
	suite.fs.MkdirAll(suite.basePath, 0755)
	suite.fs.MkdirAll(filepath.Join("base","cmd1"), 0755)
	suite.fs.MkdirAll(filepath.Join("base","cmd1","cmd3"), 0755)
	suite.fs.MkdirAll(filepath.Join("base","cmd2"), 0755)
	arr, err := suite.reader.ReadDirs(suite.basePath)
	assert.NoError(suite.T(),err)
	assert.ElementsMatch(suite.T(),[]string{"cmd1","cmd2"},arr)
}

func (suite *FileReaderTestSuite) TestFileReader_ReadDirs_empty(){
	suite.fs.MkdirAll(suite.basePath, 0755)
	arr, err := suite.reader.ReadDirs(suite.basePath)
	assert.NoError(suite.T(),err)
	assert.ElementsMatch(suite.T(),[]string{} , arr)
}

func (suite *FileReaderTestSuite) TestFileReader_ReadDirs_error() {
	_, err := suite.reader.ReadDirs(suite.basePath)
	assert.EqualErrorf(suite.T(), err, "open base: file does not exist",suite.basePath)
}

func (suite *FileReaderTestSuite) TestFileReader_ReadFilesWithExtension_success() {
	suite.fs.MkdirAll(suite.basePath, 0755)
	suite.fs.MkdirAll(filepath.Join("base1"), 0755)
	afero.WriteFile( suite.fs,filepath.Join("base1","a.go"), []byte("file a"), 0644)
	afero.WriteFile( suite.fs,filepath.Join("base","a.go"), []byte("file a"), 0644)
	afero.WriteFile( suite.fs,filepath.Join("base","b.go"), []byte("file b"), 0644)
	afero.WriteFile( suite.fs,filepath.Join("base","c.go1"), []byte("file c"), 0644)
	afero.WriteFile( suite.fs,filepath.Join("base","go.b"), []byte("file go"), 0644)

	arr, err := suite.reader.ReadFilesWithExtension(suite.basePath,"*.go")
	assert.NoError(suite.T(),err)
	assert.ElementsMatch(suite.T(),arr,[]string{filepath.Join("base","a.go"),filepath.Join("base","b.go")} )

}

func (suite *FileReaderTestSuite) TestFileReader_ReadFileAsString_error() {
	_, err := suite.reader.ReadFileAsString(suite.basePath)
	assert.EqualErrorf(suite.T(), err, "open base: file does not exist",suite.basePath)
}

func (suite *FileReaderTestSuite) TestFileReader_ReadFileAsString_success() {
	suite.fs.MkdirAll(suite.basePath, 0755)
	afero.WriteFile( suite.fs,filepath.Join("base","a.go"), []byte("file a"), 0644)
	arr, err := suite.reader.ReadFileAsString(filepath.Join("base","a.go"))
	assert.NoError(suite.T(),err)
	assert.Equal(suite.T(),"file a" , arr)

}
func TestFileReaderTestSuite(t *testing.T) {
	suite.Run(t, new(FileReaderTestSuite))
}