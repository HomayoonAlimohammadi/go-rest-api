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

// Login attemps to log user in with the given username and password
func Login(username string, password string) map[string]interface{} {

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

	// Setup response
	responseUser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}

	// Sign JWT Token
	tokenContent := jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}
	godotenv.Load("./users/.env")
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte(os.Getenv("JWT_PASSWORD")))
	helpers.HandleError(err)

	// Prepare response
	var response = map[string]interface{}{"message": "All is fine"}
	response["jwt"] = token
	response["data"] = responseUser

	return response
}
