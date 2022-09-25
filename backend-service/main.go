package main

import (
	"backend/api"
	"backend/database"
	"backend/migrations"
)

func main() {
	database.InitDatabase()
	migrations.Migrate()
	migrations.CreateDefaultAccounts()
	api.StartApp()

}
