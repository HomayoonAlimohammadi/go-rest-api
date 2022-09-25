package transactions

import (
	"backend/helpers"
	"backend/interfaces"
)

func CreateTransaction(From uint, To uint, Amount int) {
	db := helpers.ConnectDB()
	defer db.Close()

	transaction := &interfaces.Transaction{From: From, To: To, Amount: Amount}
	db.Create(&transaction)
}
