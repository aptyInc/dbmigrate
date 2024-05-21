package target

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

// GetDBDir return migration dir
func GetDBMigrationDir() string {
	return os.Getenv("DB_MIGRATION_DIR")
}

// GetDBURL returns DB URL string
func GetDBURL() string {
	host := os.Getenv("DB_HOST")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	isSSLEnabled, _ := strconv.ParseBool(os.Getenv("DB_SSL_ENABLED"))
	caPath := os.Getenv("DB_SSL_CA_PATH")
	certPath := os.Getenv("DB_SSL_CERT_PATH")
	keyPath := os.Getenv("DB_SSL_KEY_PATH")
	sslmode := os.Getenv("DB_SSL_MODE")
	schema := os.Getenv("DB_SCHEMA")
	pgStr := "postgres://%s:%s@%s:%s/%s?search_path=%s&sslmode=%s"
	var connStr = fmt.Sprintf(pgStr, user, password, host, port, dbname, schema, sslmode)
	if isSSLEnabled {
		connStr = fmt.Sprintf(pgStr+"&sslmode=%s&sslrootcert=%s&sslkey=%s&sslcert=%s", user, password, host, port, dbname, schema, sslmode, caPath, keyPath, certPath)
	}
	return connStr
}

// GetDatabase returns DB implementation
func GetDatabase() (*DatabaseImplementation, error) {
	dbMaxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	connStr := GetDBURL()
	if os.Getenv("DB_ENVIRONMENT") == "development" {
		fmt.Println(connStr)
	}
	db, err4 := sql.Open("postgres", connStr)
	db.SetMaxIdleConns(dbMaxConn)
	if err4 != nil {
		return nil, err4
	}
	err5 := db.Ping()
	if err5 != nil {
		return nil, err5
	}
	impl := &DatabaseImplementation{
		mq: Postgres{},
		DB: db,
	}
	return impl, nil

}
