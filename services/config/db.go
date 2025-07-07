package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func ConnectDB() *bun.DB {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1" // Default to localhost
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306" // Default MySQL port
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "casaos" // Default user
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "casaos" // Default password
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "casaos" // Default database name
	}

	sqlconstr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	sqldb, err := sql.Open("mysql", sqlconstr)
	if err != nil {
		panic(err)
	}
	return bun.NewDB(sqldb, mysqldialect.New())
}
