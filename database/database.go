package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func LoadDB(){
	Connect()

	// adding all schemas to db
	Database.AutoMigrate(&User{})
	Database.AutoMigrate(&Options{})
	Database.AutoMigrate(&BalanceAttribute{})
	Database.AutoMigrate(&MonthlyExpense{})
	Database.AutoMigrate(&Expense{})
}

func Connect(){
	var err error
	host := os.Getenv("POSTGRES_HOSTNAME")
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	credentials := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=EST", host, username, password, dbName, port)
	Database, err = gorm.Open(postgres.Open(credentials), &gorm.Config{})

	if err != nil{
		panic(err)
	}// else{
	// 	fmt.Println("Successfully connected to database")
	// }
}