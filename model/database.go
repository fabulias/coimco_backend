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
	db.AutoMigrate(Customer{}, Provider{}, Product{}, User_acc{}, Tag_customer{}, Tag{})

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

/*
type Profile struct {
    Refer int
    Name  string
}
//Tag_customer
type User struct {
    Profile   Profile `gorm:"ForeignKey:ProfileID;AssociationForeignKey:Refer"`
    ProfileID int
}*/
