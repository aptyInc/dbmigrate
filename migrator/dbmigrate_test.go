package migrator

import (
	"fmt"
	sourceMocks "github.com/apty/dbmigrate/source/mocks"
	targetMocks "github.com/apty/dbmigrate/target/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type DBMigrationTestSuite struct {
	suite.Suite
	src *sourceMocks.MigrationSource
	tgt *targetMocks.Database
	migration DBMigration
}

const customError  ="MyError"

func (suite *DBMigrationTestSuite) SetupTest() {
	suite.src = new(sourceMocks.MigrationSource)
	suite.tgt = new(targetMocks.Database)
	suite.migration =  &DBMigrationImplementation{
		suite.src,
		suite.tgt,
	}
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaUp_error1() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test1").Return(false,err)
	err1 := suite.migration.MigrateSchemaUp("Test1")
	assert.EqualError(suite.T(),err1,customError,"Extend Error")
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaUp_error2() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test").Return(false,nil)
	suite.tgt.On("CreateMigrationTable","Test").Return(err)
	err1 := suite.migration.MigrateSchemaUp("Test")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaUp_error3() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test").Return(true,nil)
	suite.tgt.On("GetMaxSequence","Test").Return(-1,err)
	err1 := suite.migration.MigrateSchemaUp("Test")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaUp_error4() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test").Return(true,nil)
	suite.tgt.On("GetMaxSequence","Test").Return(-1,nil)
	suite.src.On("GetSortedVersions","Test").Return(nil,err)

	err1 := suite.migration.MigrateSchemaUp("Test")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaUp_error5() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test").Return(true,nil)
	suite.tgt.On("GetMaxSequence","Test").Return(-1,nil)
	suite.src.On("GetSortedVersions","Test").Return([]int{1,2},nil)
	suite.src.On("GetMigrationUpFile","Test",1).Return("","",err)
	err1 := suite.migration.MigrateSchemaUp("Test")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaUp_error6() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test").Return(true,nil)
	suite.tgt.On("GetMaxSequence","Test").Return(-1,nil)
	suite.src.On("GetSortedVersions","Test").Return([]int{1,2},nil)
	suite.src.On("GetMigrationUpFile","Test",1).Return("1.file.UP.sql","command",nil)
	suite.tgt.On("ExecuteMigration","Test","command").Return(err)
	err1 := suite.migration.MigrateSchemaUp("Test")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaUp_error7() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test").Return(true,nil)
	suite.tgt.On("GetMaxSequence","Test").Return(-1,nil)
	suite.src.On("GetSortedVersions","Test").Return([]int{1},err)
	suite.src.On("GetMigrationUpFile",1,"b1").Return("1.file.UP.sql","command",nil)
	suite.tgt.On("ExecuteMigration","Test","command").Return(nil)
	suite.tgt.On("InsertMigrationLog","Test",1,"1.file.UP.sql","b1").Return(err)
	err1 := suite.migration.MigrateSchemaUp("Test")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaUp_ForNew() {
	suite.tgt.On("DoMigrationTableExists","Test").Return(true,nil)
	suite.tgt.On("GetMaxSequence","Test").Return(-1,nil)
	suite.src.On("GetSortedVersions","Test").Return([]int{1},nil)
	suite.src.On("GetMigrationUpFile","Test",1).Return("1.file.UP.sql","command",nil)
	suite.tgt.On("ExecuteMigration","Test","command").Return(nil)
	suite.tgt.On("InsertMigrationLog","Test",1,"1.file.UP.sql",mock.AnythingOfType("string")).Return(nil)
	err1 := suite.migration.MigrateSchemaUp("Test")
	assert.NoError(suite.T(),err1)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaUp_ForNoUpdates() {
	suite.tgt.On("DoMigrationTableExists","Test").Return(true,nil)
	suite.tgt.On("GetMaxSequence","Test").Return(5,nil)
	suite.src.On("GetSortedVersions","Test").Return([]int{1,2,3,4},nil)
	err1 := suite.migration.MigrateSchemaUp("Test")
	assert.NoError(suite.T(),err1)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateUp_Error1() {
	err := fmt.Errorf(customError)
	suite.src.On("GetSchemaList").Return(nil,err)
	err1 := suite.migration.MigrateUp()
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateUp_Error2() {
	err := fmt.Errorf(customError)
	suite.src.On("GetSchemaList").Return([]string{"Test"},nil)
	suite.tgt.On("DoMigrationTableExists","Test").Return(false,err)
	err1 := suite.migration.MigrateUp()
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateUp_ForNew() {
	suite.src.On("GetSchemaList").Return([]string{"Test"},nil)
	suite.tgt.On("DoMigrationTableExists","Test").Return(true,nil)
	suite.tgt.On("GetMaxSequence","Test").Return(-1,nil)
	suite.src.On("GetSortedVersions","Test").Return([]int{1},nil)
	suite.src.On("GetMigrationUpFile","Test",1).Return("1.file.UP.sql","command",nil)
	suite.tgt.On("ExecuteMigration","Test","command").Return(nil)
	suite.tgt.On("InsertMigrationLog","Test",1,"1.file.UP.sql",mock.AnythingOfType("string")).Return(nil)
	err1 := suite.migration.MigrateUp()
	assert.NoError(suite.T(),err1)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateUp_ForNoUpdates() {
	suite.src.On("GetSchemaList").Return([]string{"Test"},nil)
	suite.tgt.On("DoMigrationTableExists","Test").Return(true,nil)
	suite.tgt.On("GetMaxSequence","Test").Return(5,nil)
	suite.src.On("GetSortedVersions","Test").Return([]int{1,2,3,4},nil)
	err1 := suite.migration.MigrateUp()
	assert.NoError(suite.T(),err1)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaDown_error1() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test1").Return(false,err)
	err1 := suite.migration.MigrateSchemaDown("Test1")
	assert.EqualError(suite.T(),err1,customError,"Extend Error")
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaDown_MigrationDoesntExist() {
	suite.tgt.On("DoMigrationTableExists","Test1").Return(false,nil)
	err1 := suite.migration.MigrateSchemaDown("Test1")
	assert.EqualError(suite.T(),err1,"no Migrations in schema")
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaDown_error2() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test1").Return(true,nil)
	suite.tgt.On("GetLatestBatch","Test1").Return("",err)
	err1 := suite.migration.MigrateSchemaDown("Test1")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaDown_error3() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test1").Return(true,nil)
	suite.tgt.On("GetLatestBatch","Test1").Return("batch1",nil)
	suite.tgt.On("GetSequenceByBatch","Test1","batch1").Return(nil,err)
	err1 := suite.migration.MigrateSchemaDown("Test1")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaDown_error4() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test1").Return(true,nil)
	suite.tgt.On("GetLatestBatch","Test1").Return("batch1",nil)
	suite.tgt.On("GetSequenceByBatch","Test1","batch1").Return([]int{1},nil)
	suite.src.On("GetMigrationDownFile","Test1",1).Return("","",err)

	err1 := suite.migration.MigrateSchemaDown("Test1")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaDown_error5() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test1").Return(true,nil)
	suite.tgt.On("GetLatestBatch","Test1").Return("batch1",nil)
	suite.tgt.On("GetSequenceByBatch","Test1","batch1").Return([]int{1},nil)
	suite.src.On("GetMigrationDownFile","Test1",1).Return("","command",nil)
	suite.tgt.On("ExecuteMigration","Test1","command").Return(err)
	err1 := suite.migration.MigrateSchemaDown("Test1")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaDown_error6() {
	err := fmt.Errorf(customError)
	suite.tgt.On("DoMigrationTableExists","Test1").Return(true,nil)
	suite.tgt.On("GetLatestBatch","Test1").Return("batch1",nil)
	suite.tgt.On("GetSequenceByBatch","Test1","batch1").Return([]int{1},nil)
	suite.src.On("GetMigrationDownFile","Test1",1).Return("","command",nil)
	suite.tgt.On("ExecuteMigration","Test1","command").Return(nil)
	suite.tgt.On("DeleteMigrationLog","Test1","batch1").Return(err)

	err1 := suite.migration.MigrateSchemaDown("Test1")
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateSchemaDown_success() {
	suite.tgt.On("DoMigrationTableExists","Test1").Return(true,nil)
	suite.tgt.On("GetLatestBatch","Test1").Return("batch1",nil)
	suite.tgt.On("GetSequenceByBatch","Test1","batch1").Return([]int{1},nil)
	suite.src.On("GetMigrationDownFile","Test1",1).Return("","command",nil)
	suite.tgt.On("ExecuteMigration","Test1","command").Return(nil)
	suite.tgt.On("DeleteMigrationLog","Test1","batch1").Return(nil)

	err1 := suite.migration.MigrateSchemaDown("Test1")
	assert.NoError(suite.T(),err1)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateDown_success() {
	suite.src.On("GetSchemaList").Return([]string{"Test1"},nil)
	suite.tgt.On("DoMigrationTableExists","Test1").Return(true,nil)
	suite.tgt.On("GetLatestBatch","Test1").Return("batch1",nil)
	suite.tgt.On("GetSequenceByBatch","Test1","batch1").Return([]int{1},nil)
	suite.src.On("GetMigrationDownFile","Test1",1).Return("","command",nil)
	suite.tgt.On("ExecuteMigration","Test1","command").Return(nil)
	suite.tgt.On("DeleteMigrationLog","Test1","batch1").Return(nil)

	err1 := suite.migration.MigrateDown()
	assert.NoError(suite.T(),err1)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateDown_ErrorWithMigration() {
	err := fmt.Errorf(customError)
	suite.src.On("GetSchemaList").Return([]string{"Test1"},nil)
	suite.tgt.On("DoMigrationTableExists","Test1").Return(false,err)

	err1 := suite.migration.MigrateDown()
	assert.EqualError(suite.T(),err1,customError)
}

func (suite *DBMigrationTestSuite) TestDBMigrationTest_MigrateDown_ErrorWithSchemaList() {
	err := fmt.Errorf(customError)
	suite.src.On("GetSchemaList").Return(nil,err)
	err1 := suite.migration.MigrateDown()
	assert.EqualError(suite.T(),err1,customError)
}

func TestDBMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(DBMigrationTestSuite))
}