package target

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Database interface{
	DoMigrationTableExists(schema string) (bool, error)
	CreateMigrationTable(schema string) error
	GetMaxSequence(schema string) (int, error)
	GetLatestBatch(schema string) (string, error)
	InsertMigrationLog(schema string, version int, name string, batch string) error
	DeleteMigrationLog(schema string, batch string) error
	ExecuteMigration(schema string, command string) error
	GetSequenceByBatch(schema string, batch string) ([]int, error)
	DoesSchemaExists(schema string) (bool, error)
	CreateSchema(schema string) error
}

//DatabaseImplementation to access the database
type DatabaseImplementation struct {
	mq MigrationQueries
	DB *sql.DB
}

//DoMigrationTableExists Check if migration table exsists
func (base *DatabaseImplementation) DoMigrationTableExists(schema string) (bool, error) {
	var count int
	err := base.DB.QueryRow(base.mq.CountMigrationTableSQL(), schema).Scan(&count)
	return count > 0, err
}

//CreateMigrationTable Create a migration table
func (base *DatabaseImplementation) CreateMigrationTable(schema string) error {
	_, err := base.DB.Exec(base.mq.CreateMigrationTableSQL(schema))
	return err
}

//GetMaxSequence Gets the max sequence in the schema
func (base *DatabaseImplementation) GetMaxSequence(schema string) (int, error) {
	var num int
	err := base.DB.QueryRow(base.mq.GetMaxSequenceSQL(schema)).Scan(&num)
	return num, err
}

//GetLatestBatch Gets the latest batch in the schema
func (base *DatabaseImplementation) GetLatestBatch(schema string) (string, error) {
	var batch string
	row := base.DB.QueryRow(base.mq.GetLatestBatchSQL(schema))

	switch err := row.Scan( &batch); err {
	case sql.ErrNoRows:
		fmt.Println("No batch were returned!")
		return "",nil
	case nil:
		return batch, err
	default:
		return "", err
	}
	//.Scan(&batch)

}

//InsertMigrationLog inserts a migration log
func (base *DatabaseImplementation) InsertMigrationLog(schema string, version int, name string, batch string) error {
	stmt, err := base.DB.Prepare(base.mq.InsertMigrationLogSQL(schema))

	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err1 := stmt.Exec(version, name, batch)
	return err1
}

//DeleteMigrationLog Deletes a batch of migration logs
func (base *DatabaseImplementation) DeleteMigrationLog(schema string, batch string) error {
	stmt, err := base.DB.Prepare(base.mq.DeleteMigrationLogSQL(schema))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err1 := stmt.Exec(batch)
	return err1
}

//ExecuteMigration executes the given SQL as script
func (base *DatabaseImplementation) ExecuteMigration(schema string, command string) error {
	_, err1 := base.DB.Exec(base.mq.SetSchemaSQL(schema))
	if err1 != nil {
		return err1
	}
	_, err2 := base.DB.Exec(command)
	return err2

}

//GetSequenceByBatch get the sequence ids by batch
func (base *DatabaseImplementation) GetSequenceByBatch(schema string, batch string) ([]int, error) {
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
func (base *DatabaseImplementation) DoesSchemaExists(schema string) (bool, error) {
	var count int
	err := base.DB.QueryRow(base.mq.CountSchemaSQL(), schema).Scan(&count)
	return count > 0, err
}

//CreateSchema Create the schema
func (base *DatabaseImplementation) CreateSchema(schema string) error {
	_, err := base.DB.Exec(base.mq.CreateSchemaSQL(schema))
	return err
}
