package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"rest-api/repositories"
	"rest-api/requests"
	"rest-api/utils/caching"
	"rest-api/utils/crypto"
)

type AuthController struct {
	DB    *sql.DB
	Cache caching.Cache
}

func NewAuthController(db *sql.DB, c caching.Cache) *AuthController {
	return &AuthController{
		DB:    db,
		Cache: c,
	}
}

func (uc *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var rr requests.RegisterRequest
	err := decoder.Decode(&rr)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	isExist, err := repositories.GetExistingUserByEmail(uc.DB, rr.Email)
	if err != nil {
		log.Fatalf("Internal Error: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if isExist {
		http.Error(w, "Email existing", http.StatusBadRequest)
		return
	}

	err = repositories.CreateUser(uc.DB, rr.Email, rr.Name, rr.Password)
	if err != nil {
		log.Fatalf("Add user to database error: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	p := map[string]string{
		"message": "Registration success, please login",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (uc *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var lr requests.LoginRequest
	err := decoder.Decode(&lr)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := repositories.GetPrivateUserDetailsByEmail(uc.DB, lr.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}
		log.Fatalf("Login Error: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	password := crypto.HashPassword(lr.Password, user.Salt)
	if user.Password != password {
		http.Error(w, "Invalid password", http.StatusBadRequest)
		return
	}
	token, err := crypto.GenerateToken(user.ID)
	if err != nil {
		log.Fatalf("Login Error: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	oneMonth := time.Duration(60*60*24*30) * time.Second
	err = uc.Cache.Set(fmt.Sprintf("token_%s", token), strconv.Itoa(user.ID), oneMonth)
	if err != nil {
		log.Fatalf("Login Error: %s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	p := map[string]string{
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
