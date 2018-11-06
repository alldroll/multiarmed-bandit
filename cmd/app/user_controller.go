package main

import (
	"encoding/json"
	mb "github.com/alldroll/multiarmed-bandit"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type userController struct {
	algo mb.Algorithm
}

func newUserController(algo mb.Algorithm) *userController {
	return &userController{
		algo: algo,
	}
}

//
func (c *userController) bindRoutes(r *mux.Router) {
	r.HandleFunc("/experiment/{experiment}/{choice}/", c.rewardHandler).Methods("PUT")
	r.HandleFunc("/experiment/{experiment}/", c.choiceHandler).Methods("GET")
}

//
func (c *userController) rewardHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars       = mux.Vars(r)
		experiment = vars["experiment"]
		variant    = vars["choice"]
	)

	v, err := strconv.ParseUint(variant, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.algo.Reward(experiment, uint32(v))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(struct {
		bool `json:"success"`
	}{
		true,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

//
func (c *userController) choiceHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars       = mux.Vars(r)
		experiment = vars["experiment"]
	)

	variant, err := c.algo.Suggest(experiment)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := c.algo.Show(experiment, variant); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(struct {
		Variant uint32 `json:"variant"`
	}{
		variant,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
