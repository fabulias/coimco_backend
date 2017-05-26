package model

import (
	"database/sql"
	"log"

	"coimco_backend/hash"
	"github.com/kimxilxyong/gorp"
	_ "github.com/mattn/go-sqlite3"
)

var err error
var dbmap = initDb()

//Initialize database
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
	dbmap.AddTableWithName(User_acc{}, "user_acc")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	//Create admin account
	var in User_acc
	in.Name = Name
	in.Lastname = Lastname
	in.Mail = Mail
	hash_pass, _ := hash.HashPassword(Pass)
	in.Pass = hash_pass
	in.Rut = Rut
	in.Role = Role
	in.Active = Active
	err = dbmap.Insert(&in)
	if err != nil {
		log.Println(ErrorAdminAccount)
	}
	return dbmap
}
