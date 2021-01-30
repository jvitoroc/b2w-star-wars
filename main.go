package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/jvitoroc/b2w-star-wars/resources"
	"github.com/jvitoroc/b2w-star-wars/utils"
)

// middleware para configurar os headers basicos
func setBasics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}

	ctx, cancel := utils.WithTimeout(10)
	defer cancel()
	db, disconnect, err := utils.ConnectMongoDB(ctx)

	if err != nil {
		panic(err)
	}

	defer disconnect()

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"POST, GET, PATCH, DELETE"})

	r := mux.NewRouter()
	r.StrictSlash(false)
	r.Use(setBasics)

	// inicializa os handlers e models da API
	resources.Initialize(r, db)

	http.ListenAndServe(":"+os.Getenv("API_PORT"), handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(r))
}
