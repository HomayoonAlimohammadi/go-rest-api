package useraccounts

import (
	"backend/helpers"
	"backend/interfaces"
	"strconv"
)

func updateAccount(id uint, amount int) interfaces.ResponseAccount {
	db := helpers.ConnectDB()
	defer db.Close()
	account := interfaces.Account{}
	responseAccount := interfaces.ResponseAccount{}

	// db.Model(&interfaces.Account{}).Where("id = ?", id).Update("balance", amount)

	db.Where("id = ?", id).First(&account)
	account.Balance = uint(amount)
	db.Save(&account)

	responseAccount.ID = account.ID
	responseAccount.Name = account.Name
	responseAccount.Balance = account.Balance
	return responseAccount
}

func getAccount(id uint) *interfaces.Account {
	db := helpers.ConnectDB()
	defer db.Close()
	account := &interfaces.Account{}
	if db.Where("id = ?", id).First(&account).RecordNotFound() {
		return nil
	}
	return account
}

func Transaction(userID uint, from uint, to uint, amount int, jwt string) map[string]interface{} {
	userIDString := strconv.Itoa(int(userID))
	isValid := helpers.ValidateToken(userIDString, jwt)
	if !isValid {
		return map[string]interface{}{"message": "Not valid token"}
	}

	fromAccount := getAccount(from)
	toAccount := getAccount(to)

	if fromAccount == nil || toAccount == nil {
		return map[string]interface{}{"message": "Accounts not found"}
	} else if fromAccount.UserID != userID {
		return map[string]interface{}{"message": "You are not the owner of the account"}
	} else if int(fromAccount.Balance) < amount {
		return map[string]interface{}{"message": "Not enough money in your account"}
	}

	updatedAccount := updateAccount(from, int(fromAccount.Balance)-amount)
	updateAccount(to, int(toAccount.Balance)+amount)

	// transactions.CreateTransaction(from, to, amount)

	var response = map[string]interface{}{"message": "All is fine"}
	response["data"] = updatedAccount
	return response
}
