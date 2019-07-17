package target

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/apty/dbmigrate/target/mocks"
)

type PGDatabaseTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	mockDB  *sql.DB
	db *Database
}
const sampleSchema = "TestSchema"
func (suite *PGDatabaseTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		suite.T().Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	suite.mockDB = mockDB
	suite.sqlMock = sqlMock
	mq := new(mocks.MigrationQueries)
	suite.db = &Database{mq, mockDB}
	mq.On("CountMigrationTableSQL").Return("CountMigrationTableSQL")
	mq.On("CreateMigrationTableSQL", sampleSchema).Return("CreateMigrationTableSQL")
	mq.On("GetMaxSequenceSQL", sampleSchema).Return("GetMaxSequenceSQL")
	mq.On("GetLatestBatchSQL", sampleSchema).Return("GetLatestBatchSQL")
	mq.On("InsertMigrationLogSQL", "test").Return(`INSERT INTO test.db_migrations (sequence,name,batch) values ($1,$2, $3)`)
	mq.On("DeleteMigrationLogSQL", "test").Return(`DELETE FROM test.db_migrations WHERE batch = $1`)
	mq.On("GetSequenceByBatchSQL", sampleSchema).Return("GetSequenceByBatchSQL")
	mq.On("GetLatestBatchSQL", sampleSchema).Return("GetLatestBatchSQL")
	mq.On("SetSchemaSQL", sampleSchema).Return("SetSchemaSQL")
	mq.On("CountSchemaSQL").Return("CountSchemaSQL")
	mq.On("CreateSchemaSQL", sampleSchema).Return("CreateSchemaSQL")
}

func (suite *PGDatabaseTestSuite) TearDownSuite() {
	suite.mockDB.Close()
}


func (suite *PGDatabaseTestSuite) TestDatabase_DoMigrationTableExists_Error() {
	err:=fmt.Errorf("error")
	suite.sqlMock.ExpectQuery("CountMigrationTableSQL").WillReturnError(err)
	_,returnErr := suite.db.DoMigrationTableExists(sampleSchema)
	assert.Error(suite.T(),err,returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_DoMigrationTableExists_True() {
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow("1")
	suite.sqlMock.ExpectQuery("CountMigrationTableSQL").WillReturnRows(rows)
	val,returnErr := suite.db.DoMigrationTableExists(sampleSchema)
	assert.Nil(suite.T(),returnErr)
	assert.True(suite.T(),val)
}

func (suite *PGDatabaseTestSuite) TestDatabase_DoMigrationTableExists_False() {
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow("0")
	suite.sqlMock.ExpectQuery("CountMigrationTableSQL").WillReturnRows(rows)
	val,returnErr := suite.db.DoMigrationTableExists(sampleSchema)
	assert.Nil(suite.T(),returnErr)
	assert.False(suite.T(),val)
}

func (suite *PGDatabaseTestSuite) TestDatabase_CreateMigrationTable_Error() {
	err:=fmt.Errorf("error")
	suite.sqlMock.ExpectExec("CreateMigrationTableSQL").WillReturnError(err)
	returnErr := suite.db.CreateMigrationTable(sampleSchema)
	assert.Error(suite.T(),err,returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_CreateMigrationTable_Success() {
	suite.sqlMock.ExpectExec("CreateMigrationTableSQL").WillReturnResult(sqlmock.NewResult(1, 1))
	returnErr := suite.db.CreateMigrationTable(sampleSchema)
     	assert.Nil(suite.T(),returnErr)
}



func (suite *PGDatabaseTestSuite) TestDatabase_GetMaxSequence_Error() {
	err:=fmt.Errorf("error")
	suite.sqlMock.ExpectQuery("GetMaxSequenceSQL").WillReturnError(err)
	_,returnErr := suite.db.GetMaxSequence(sampleSchema)
	assert.Error(suite.T(),err,returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_GetMaxSequence_Success() {
	rows := sqlmock.NewRows([]string{"sequence"}).
		AddRow("100")
	suite.sqlMock.ExpectQuery("GetMaxSequenceSQL").WillReturnRows(rows)
	val,returnErr := suite.db.GetMaxSequence(sampleSchema)
	assert.Nil(suite.T(),returnErr)
	assert.Equal(suite.T(),100,val)
}


func (suite *PGDatabaseTestSuite) TestDatabase_GetLatestBatch_Error() {
	err:=fmt.Errorf("error")
	suite.sqlMock.ExpectQuery("GetLatestBatchSQL").WillReturnError(err)
	_,returnErr := suite.db.GetLatestBatch(sampleSchema)
	assert.Error(suite.T(),err,returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_GetLatestBatch_Success() {
	rows := sqlmock.NewRows([]string{"batch"}).
		AddRow("batchTest")
	suite.sqlMock.ExpectQuery("GetLatestBatchSQL").WillReturnRows(rows)
	val,returnErr := suite.db.GetLatestBatch(sampleSchema)
	assert.Nil(suite.T(),returnErr)
	assert.Equal(suite.T(),"batchTest",val)
}

func (suite *PGDatabaseTestSuite) TestDatabase_InsertMigrationLog_Error() {
	err:=fmt.Errorf("error")
	suite.sqlMock.ExpectPrepare("INSERT INTO test.db_migrations (.+)").ExpectExec().WillReturnError(err)
	returnErr := suite.db.InsertMigrationLog("test",1,"Test","b1")
	assert.Error(suite.T(),err,returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_InsertMigrationLog_Success() {
	suite.sqlMock.ExpectPrepare("INSERT INTO test.db_migrations (.+)").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	returnErr := suite.db.InsertMigrationLog("test",1,"Test","b1")
	assert.Nil(suite.T(),returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_DeleteMigrationLog_Error() {
	err:=fmt.Errorf("error")
	suite.sqlMock.ExpectPrepare("DELETE FROM test.db_migrations(.+)").ExpectExec().WillReturnError(err)
	returnErr := suite.db.DeleteMigrationLog("test","Test")
	assert.Error(suite.T(),err,returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_DeleteMigrationLog_Success() {
	suite.sqlMock.ExpectPrepare("DELETE FROM test.db_migrations(.+)").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	returnErr := suite.db.DeleteMigrationLog("test","Test")
	assert.Nil(suite.T(),returnErr)
}


func (suite *PGDatabaseTestSuite) TestDatabase_ExecuteMigration_Error_settingSchema() {
	err:=fmt.Errorf("error")
	suite.sqlMock.ExpectExec("SetSchemaSQL").WillReturnError(err)
	returnErr := suite.db.ExecuteMigration(sampleSchema,"Test")
	assert.Error(suite.T(),err,returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_ExecuteMigration_Error_Command() {
	err:=fmt.Errorf("error")
	suite.sqlMock.ExpectExec("SetSchemaSQL").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectExec("Test").WillReturnError(err)
	returnErr := suite.db.ExecuteMigration(sampleSchema,"Test")
	assert.Error(suite.T(),err,returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_ExecuteMigration_Success() {
	suite.sqlMock.ExpectExec("SetSchemaSQL").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectExec("Test").WillReturnResult(sqlmock.NewResult(1, 1))
	returnErr := suite.db.ExecuteMigration(sampleSchema,"Test")
	assert.Nil(suite.T(),returnErr)
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(PGDatabaseTestSuite))
}



