package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"rest-api/repositories"
	"rest-api/utils/caching"
)

type UserController struct {
	DB    *sql.DB
	Cache caching.Cache
}

func NewUserController(db *sql.DB, c caching.Cache) *UserController {
	return &UserController{
		DB:    db,
		Cache: c,
	}
}

func (uc *UserController) Lists(w http.ResponseWriter, r *http.Request) {
	var err error
	page := 1
	pageStr, ok := r.URL.Query()["page"]
	if ok {
		page, err = strconv.Atoi(pageStr[0])
		if err != nil {
			page = 1
		}
	}

	resultsPerPage := 10
	resultsPerPageStr, ok := r.URL.Query()["results_per_page"]
	if ok {
		resultsPerPage, err = strconv.Atoi(resultsPerPageStr[0])
		if err != nil {
			resultsPerPage = 1
		}
	}

	users, err := repositories.GetUsers(uc.DB, page, resultsPerPage)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}
