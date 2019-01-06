package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"rest-api/repositories"
	"rest-api/requests"
	"rest-api/utils/caching"
	"rest-api/utils/upload"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
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

func (uc *UserController) Update(w http.ResponseWriter, r *http.Request) {
	jobID, err := strconv.Atoi(path.Base(r.URL.Path))
	decoder := json.NewDecoder(r.Body)
	var ur requests.UpdateUserRequest
	err = decoder.Decode(&ur)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	target_url := "http://localhost:9090/upload"
	filename := "test.pdf"
	upload.PostFile(filename, target_url)

	err = repositories.UpdateUser(uc.DB, jobID, ur.Email, ur.Name, ur.Password)

	if err != nil {
		http.Error(w, "Error parsing token", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)

}

func (uc *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	var token string
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
		token = strings.TrimPrefix(token, "Bearer ")
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, msg
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		http.Error(w, "Error parsing token", http.StatusUnauthorized)
		return
	}
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userID, _ := claims["user_id"].(float64)
		user, err := repositories.GetUserByID(uc.DB, int(userID))
		if err != nil {
			http.Error(w, "Error parsing token", http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(user)
	} else {
		fmt.Println(err)
	}
}
