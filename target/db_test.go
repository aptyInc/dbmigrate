package target

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"

	"github.com/apty/dbmigrate/target/mocks"
)

type PGDatabaseTestSuite struct {
	suite.Suite
	sqlMock sqlmock.Sqlmock
	mockDB  *sql.DB
	db      *Database
}

const sampleSchema = "TestSchema"

func (suite *PGDatabaseTestSuite) SetupSuite() {
	mockDB, sqlMock, err := sqlmock.New()
	require.NoError(suite.T(), err)
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
	mq.On("GetSequenceByBatchSQL", "test").Return(`SELECT sequence FROM test.db_migrations WHERE batch = $1`)
	mq.On("GetLatestBatchSQL", sampleSchema).Return("GetLatestBatchSQL")
	mq.On("SetSchemaSQL", sampleSchema).Return("SetSchemaSQL")
	mq.On("CountSchemaSQL").Return("CountSchemaSQL")
	mq.On("CreateSchemaSQL", sampleSchema).Return("CreateSchemaSQL")
}

func (suite *PGDatabaseTestSuite) TearDownSuite() {
	suite.mockDB.Close()
}

func (suite *PGDatabaseTestSuite) AfterTest(_, _ string) {
	require.NoError(suite.T(), suite.sqlMock.ExpectationsWereMet())
}

func (suite *PGDatabaseTestSuite) TestDatabase_DoMigrationTableExists_Error() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectQuery("CountMigrationTableSQL").WillReturnError(err)
	_, returnErr := suite.db.DoMigrationTableExists(sampleSchema)
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_DoMigrationTableExists_True() {
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow("1")
	suite.sqlMock.ExpectQuery("CountMigrationTableSQL").WillReturnRows(rows)
	val, returnErr := suite.db.DoMigrationTableExists(sampleSchema)
	assert.NoError(suite.T(), returnErr)
	assert.True(suite.T(), val)
}

func (suite *PGDatabaseTestSuite) TestDatabase_DoMigrationTableExists_False() {
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow("0")
	suite.sqlMock.ExpectQuery("CountMigrationTableSQL").WillReturnRows(rows)
	val, returnErr := suite.db.DoMigrationTableExists(sampleSchema)
	assert.NoError(suite.T(), returnErr)
	assert.False(suite.T(), val)
}

