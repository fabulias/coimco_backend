package model

import (
	"database/sql"
	"log"
	"strings"

	"github.com/kimxilxyong/gorp"
	_ "github.com/mattn/go-sqlite3"
)

var err error
var dbmap = initDb()

//Return true in case of that all params are okay
func CheckInCustomer(in Customer) bool {
	if strings.Compare(in.Name, "") != 0 && strings.Compare(in.Rut, "") != 0 && strings.Compare(in.Mail, "") != 0 {
		return true
	} else {
		return false
	}
}

//Return true in case of that all params are okay
func CheckInAccount(in User_acc) bool {
	var flag bool = false
	if strings.Compare(in.Name, "") == 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Lastname, "") == 0 {
		flag = true
		return flag
	} else if in.Role != false && in.Role != true {
		flag = true
		return flag
	} else if strings.Compare(in.Mail, "") == 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Rut, "") == 0 {
		flag = true
		return flag
	} else if strings.Compare(in.Pass, "") == 0 {
		flag = true
		return flag
	}
	return flag
}

//Return true in case of that all params are okay
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

func checkErr(err error, msg string) {
	if err != nil {
		log.Println(msg)
	}
}

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
	//dbmap.AddTableWithName(Product{}, "product")
	//dbmap.AddTableWithName(Provider{}, "provider")
	//dbmap.AddTableWithName(User_acc{}, "user_acc")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")

	//Create admin account
	var in User_acc
	in.Name = Name
	in.Lastname = Lastname
	in.Mail = Mail
	in.Pass = Pass
	in.Rut = Rut
	in.Role = Role
	in.Active = Active
	err = dbmap.Insert(&in)
	if err != nil {
		log.Println(ErrorAdminAccount)
	}
	return dbmap
}
