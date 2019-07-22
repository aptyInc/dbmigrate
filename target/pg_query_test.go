package target

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PgQueryTestSuite struct {
	suite.Suite
	pg *Postgres
}

//TODO Check the Arguments in each string
func (suite *PgQueryTestSuite) SetupSuite() {
	suite.pg = &Postgres{}
}

func (suite *PgQueryTestSuite) TearDownSuite() {
}

func (suite *PgQueryTestSuite) TestPostgres_CountMigrationTableSQL() {
	returnValue := suite.pg.CountMigrationTableSQL()
	assert.Equal(suite.T(), "SELECT count(1) FROM pg_catalog.pg_tables WHERE tablename = 'db_migrations' AND schemaname = $1", returnValue)
}

func (suite *PgQueryTestSuite) TestPostgres_CreateMigrationTableSQL() {
	returnValue := suite.pg.CreateMigrationTableSQL("test")
	assert.Equal(suite.T(), `CREATE TABLE "test".db_migrations (
									id           SERIAL       PRIMARY KEY,
									sequence		 INT UNIQUE NOT NULL,
									name				 VARCHAR(255)  UNIQUE NOT NULL,
									batch				 VARCHAR(150)  NOT NULL,
                  run_on 			 TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc',now()) NOT NULL
                )`, returnValue)
}

func (suite *PgQueryTestSuite) TestPostgres_GetMaxSequenceSQL() {
	returnValue := suite.pg.GetMaxSequenceSQL("test")
	assert.Equal(suite.T(), `SELECT coalesce(max(sequence),-1) FROM  "test".db_migrations`, returnValue)
}

func (suite *PgQueryTestSuite) TestPostgres_GetLatestBatchSQL() {
	returnValue := suite.pg.GetLatestBatchSQL("test")
	assert.Equal(suite.T(), `SELECT batch FROM "test".db_migrations ORDER  BY run_on DESC LIMIT 1`, returnValue)
}

func (suite *PgQueryTestSuite) TestPostgres_InsertMigrationLogSQL() {
	returnValue := suite.pg.InsertMigrationLogSQL("test")
	assert.Equal(suite.T(), `INSERT INTO "test".db_migrations (sequence,name,batch) values ($1,$2, $3)`, returnValue)
}

func (suite *PgQueryTestSuite) TestPostgres_DeleteMigrationLogSQL() {
	returnValue := suite.pg.DeleteMigrationLogSQL("test")
	assert.Equal(suite.T(), `DELETE FROM "test".db_migrations WHERE batch = $1`, returnValue)
}

func (suite *PgQueryTestSuite) TestPostgres_SetSchemaSQL() {
	returnValue := suite.pg.SetSchemaSQL("test")
	assert.Equal(suite.T(), `SET search_path TO test`, returnValue)
}

func (suite *PgQueryTestSuite) TestPostgres_GetSequenceByBatchSQL() {
	returnValue := suite.pg.GetSequenceByBatchSQL("test")
	assert.Equal(suite.T(), `SELECT sequence FROM "test".db_migrations WHERE batch = $1`, returnValue)
}

func (suite *PgQueryTestSuite) TestPostgres_CountSchemaSQL() {
	returnValue := suite.pg.CountSchemaSQL()
	assert.Equal(suite.T(), `SELECT count(1) FROM information_schema.schemata WHERE schema_name  = $1`, returnValue)
}
func (suite *PgQueryTestSuite) TestPostgres_CreateSchemaSQL() {
	returnValue := suite.pg.CreateSchemaSQL("test")
	assert.Equal(suite.T(), `CREATE SCHEMA "test"`, returnValue)
}

func TestPgQueryTestSuite(t *testing.T) {
	suite.Run(t, new(PgQueryTestSuite))
}
