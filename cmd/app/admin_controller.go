package main

import (
	mb "github.com/alldroll/multiarmed-bandit"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var (
	listTemplate = template.Must(template.ParseFiles("cmd/app/public/list.html"))
)

type adminController struct {
	storage mb.Storage
}

func newAdminController(storage mb.Storage) *adminController {
	return &adminController{
		storage: storage,
	}
}

//
func (c *adminController) bindRoutes(r *mux.Router) {
	r.HandleFunc("/admin/experiment/{experiment}/stat/", c.statHandler).Methods("GET")
	r.HandleFunc("/admin/experiment/", c.addHandler).Methods("POST")
	r.HandleFunc("/admin/experiments/", c.listHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./cmd/app/public/")))
}

func (c *adminController) statHandler(w http.ResponseWriter, r *http.Request) {

}

func (c *adminController) addHandler(w http.ResponseWriter, r *http.Request) {

}

func (c *adminController) listHandler(w http.ResponseWriter, r *http.Request) {
	list, err := c.storage.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listTemplate.Execute(w, struct {
		Experiments []mb.Experiment
	}{
		Experiments: list,
	})
}
