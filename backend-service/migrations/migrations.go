package migrations

import (
	"backend/helpers"

	"backend/interfaces"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func CreateDefaultAccounts() {
	db := helpers.ConnectDB()
	defer db.Close()
	users := [2]interfaces.User{
		{Username: "Homayoon", Email: "homayoon.alimohammadi@divar.ir"},
		{Username: "Nooshin", Email: "nooshin.rajabi@sharif.edu"},
	}

	for _, user := range users {
		generatedPassword := helpers.HashAndSalt([]byte(user.Username))
		user.Password = generatedPassword
		db.Create(&user)

		account := interfaces.Account{Type: "Daily Account", Name: string(user.Username + "'s Account"), Balance: uint(1000 * user.ID), UserID: user.ID}
		db.Create(&account)
	}
}

func Migrate() {
	db := helpers.ConnectDB()
	defer db.Close()
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec("DROP TABLE IF EXISTS accounts")
	db.AutoMigrate(&interfaces.User{}, &interfaces.Account{})
}
