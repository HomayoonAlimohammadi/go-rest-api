package helpers

import (
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func HandleError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func HashAndSalt(password []byte) string {
	hashed, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	HandleError(err)
	return string(hashed)
}

func ConnectDB() *gorm.DB {
	workingDir, err := os.Getwd()
	HandleError(err)
	err = godotenv.Load(workingDir + "/migrations/.env")
	HandleError(err)
	db, err := gorm.Open("postgres", os.Getenv("DSN"))
	HandleError(err)
	return db
}
