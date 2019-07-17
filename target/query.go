package target

//MigrationQueries interface to get all the required queries.
type MigrationQueries interface {
	CountMigrationTableSQL() string
	CreateMigrationTableSQL(schema string) string
	GetMaxSequenceSQL(schema string) string
	GetLatestBatchSQL(schema string) string
	InsertMigrationLogSQL(schema string) string
	DeleteMigrationLogSQL(schema string) string
	SetSchemaSQL(schema string) string
	GetSequenceByBatchSQL(schema string) string
	CountSchemaSQL() string
	CreateSchemaSQL(schema string)string
}
