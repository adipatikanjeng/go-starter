package main

import (
	"log"
	"net/http"
	"os"
	"rest-api/utils/middleware"

	"rest-api/controllers"
	"rest-api/routes"
	"rest-api/utils/caching"
	"rest-api/utils/database"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

	mux := mux.NewRouter()
	amw := middleware.AuthMiddleware(cache)
	mux.Use(amw.Middleware) //middleware for api
	routes.CreateRoutes(mux, authController, userController, jobController)

	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}
