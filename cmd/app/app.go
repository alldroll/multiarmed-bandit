package main

import (
	"database/sql"
	"fmt"
	mb "github.com/alldroll/multiarmed-bandit"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

//
func runApp() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	db, err := createDB(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_NAME"),
	)

	if err != nil {
		log.Fatal(err)
	}

	service, err := mb.NewService(db)

	if err != nil {
		log.Fatal(err)
	}

	runScheduler(service)

	userController := &userController{
		algo: service.GetAlgorithm(),
	}

	r := mux.NewRouter()
	userController.bindRoutes(r)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

//
func createDB(user, password, host, db string) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		user,
		password,
		host,
		db,
	)

	return sql.Open("mysql", connectionString)
}

//
func runScheduler(service mb.Service) {
	go func() {
		for range time.NewTicker(2 * time.Second).C {
			service.Update()
		}
	}()
}
