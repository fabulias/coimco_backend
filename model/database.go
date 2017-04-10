package model

import (
	"database/sql"
	// _ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"
	"log"
)

var err error
var dbmap = initDb()

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg)
	}
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("sqlite3", "db.local")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'XXX' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Client{}, "customer")
	dbmap.AddTableWithName(Product{}, "product")
	dbmap.AddTableWithName(Provider{}, "provider")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
	return dbmap
}
