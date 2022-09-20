package useraccounts

import (
	"backend/helpers"
	"backend/interfaces"
)

func updateAccout(id uint, amount int) {
	db := helpers.ConnectDB()
	defer db.Close()

	db.Model(&interfaces.Account{}).Where("id = ?", id).Update("balance", amount)
}
