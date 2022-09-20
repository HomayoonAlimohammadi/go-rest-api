package users

import (
	"backend/helpers"
	"backend/interfaces"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func prepareToken(user *interfaces.User) string {
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}
	godotenv.Load("./users/.env")
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte(os.Getenv("JWT_PASSWORD")))
	helpers.HandleError(err)

	return token
}

func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount) map[string]interface{} {

	// Setup response
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	var response = map[string]interface{}{"message": "All is fine"}
	var token = prepareToken(user)
	response["jwt"] = token
	response["data"] = responseUser

	return response
}

// Login attemps to log user in with the given username and password
func Login(username string, password string) map[string]interface{} {

	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: password, Valid: "password"},
		})

	if !valid {
		return map[string]interface{}{"message": "Invalid values"}
	}

	// Connect to DataBase
	db := helpers.ConnectDB()
	defer db.Close()

	// Get user from DB
	user := &interfaces.User{}
	if db.Where("username = ?", username).First(&user).RecordNotFound() {
		return map[string]interface{}{"message": "User not found"}
	}

	// Check passwords match
	passwordError := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if passwordError == bcrypt.ErrMismatchedHashAndPassword && passwordError != nil {
		return map[string]interface{}{"message": "Wrong password"}
	}

	// Find account for the user
	accounts := []interfaces.ResponseAccount{}
	db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

	var response = prepareResponse(user, accounts)

	return response
}

func Register(username string, email string, password string) map[string]interface{} {
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: password, Valid: "password"},
		})

	if !valid {
		return map[string]interface{}{"message": "Invalid values"}
	}

	db := helpers.ConnectDB()
	defer db.Close()

	generatedPassword := helpers.HashAndSalt([]byte(password))
	user := &interfaces.User{Username: username, Email: email, Password: generatedPassword}
	db.Create(&user)

	account := interfaces.Account{Type: "Daily Account", Name: string(username + "'s Account"), Balance: 0, UserID: user.ID}
	db.Create(&account)

	accounts := []interfaces.ResponseAccount{}
	responseAccount := interfaces.ResponseAccount{ID: account.ID, Name: account.Name, Balance: uint(account.Balance)}
	accounts = append(accounts, responseAccount)
	var response = prepareResponse(user, accounts)

	return response
}
