package main

import (
	"github.com/alldroll/multiarmed-bandit/web"
	"os"
)

func main() {
	config := web.AppConfig{
		Port:       os.Getenv("PORT"),
		DBUserName: os.Getenv("APP_DB_USERNAME"),
		DBPassword: os.Getenv("APP_DB_PASSWORD"),
		DBHost:     os.Getenv("APP_DB_HOST"),
		DBName:     os.Getenv("APP_DB_NAME"),
	}

	app := web.NewApp(config)
	app.Run()
}
