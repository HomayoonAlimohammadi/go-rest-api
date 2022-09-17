package migrations

import (
	"backend/helpers"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string
}

type Account struct {
	gorm.Model
	Type    string
	Name    string
	Balance uint
	UserID  uint
}

func connectDB() *gorm.DB {
	workingDir, err := os.Getwd()
	helpers.HandleError(err)
	err = godotenv.Load(workingDir + "/migrations/.env")
	helpers.HandleError(err)
	db, err := gorm.Open("postgres", os.Getenv("DSN"))
	helpers.HandleError(err)
	return db
}

func createAccounts() {
	db := connectDB()

	users := [2]User{
		{Username: "Homayoon", Email: "homayoon.alimohammadi@divar.ir"},
		{Username: "Nooshin", Email: "nooshin.rajabi@sharif.edu"},
	}

	for _, user := range users {
		generatedPassword := helpers.HashAndSalt([]byte(user.Username))
		user.Password = generatedPassword
		db.Create(&user)

		account := Account{Type: "Daily Account", Name: string(user.Username + "'s Account"), Balance: uint(1000 * user.ID), UserID: user.ID}
		db.Create(&account)
	}
	defer db.Close()
}

func Migrate() {
	db := connectDB()
	db.AutoMigrate(&User{}, &Account{})
	defer db.Close()

	createAccounts()
}
