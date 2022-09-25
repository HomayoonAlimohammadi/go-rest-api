package database

import (
	"backend/helpers"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func InitDatabase() {
	workingDir, err := os.Getwd()
	helpers.HandleError(err)
	err = godotenv.Load(workingDir + "/migrations/.env")
	helpers.HandleError(err)
	db, err := gorm.Open("postgres", os.Getenv("DSN"))
	helpers.HandleError(err)

	db.DB().SetMaxIdleConns(20)
	db.DB().SetMaxOpenConns(200)
	DB = db

}
