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
	src *source.FileSource
	tgt *target.DatabaseImplementation
}

func (dbm *DBMigrationImplementation) initGetMaxSequence(schema string) (int, error) {
	isFound, err1 := dbm.tgt.DoMigrationTableExists(schema)
	if err1 != nil {
		return -1, err1
	}
	if !isFound {
		err1 := dbm.tgt.CreateMigrationTable(schema)
		return -1, err1
	}
	return dbm.tgt.GetMaxSequence(schema)
}
func (dbm *DBMigrationImplementation) migrateSequenceUp(schema string, version int, batch string) error {
	name, commands, err1 := dbm.src.GetMigrationUpFile(schema, version)

	if err1 != nil {
		return err1
	}
	err2 := dbm.tgt.ExecuteMigration(schema, commands)
	if err2 != nil {
		return err2
	}

	return dbm.tgt.InsertMigrationLog(schema, version, name, batch)

}

func (dbm *DBMigrationImplementation) migrateSequenceDown(schema string, version int) error {
	_, commands, err1 := dbm.src.GetMigrationDownFile(schema, version)

	if err1 != nil {
		return err1
	}
	return dbm.tgt.ExecuteMigration(schema, commands)

}

func (dbm *DBMigrationImplementation) MigrateSchemaUp(schema string) error {

	maxSequence, err1 := dbm.initGetMaxSequence(schema)
	if err1 != nil {
		return err1
	}
	allVersions, err2 := dbm.src.GetSortedVersions(schema)
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

	doExisit, err1 := dbm.tgt.DoMigrationTableExists(schema)
	if err1 != nil {
		return err1
	}
	if !doExisit {
		return fmt.Errorf("no Migrations in schema")
	}

	batch, err2 := dbm.tgt.GetLatestBatch(schema)
	if err2 != nil {
		return err2
	}
	sequences, err3 := dbm.tgt.GetSequenceByBatch(schema, batch)
	if err3 != nil {
		return err3
	}
	for _, sequence := range sequences {
		err4 := dbm.migrateSequenceDown(schema, sequence)
		if err4 != nil {
			return err4
		}

	}
	return dbm.tgt.DeleteMigrationLog(schema, batch)
}

//MigrateUp  Migrates all the schemas Up
func (dbm *DBMigrationImplementation) MigrateUp() error {
	schemasT0Migrate, err1 := dbm.src.GetSchemaList()
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
	schemasT0Migrate, err1 := dbm.src.GetSchemaList()
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
