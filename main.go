package main

import (
	"log"
	"net/http"
	"os"

	"rest-api/controllers"
	"rest-api/routes"
	"rest-api/utils/caching"
	"rest-api/utils/database"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/urfave/negroni"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.Connect(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	if err != nil {
		log.Fatal(err)
	}
	cache := &caching.Redis{
		Client: caching.Connect(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), 0),
	}

	authController := controllers.NewAuthController(db, cache)
	userController := controllers.NewUserController(db, cache)
	jobController := controllers.NewJobController(db, cache)

	apiMux := mux.NewRouter()
	mux := mux.NewRouter()

	mw := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	routes.CreateAuthRoutes(mux, authController)
	routes.CreateAPIRoutes(apiMux, userController, jobController)

	an := negroni.New(negroni.HandlerFunc(mw.HandlerWithNext), negroni.Wrap(apiMux))
	mux.PathPrefix("/api/v1").Handler(an)
	n := negroni.Classic()
	n.UseHandler(mux)

	if err := http.ListenAndServe(":8081", n); err != nil {
		log.Fatal(err)
	}
}
