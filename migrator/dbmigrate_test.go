package migrator

import (
	sourceMocks "github.com/apty/dbmigrate/source/mocks"
	"github.com/apty/dbmigrate/target"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DBMigrationTestSuite struct {
	suite.Suite
	src *sourceMocks.MigrationSource
	tgt *target.DatabaseImplementation
	migration *DBMigrationImplementation
}

func (suite *DBMigrationTestSuite) SetupSuite() {
	suite.src = new(sourceMocks.MigrationSource)
	//suite.reader = new(targetMocks.MigrationQueries)
	//suite.basePath =  filepath.Join("base")
}

func (suite *DBMigrationTestSuite) TestBMigrationTest_ReadFileAsString_error() {
	assert.Equal(suite.T(),1,1)
}

func TestDBMigrationTestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(DBMigrationTestSuite))
}