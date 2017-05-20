package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v2"
	"log"
	"strings"
)

var err error
var dbmap = initDb()

func CheckInCustomer(in Customer) bool {
	if strings.Compare(in.Name, "") != 0 && strings.Compare(in.Rut, "") != 0 && strings.Compare(in.Mail, "") != 0 {
		return true
	} else {
		return false
	}
}

func CheckInProduct(in Product) bool {
	var flag bool = false
	if strings.Compare(in.Name, "") != 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Details, "") != 0 {
		flag = true
		return flag
	} else if in.Stock < 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Brand, "") != 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Category, "") != 0 {
		flag = true
		return flag
	}
	return flag
}

/*
Stock Brand Category
*/

func checkErr(err error, msg string) {
	if err != nil {
		log.Println(msg)
	}
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	log.Println("Initialize database")
	db, err := sql.Open("sqlite3", "local.db")
	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	//return dbmap
	// add a table, setting the table name to 'XXX' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Customer{}, "customer")
	dbmap.AddTableWithName(Product{}, "product")
	dbmap.AddTableWithName(Provider{}, "provider")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
	return dbmap
}
