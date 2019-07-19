package source

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type FileReaderTestSuite struct {
	suite.Suite
	fs afero.Fs
	reader FileReader
}

func (suite *FileReaderTestSuite) SetupSuite() {
	suite.fs = afero.NewMemMapFs()
	suite.reader = &ReaderImplementation{suite.fs}
}

func (suite *FileReaderTestSuite) SetupTest() {
	suite.fs.RemoveAll("/base")
	suite.fs.RemoveAll("/base1")
}

func (suite *FileReaderTestSuite) TestFileReader_ReadDirs_dirs(){
	suite.fs.MkdirAll("/base", 0755)
	suite.fs.MkdirAll("/base/cmd1", 0755)
	suite.fs.MkdirAll("/base/cmd1/cmd3", 0755)
	suite.fs.MkdirAll("/base/cmd2", 0755)
	arr, err := suite.reader.ReadDirs("/base")
	assert.NoError(suite.T(),err)
	assert.ElementsMatch(suite.T(),[]string{"cmd1","cmd2"},arr)
}

func (suite *FileReaderTestSuite) TestFileReader_ReadDirs_empty(){
	suite.fs.MkdirAll("/base", 0755)
	arr, err := suite.reader.ReadDirs("/base")
	assert.NoError(suite.T(),err)
	assert.ElementsMatch(suite.T(),[]string{} , arr)
}

func (suite *FileReaderTestSuite) TestFileReader_ReadDirs_error() {
	_, err := suite.reader.ReadDirs("/base")
	assert.EqualError(suite.T(), err, "open /base: file does not exist")
}

func (suite *FileReaderTestSuite) TestFileReader_ReadFilesWithExtension_success() {
	suite.fs.MkdirAll("/base", 0755)
	suite.fs.MkdirAll("/base1", 0755)
	afero.WriteFile( suite.fs,"/base1/a.go", []byte("file a"), 0644)
	afero.WriteFile( suite.fs,"/base/a.go", []byte("file a"), 0644)
	afero.WriteFile( suite.fs,"/base/b.go", []byte("file b"), 0644)
	afero.WriteFile( suite.fs,"/base/c.go1", []byte("file c"), 0644)
	afero.WriteFile( suite.fs,"/base/go.b", []byte("file go"), 0644)

	arr, err := suite.reader.ReadFilesWithExtension("/base","*.go")
	assert.NoError(suite.T(),err)
	assert.ElementsMatch(suite.T(),[]string{"/base/a.go","/base/b.go"} , arr)

}

func (suite *FileReaderTestSuite) TestFileReader_ReadFileAsString_error() {
	_, err := suite.reader.ReadFileAsString("/base")
	assert.EqualError(suite.T(), err, "open /base: file does not exist")
}

func (suite *FileReaderTestSuite) TestFileReader_ReadFileAsString_success() {
	suite.fs.MkdirAll("/base", 0755)
	afero.WriteFile( suite.fs,"/base/a.go", []byte("file a"), 0644)
	arr, err := suite.reader.ReadFileAsString("/base/a.go")
	assert.NoError(suite.T(),err)
	assert.Equal(suite.T(),"file a" , arr)

}
func TestFileReaderTestSuite(t *testing.T) {
	suite.Run(t, new(FileReaderTestSuite))
}