func (suite *PGDatabaseTestSuite) TestDatabase_CreateMigrationTable_Error() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectExec("CreateMigrationTableSQL").WillReturnError(err)
	returnErr := suite.db.CreateMigrationTable(sampleSchema)
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_CreateMigrationTable_Success() {
	suite.sqlMock.ExpectExec("CreateMigrationTableSQL").WillReturnResult(sqlmock.NewResult(1, 1))
	returnErr := suite.db.CreateMigrationTable(sampleSchema)
	assert.NoError(suite.T(), returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_GetMaxSequence_Error() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectQuery("GetMaxSequenceSQL").WillReturnError(err)
	_, returnErr := suite.db.GetMaxSequence(sampleSchema)
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_GetMaxSequence_Success() {
	rows := sqlmock.NewRows([]string{"sequence"}).
		AddRow("100")
	suite.sqlMock.ExpectQuery("GetMaxSequenceSQL").WillReturnRows(rows)
	val, returnErr := suite.db.GetMaxSequence(sampleSchema)
	assert.NoError(suite.T(), returnErr)
	assert.Equal(suite.T(), 100, val)
}

func (suite *PGDatabaseTestSuite) TestDatabase_GetLatestBatch_Error() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectQuery("GetLatestBatchSQL").WillReturnError(err)
	_, returnErr := suite.db.GetLatestBatch(sampleSchema)
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_GetLatestBatch_Success() {
	rows := sqlmock.NewRows([]string{"batch"}).
		AddRow("batchTest")
	suite.sqlMock.ExpectQuery("GetLatestBatchSQL").WillReturnRows(rows)
	val, returnErr := suite.db.GetLatestBatch(sampleSchema)
	assert.NoError(suite.T(), returnErr)
	assert.Equal(suite.T(), "batchTest", val)
}

func (suite *PGDatabaseTestSuite) TestDatabase_InsertMigrationLog_Error() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectPrepare("INSERT INTO test.db_migrations (.+)").ExpectExec().WillReturnError(err)
	returnErr := suite.db.InsertMigrationLog("test", 1, "Test", "b1")
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_InsertMigrationLog_ErrorWithPrepare() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectPrepare("INSERT INTO test.db_migrations (.+)").WillReturnError(err)
	returnErr := suite.db.InsertMigrationLog("test", 1, "Test", "b1")
	assert.EqualError(suite.T(), returnErr, "custom error")
}
func (suite *PGDatabaseTestSuite) TestDatabase_InsertMigrationLog_Success() {
	suite.sqlMock.ExpectPrepare("INSERT INTO test.db_migrations (.+)").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	returnErr := suite.db.InsertMigrationLog("test", 1, "Test", "b1")
	assert.NoError(suite.T(), returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_DeleteMigrationLog_Error() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectPrepare("DELETE FROM test.db_migrations(.+)").ExpectExec().WillReturnError(err)
	returnErr := suite.db.DeleteMigrationLog("test", "Test")
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_DeleteMigrationLog_ErrorWithStatement() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectPrepare("DELETE FROM test.db_migrations(.+)").WillReturnError(err)
	returnErr := suite.db.DeleteMigrationLog("test", "Test")
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_DeleteMigrationLog_Success() {
	suite.sqlMock.ExpectPrepare("DELETE FROM test.db_migrations(.+)").ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	returnErr := suite.db.DeleteMigrationLog("test", "Test")
	assert.NoError(suite.T(), returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_ExecuteMigration_Error_settingSchema() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectExec("SetSchemaSQL").WillReturnError(err)
	returnErr := suite.db.ExecuteMigration(sampleSchema, "Test")
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_ExecuteMigration_Error_Command() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectExec("SetSchemaSQL").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectExec("Test").WillReturnError(err)
	returnErr := suite.db.ExecuteMigration(sampleSchema, "Test")
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_ExecuteMigration_Success() {
	suite.sqlMock.ExpectExec("SetSchemaSQL").WillReturnResult(sqlmock.NewResult(1, 1))
	suite.sqlMock.ExpectExec("Test").WillReturnResult(sqlmock.NewResult(1, 1))
	returnErr := suite.db.ExecuteMigration(sampleSchema, "Test")
	assert.NoError(suite.T(), returnErr)
}

func (suite *PGDatabaseTestSuite) TestDatabase_GetSequenceByBatch_ErrorWithStatement() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectPrepare("SELECT sequence FROM test.db_migrations WHERE(.+)").WillReturnError(err)
	_, returnErr := suite.db.GetSequenceByBatch("test", "Test")
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_GetSequenceByBatch_ErrorWithQuery() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectPrepare("SELECT sequence FROM test.db_migrations WHERE (.+)").ExpectQuery().WillReturnError(err)
	_, returnErr := suite.db.GetSequenceByBatch("test", "Test")
	assert.EqualError(suite.T(), returnErr, "custom error")
}
func (suite *PGDatabaseTestSuite) TestDatabase_GetSequenceByBatch_SuccessWithOneElement() {
	rows := sqlmock.NewRows([]string{"sequence"}).
		AddRow("100")
	suite.sqlMock.ExpectPrepare("SELECT sequence FROM test.db_migrations WHERE(.+)").ExpectQuery().WillReturnRows(rows)
	array, returnErr := suite.db.GetSequenceByBatch("test", "Test")
	assert.NoError(suite.T(), returnErr)
	assert.ElementsMatch(suite.T(), []int{100}, array)

}
func (suite *PGDatabaseTestSuite) TestDatabase_GetSequenceByBatch_ErrorWithRows() {
	rows := sqlmock.NewRows([]string{"sequence"}).
		AddRow("error")
	suite.sqlMock.ExpectPrepare("SELECT sequence FROM test.db_migrations WHERE(.+)").ExpectQuery().WillReturnRows(rows)
	_, returnErr := suite.db.GetSequenceByBatch("test", "Test")
	assert.EqualError(suite.T(), returnErr, "sql: Scan error on column index 0, name \"sequence\": converting driver.Value type string (\"error\") to a int: invalid syntax")

}

func (suite *PGDatabaseTestSuite) TestDatabase_GetSequenceByBatch_SuccessWithWEmptyElement() {
	rows := sqlmock.NewRows([]string{"sequence"})
	suite.sqlMock.ExpectPrepare("SELECT sequence FROM test.db_migrations WHERE(.+)").ExpectQuery().WillReturnRows(rows)
	array, returnErr := suite.db.GetSequenceByBatch("test", "Test")
	assert.NoError(suite.T(), returnErr)
	assert.ElementsMatch(suite.T(), []int{}, array)

}

func (suite *PGDatabaseTestSuite) TestDatabase_DoesSchemaExists_Error() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectQuery("CountSchemaSQL").WillReturnError(err)
	_, returnErr := suite.db.DoesSchemaExists(sampleSchema)
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_DoesSchemaExists_True() {
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow("1")
	suite.sqlMock.ExpectQuery("CountSchemaSQL").WillReturnRows(rows)
	val, returnErr := suite.db.DoesSchemaExists(sampleSchema)
	assert.NoError(suite.T(), returnErr)
	assert.True(suite.T(), val)
}

func (suite *PGDatabaseTestSuite) TestDatabase_DoesSchemaExists_False() {
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow("0")
	suite.sqlMock.ExpectQuery("CountSchemaSQL").WillReturnRows(rows)
	val, returnErr := suite.db.DoesSchemaExists(sampleSchema)
	assert.NoError(suite.T(), returnErr)
	assert.False(suite.T(), val)
}

func (suite *PGDatabaseTestSuite) TestDatabase_CreateSchema_Error() {
	err := fmt.Errorf("custom error")
	suite.sqlMock.ExpectExec("CreateSchemaSQL").WillReturnError(err)
	returnErr := suite.db.CreateSchema(sampleSchema)
	assert.EqualError(suite.T(), returnErr, "custom error")
}

func (suite *PGDatabaseTestSuite) TestDatabase_CreateSchema_Success() {
	suite.sqlMock.ExpectExec("CreateSchemaSQL").WillReturnResult(sqlmock.NewResult(1, 1))
	returnErr := suite.db.CreateSchema(sampleSchema)
	assert.NoError(suite.T(), returnErr)
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(PGDatabaseTestSuite))
}
