package target

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

// GetDBDir return migration dir
func GetDBMigrationDir() string {
	return getEnv("DB_MIGRATION_DIR", "./migrations")
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDBURL returns DB URL string
func GetDBURL() string {
	connStr := os.Getenv("DB_URL")
	if connStr != "" {
		return connStr
	}
	host := getEnv("DB_HOST", "localhost")
	password := getEnv("DB_PASSWORD", "postgres")
	user := getEnv("DB_USER", "postgres")
	dbname := getEnv("DB_NAME", "postgres")
	port := getEnv("DB_PORT", "5432")
	isSSLEnabled, _ := strconv.ParseBool(os.Getenv("DB_SSL_ENABLED"))
	caPath := os.Getenv("DB_SSL_CA_PATH")
	certPath := os.Getenv("DB_SSL_CERT_PATH")
	keyPath := os.Getenv("DB_SSL_KEY_PATH")
	sslmode := getEnv("DB_SSL_MODE", "disable")
	schema := os.Getenv("DB_SCHEMA")
	pgStr := "postgres://%s:%s@%s:%s/%s?search_path=%s&sslmode=%s"
	connStr = fmt.Sprintf(pgStr, user, password, host, port, dbname, schema, sslmode)
	if isSSLEnabled {
		connStr = fmt.Sprintf(pgStr+"&sslrootcert=%s&sslkey=%s&sslcert=%s", user, password, host, port, dbname, schema, sslmode, caPath, keyPath, certPath)
	}
	return connStr
}

// GetDatabase returns DB implementation
func GetDatabase() (*DatabaseImplementation, error) {
	dbMaxConn, _ := strconv.Atoi(getEnv("DB_MAX_CONNECTIONS", "1"))
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
