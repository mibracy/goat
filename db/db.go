package db

import (
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

func ConnectDB() *bun.DB {
	sqlconstr := fmt.Sprintf("%s:%s@/casaos", "casaos", "casaos")
	sqldb, err := sql.Open("mysql", sqlconstr)
	if err != nil {
		panic(err)
	}
	return bun.NewDB(sqldb, mysqldialect.New())
}
