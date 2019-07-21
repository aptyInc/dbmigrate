package migrator

import (
	"fmt"
	"time"

	"github.com/apty/dbmigrate/source"
	"github.com/apty/dbmigrate/target"
)

type DBMigration interface{
	MigrateSchemaUp(schema string) error
	MigrateSchemaDown(schema string) error
	MigrateUp() error
	MigrateDown() error
}
//DBMigrationImplementation the main service to execute the migration
type DBMigrationImplementation struct {
	Src source.MigrationSource
	Tgt target.Database
}

func (dbm *DBMigrationImplementation) initGetMaxSequence(schema string) (int, error) {
	isFound, err1 := dbm.Tgt.DoMigrationTableExists(schema)
	if err1 != nil {
		return -1, err1
	}
	if !isFound {
		err1 := dbm.Tgt.CreateMigrationTable(schema)
		return -1, err1
	}
	return dbm.Tgt.GetMaxSequence(schema)
}
func (dbm *DBMigrationImplementation) migrateSequenceUp(schema string, version int, batch string) error {
	name, commands, err1 := dbm.Src.GetMigrationUpFile(schema, version)

	if err1 != nil {
		return err1
	}
	err2 := dbm.Tgt.ExecuteMigration(schema, commands)
	if err2 != nil {
		return err2
	}

	return dbm.Tgt.InsertMigrationLog(schema, version, name, batch)

}

func (dbm *DBMigrationImplementation) migrateSequenceDown(schema string, sequence int) error {
	_, commands, err1 := dbm.Src.GetMigrationDownFile(schema, sequence)

	if err1 != nil {
		return err1
	}
	return dbm.Tgt.ExecuteMigration(schema, commands)

}

func (dbm *DBMigrationImplementation) MigrateSchemaUp(schema string) error {
	doesSchemaExists, err := dbm.Tgt.DoesSchemaExists(schema)
	if err != nil {
		return err
	}
	if !doesSchemaExists {
		dbm.Tgt.CreateSchema(schema)
	}
	maxSequence, err1 := dbm.initGetMaxSequence(schema)
	if err1 != nil {
		return err1
	}
	allVersions, err2 := dbm.Src.GetSortedVersions(schema)
	if err2 != nil {
		return err2
	}
	batch := time.Now().String()
	for _, version := range allVersions {
		if version > maxSequence {
			err3 := dbm.migrateSequenceUp(schema, version, batch)
			if err3 != nil {
				return err3
			}
		}
	}
	return nil
}

func (dbm *DBMigrationImplementation) MigrateSchemaDown(schema string) error {

	doExist, err1 := dbm.Tgt.DoMigrationTableExists(schema)
	if err1 != nil {
		return err1
	}
	if !doExist {
		return fmt.Errorf("no Migrations in schema")
	}

	batch, err2 := dbm.Tgt.GetLatestBatch(schema)
	if err2 != nil {
		return err2
	}
	sequences, err3 := dbm.Tgt.GetSequenceByBatch(schema, batch)
	if err3 != nil {
		return err3
	}
	for _, sequence := range sequences {
		err4 := dbm.migrateSequenceDown(schema, sequence)
		if err4 != nil {
			return err4
		}

	}
	return dbm.Tgt.DeleteMigrationLog(schema, batch)
}

//MigrateUp  Migrates all the schemas Up
func (dbm *DBMigrationImplementation) MigrateUp() error {
	schemasT0Migrate, err1 := dbm.Src.GetSchemaList()
	if err1 != nil {
		return err1
	}
	for _, schema := range schemasT0Migrate {
		err2 := dbm.MigrateSchemaUp(schema)
		if err2 != nil {
			return err2
		}
	}
	return nil
}

//MigrateDown  Migrates all the schemas Down
func (dbm *DBMigrationImplementation) MigrateDown() error {
	schemasT0Migrate, err1 := dbm.Src.GetSchemaList()
	if err1 != nil {
		return err1
	}
	for _, schema := range schemasT0Migrate {
		err2 := dbm.MigrateSchemaDown(schema)
		if err2 != nil {
			return err2
		}
	}
	return nil
}
