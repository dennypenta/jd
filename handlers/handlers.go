package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"fmt"
	"io"
	"io/ioutil"
	"jd/models"
	"strconv"
)

func ListVacancy(w http.ResponseWriter, r *http.Request) {
	var vacansies []models.Vacancy
	models.Database.Connect.Find(&vacansies)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(vacansies); err != nil {
		panic(err)
	}
}

func RetrieveVacancy(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	i, err := strconv.Atoi(id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	vacancy := &models.Vacancy{}
	models.Database.Get(vacancy, i)

	if vacancy.ID == 0 {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(vacancy); err != nil {
		panic(err)
	}
}

func CreateVacancy(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	vacancy := models.Vacancy{}
	if err := json.Unmarshal(body, &vacancy); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	models.Database.Connect.Create(&vacancy)
	if models.Database.Connect.NewRecord(vacancy) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, body)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(vacancy); err != nil {
		panic(err)
	}
}
