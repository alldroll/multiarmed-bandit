package web

import (
	"database/sql"
	"fmt"
	mb "github.com/alldroll/multiarmed-bandit"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type App struct {
	config AppConfig
}

type AppConfig struct {
	Port       string
	DBUserName string
	DBPassword string
	DBHost     string
	DBName     string
}

//
func NewApp(config AppConfig) App {
	return App{
		config: config,
	}
}

//
func (a App) Run() {
	port := a.config.Port
	if port == "" {
		log.Fatal("Port must be set")
	}

	db, err := a.createDB()

	if err != nil {
		log.Fatal(err)
	}

	service, err := mb.NewService(db)

	if err != nil {
		log.Fatal(err)
	}

	a.runScheduler(service)

	userController := newUserController(service.GetAlgorithm())
	adminController := newAdminController(service.GetStorage())

	r := mux.NewRouter()
	r.StrictSlash(true)

	userController.bindRoutes(r)
	adminController.bindRoutes(r)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

//
func (a App) createDB() (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s",
		a.config.DBUserName,
		a.config.DBPassword,
		a.config.DBHost,
		a.config.DBName,
	)

	return sql.Open("mysql", connectionString)
}

//
func (a App) runScheduler(service mb.Service) {
	go func() {
		for range time.NewTicker(2 * time.Second).C {
			service.Update()
		}
	}()
}
