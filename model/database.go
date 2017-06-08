package model

import (
	"log"
	"os"
	"strconv"

	"github.com/fabulias/coimco_backend/hash"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var err error
var Dbmap = initDb()

//Initialize database
func initDb() *gorm.DB {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	log.Println("Initialize database")
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	//"postgres://losaieljggcviq:94c6c9315e714fab"+
	//"5415ed1be76d4a2037881447b75770f62842d5ff4a0f1dac@ec2-107-22-244-62."+
	//"compute-1.amazonaws.com:5432/d2pkqjvdn5eiha")
	//db, err := gorm.Open("sqlite3", "local.db")
	//LogMode is active
	db.LogMode(true)
	//defer db.Close()
	if err != nil {
		checkErr(err, err.Error())
	}
	db.SingularTable(true)
	db.AutoMigrate(Customer{}, Provider{}, Product{},
		UserAcc{}, Tag{}, TagCustomer{}, Sale{},
		SaleDetail{}, Purchase{}, PurchaseDetail{})

	db.Model(&TagCustomer{}).AddForeignKey("tag_id", "tag(id)",
		"RESTRICT", "RESTRICT")
	db.Model(&TagCustomer{}).AddForeignKey("customer_id", "customer(rut)",
		"RESTRICT", "RESTRICT")

	db.Model(&Sale{}).AddForeignKey("customer_id", "customer(rut)",
		"RESTRICT", "RESTRICT")
	db.Model(&Sale{}).AddForeignKey("user_id", "user_acc(mail)",
		"RESTRICT", "RESTRICT")

	db.Model(&SaleDetail{}).AddForeignKey("sale_id", "sale(id)",
		"RESTRICT", "RESTRICT")
	db.Model(&SaleDetail{}).AddForeignKey("product_id", "product(id)",
		"RESTRICT", "RESTRICT")

	db.Model(&Purchase{}).AddForeignKey("provider_id", "provider(rut)",
		"RESTRICT", "RESTRICT")

	db.Model(&PurchaseDetail{}).AddForeignKey("purchase_id", "purchase(id)",
		"RESTRICT", "RESTRICT")
	db.Model(&PurchaseDetail{}).AddForeignKey("product_id", "product(id)",
		"RESTRICT", "RESTRICT")

	//Create admin account
	var in UserAcc
	in.Name = os.Getenv("NAME")
	in.Lastname = os.Getenv("LASTNAME")
	in.Mail = os.Getenv("MAIL")
	hash_pass, _ := hash.HashPassword(os.Getenv("PASSWORD"))
	in.Pass = hash_pass
	in.Rut = os.Getenv("RUT")
	role, _ := strconv.ParseInt(os.Getenv("ROLE"), 10, 8)
	in.Role = int8(role)
	in.Active = true

	db.FirstOrCreate(&in)

	if err != nil {
		log.Println(ErrorAdminAccount)
	}
	return db
}
