package target

import "fmt"

//Postgres struct to return Postgres queries
type Postgres struct{}

//CountMigrationTableSQL query to check if migration table exists
func (p Postgres) CountMigrationTableSQL() string {
	return "SELECT count(1) FROM pg_catalog.pg_tables WHERE tablename = 'db_migrations' AND schemaname = $1"
}

//CreateMigrationTableSQL query to create migration table
func (p Postgres) CreateMigrationTableSQL(schema string) string {
	return fmt.Sprintf(`CREATE TABLE "%s".db_migrations (
									id           SERIAL       PRIMARY KEY,
									sequence		 INT UNIQUE NOT NULL,
									name				 VARCHAR(255)  UNIQUE NOT NULL,
									batch				 VARCHAR(150)  NOT NULL,
                  run_on 			 TIMESTAMP WITHOUT TIME ZONE DEFAULT TIMEZONE('utc',now()) NOT NULL
                )`, schema)
}

//GetMaxSequenceSQL query to return Max sequence in the schema
func (p Postgres) GetMaxSequenceSQL(schema string) string {
	return fmt.Sprintf(`SELECT coalesce(max(sequence),-1) FROM  "%s".db_migrations`, schema)
}

//GetLatestBatchSQL query to return latest batch name in the schema
func (p Postgres) GetLatestBatchSQL(schema string) string {
	return fmt.Sprintf(`SELECT batch FROM "%s".db_migrations ORDER  BY run_on DESC LIMIT 1`, schema)
}

//InsertMigrationLogSQL query to insert a migration
func (p Postgres) InsertMigrationLogSQL(schema string) string {
	return fmt.Sprintf(`INSERT INTO "%s".db_migrations (sequence,name,batch) values ($1,$2, $3)`, schema)
}

//DeleteMigrationLogSQL query to delete migration
func (p Postgres) DeleteMigrationLogSQL(schema string) string {
	return fmt.Sprintf(`DELETE FROM "%s".db_migrations WHERE batch = $1`, schema)
}

//SetSchemaSQL sets the schema to correct folder
func (p Postgres) SetSchemaSQL(schema string) string {
	return fmt.Sprintf(`SET search_path TO %s`, schema)
}

//GetSequenceByBatchSQL query to return sequences by batch
func (p Postgres) GetSequenceByBatchSQL(schema string) string {
	return fmt.Sprintf(`SELECT sequence FROM "%s".db_migrations WHERE batch = $1`, schema)
}

//GetSequenceByBatchSQL query to return sequences by batch
func (p Postgres) CountSchemaSQL() string {
	return `SELECT count(1) FROM information_schema.schemata WHERE schema_name  = $1`
}

//CreateSchemaSQL query to check if migration schema exists
func (p Postgres) CreateSchemaSQL(schema string) string {
	return fmt.Sprintf(`CREATE SCHEMA "%s"`, schema)
}
