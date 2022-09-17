package main

import (
	"backend/api"
	"backend/migrations"
)

func main() {
	migrations.Migrate()
	migrations.CreateDefaultAccounts()
	api.StartApp()

}
