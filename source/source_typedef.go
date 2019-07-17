package source

//MigrationSource Interface for Giving detials about migration versions
type MigrationSource interface {
	GetSchemaList() ([]string, error)
	GetSortedVersions(schema string) ([]int, error)
	GetMigrationUpFile(schema string, version int) (string, string, error)
	GetMigrationDownFile(schema string, version int) (string, string, error)
}
