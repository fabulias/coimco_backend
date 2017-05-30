package model

import (
	"log"

	"coimco_backend/hash"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/jinzhu/gorm/dialects/postgres"
)

var err error
var dbmap = initDb()

//Initialize database
func initDb() *gorm.DB {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	log.Println("Initialize database")
	db, err := gorm.Open("sqlite3", "local.db")
	//LogMode is active
	db.LogMode(true)
	//defer db.Close()
	if err != nil {
		checkErr(err, err.Error())
	}
	db.SingularTable(true)
	db.AutoMigrate(Customer{}, Provider{}, Product{}, User_acc{}, Tag{}, Tag_customer{})
	//db.Model(&Tag_customer{}).AddForeignKey("CustomerID", "tag(id)", "RESTRICT", "RESTRICT")
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
	if db.NewRecord(in) {
		db.Create(in)
	}
	if err != nil {
		log.Println(ErrorAdminAccount)
	}
	return db
}
