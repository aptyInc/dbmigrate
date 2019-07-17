package source

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// FileSourceInternalTestSuite Suite to test
type FileSourceInternalTestSuite struct {
	suite.Suite
}

// TestUniqueForNil check with nil
func (suite *FileSourceInternalTestSuite) TestUniqueForNil() {
	var input []int
	resut := unique(input)
	assert.Len(suite.T(), resut, 0)
}

//TestUniqueForEmptyArray check with empty array
func (suite *FileSourceInternalTestSuite) TestUniqueForEmptyArray() {
	input := []int{}
	resut := unique(input)
	assert.Len(suite.T(), resut, 0)
}

//TestUniqueForDuplicateValues check if duplicates are removed
func (suite *FileSourceInternalTestSuite) TestUniqueForDuplicateValues() {
	input := []int{1, 1}
	resut := unique(input)
	assert.Len(suite.T(), resut, 1)
}

//TestUniqueForUniqueValues checks if unique values are correct
func (suite *FileSourceInternalTestSuite) TestUniqueForUniqueValues() {
	input := []int{1, 2}
	resut := unique(input)
	assert.Len(suite.T(), resut, 2)
}

//TestUniqueForCombinedValues checks if the values are correct
func (suite *FileSourceInternalTestSuite) TestUniqueForCombinedValues() {
	input := []int{1, 2, 2}
	resut := unique(input)
	assert.Len(suite.T(), resut, 2)
}

//TestGetFileDetailsWrongInputs tests invalid inputs
func (suite *FileSourceInternalTestSuite) TestGetFileDetailsWrongInputs() {
	inputs := []string{
		"abcd.json",
		"1-test.sql",
		"test.UP.sql",
		"1.UP.SQL",
		"3.4-NEW.UP.SQL",
	}
	for _, input := range inputs {
		_, _, _, err := getFileDetails(input)
		assert.Error(suite.T(), err)
	}
}

//TestGetFileDetailsValidUP to test if up file is retured as expected
func (suite *FileSourceInternalTestSuite) TestGetFileDetailsValidUP() {

	version, isUp, isDown, err := getFileDetails("15-Testing.UP.sql")
	assert.Equal(suite.T(), 15, version)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), isUp)
	assert.False(suite.T(), isDown)

}

//TestGetFileDetailsInValidUPOrDown to test a file when it is not UP or down
func (suite *FileSourceInternalTestSuite) TestGetFileDetailsInValidUPOrDown() {

	version, isUp, isDown, err := getFileDetails("15-Testing.u.sql")
	assert.Equal(suite.T(), 15, version)
	assert.Nil(suite.T(), err)
	assert.False(suite.T(), isUp)
	assert.False(suite.T(), isDown)

}

//TestGetFileDetailsValidDown to test if Down file is retured as expected
func (suite *FileSourceInternalTestSuite) TestGetFileDetailsValidDown() {

	version, isUp, isDown, err := getFileDetails("15-Testing.DOWN.sql")
	assert.Equal(suite.T(), 15, version)
	assert.Nil(suite.T(), err)
	assert.True(suite.T(), isDown)
	assert.False(suite.T(), isUp)

}

//TestFileSourceInternalTestSuite  tests the internal functions
func TestFileSourceInternalTestSuite(t *testing.T) {
	suite.Run(t, new(FileSourceInternalTestSuite))
}
