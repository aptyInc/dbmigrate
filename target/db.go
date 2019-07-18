package target

import (
	"database/sql"
)

//Database to access the database
type Database struct {
	mq MigrationQueries
	DB *sql.DB
}

//DoMigrationTableExists Check if migration table exsists
func (base *Database) DoMigrationTableExists(schema string) (bool, error) {
	var count int
	err := base.DB.QueryRow(base.mq.CountMigrationTableSQL(), schema).Scan(&count)
	return count > 0, err
}

//CreateMigrationTable Create a migration table
func (base *Database) CreateMigrationTable(schema string) error {
	_, err := base.DB.Exec(base.mq.CreateMigrationTableSQL(schema))
	return err
}

//GetMaxSequence Gets the max sequence in the schema
func (base *Database) GetMaxSequence(schema string) (int, error) {
	var num int
	err := base.DB.QueryRow(base.mq.GetMaxSequenceSQL(schema)).Scan(&num)
	return num, err
}

//GetLatestBatch Gets the latest batch in the schema
func (base *Database) GetLatestBatch(schema string) (string, error) {
	var batch string
	err := base.DB.QueryRow(base.mq.GetLatestBatchSQL(schema)).Scan(&batch)
	return batch, err
}

//InsertMigrationLog inserts a migration log
func (base *Database) InsertMigrationLog(schema string, version int, name string, batch string) error {
	stmt, err := base.DB.Prepare(base.mq.InsertMigrationLogSQL(schema))

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err1 := stmt.Exec(version, name, batch)
	return err1
}

//DeleteMigrationLog Deletes a batch of migration logs
func (base *Database) DeleteMigrationLog(schema string, batch string) error {
	stmt, err := base.DB.Prepare(base.mq.DeleteMigrationLogSQL(schema))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err1 := stmt.Exec(batch)
	return err1
}

//ExecuteMigration executes the given SQL as script
func (base *Database) ExecuteMigration(schema string, command string) error {
	_, err1 := base.DB.Exec(base.mq.SetSchemaSQL(schema))
	if err1 != nil {
		return err1
	}
	_, err2 := base.DB.Exec(command)
	return err2

}

//GetSequenceByBatch get the sequence ids by batch
func (base *Database) GetSequenceByBatch(schema string, batch string) ([]int, error) {
	var sequences []int
	stmt, err := base.DB.Prepare(base.mq.GetSequenceByBatchSQL(schema))
	if err != nil {
		return sequences, err
	}
	defer stmt.Close()
	result, err1 := stmt.Query(batch)
	if err1 != nil {
		return sequences, err1
	}
	defer result.Close()
	var sequence int
	for result.Next() {
		err := result.Scan(&sequence)
		if err != nil {
			return sequences, err
		}
		sequences = append(sequences, sequence)
	}
	err = result.Err()
	return sequences, err
}

//DoesSchemaExists Check if schema exists
func (base *Database) DoesSchemaExists(schema string) (bool, error) {
	var count int
	err := base.DB.QueryRow(base.mq.CountSchemaSQL(), schema).Scan(&count)
	return count > 0, err
}

//CreateSchema Create the schema
func (base *Database) CreateSchema(schema string) error {
	_, err := base.DB.Exec(base.mq.CreateSchemaSQL(schema))
	return err
}
