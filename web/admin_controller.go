package web

import (
	mb "github.com/alldroll/multiarmed-bandit"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var (
	listTemplate = template.Must(template.ParseFiles(
		"web/public/experiment_list.tmpl",
		"web/public/header.tmpl",
		"web/public/footer.tmpl",
	))

	experimentTemplate = template.Must(template.ParseFiles(
		"web/public/experiment.tmpl",
		"web/public/header.tmpl",
		"web/public/footer.tmpl",
	))
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
	r.HandleFunc("/admin/experiment/{experiment}/", c.experimentHandler).Methods("GET")
	r.HandleFunc("/admin/experiment/", c.addHandler).Methods("POST")
	r.HandleFunc("/admin/experiments/", c.listHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./web/public/static")))
}

//
func (c *adminController) experimentHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars = mux.Vars(r)
		name = vars["experiment"]
	)

	experiment, err := c.storage.Find(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if experiment == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	experimentTemplate.Execute(w, struct {
		Experiment mb.Experiment
	}{
		Experiment: experiment,
	})
}

//
func (c *adminController) statHandler(w http.ResponseWriter, r *http.Request) {

}

//
func (c *adminController) addHandler(w http.ResponseWriter, r *http.Request) {

}

//
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
