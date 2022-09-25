package transactions

import (
	"backend/database"
	"backend/helpers"
	"backend/interfaces"
)

func CreateTransaction(From uint, To uint, Amount int) {
	db := database.DB

	transaction := &interfaces.Transaction{From: From, To: To, Amount: Amount}
	db.Create(&transaction)
}

func GetTransactionsByAccount(id uint) []interfaces.ResponseTransaction {
	transactions := []interfaces.ResponseTransaction{}

	database.DB.Table("transactions").Select("id, transactions.from, transactions.to, amount").Where(interfaces.Transaction{From: id}).Or(interfaces.Transaction{To: id}).Scan(&transactions)

	return transactions

}

func GetMyTransactions(id string, jwt string) map[string]interface{} {
	isValid := helpers.ValidateToken(id, jwt)
	if !isValid {
		return map[string]interface{}{"message": "Invalid token"}
	}
	accounts := []interfaces.ResponseAccount{}
	database.DB.Table("accounts").Select("id, name, balance").Where("user_id = ?", id).Scan(&accounts)
	transactions := []interfaces.ResponseTransaction{}
	for i := 0; i < len(accounts); i++ {
		accTransactions := GetTransactionsByAccount(accounts[i].ID)
		transactions = append(transactions, accTransactions...)
	}
	var response = map[string]interface{}{"message": "All is fine"}
	response["data"] = transactions

	return response
}
