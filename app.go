package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jvitoroc/b2w-star-wars/resources"
	"github.com/jvitoroc/b2w-star-wars/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type App struct {
	Router *mux.Router
	DB     *mongo.Client

	HttpHandler http.Handler
}

func (a *App) Initialize(mongoURI string, databaseName string) {
	_, err := a._ConnectMongoDB(mongoURI)
	if err != nil {
		log.Fatal(errors.New("Could not connect to the database server: " + err.Error()))
	}

	a._ConfigureRouter(databaseName)
}

func (a *App) Run(addr string) {
	defer a.DB.Disconnect(context.Background())

	// inicializa o server
	err := a._Listen(addr)
	if err != nil {
		log.Fatal(err)
	}
}

// funcoes que iniciam com _ sao "privadas"
func (a *App) _ConnectMongoDB(uri string) (utils.DisconnectFunc, error) {
	ctx, cancel := utils.WithTimeout(10)
	defer cancel()
	db, disconnect, err := utils.ConnectMongoDB(ctx, uri)

	if err == nil {
		a.DB = db
		if err := a._TestMongoDBConnection(); err != nil {
			return nil, err
		}
	}

	return disconnect, err
}

func (a *App) _TestMongoDBConnection() error {
	ctx, cancel := utils.WithTimeout(10)
	defer cancel()
	return a.DB.Ping(ctx, readpref.PrimaryPreferred())
}

func (a *App) _ConfigureRouter(databaseName string) {
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedHeaders := handlers.AllowedHeaders([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"POST, GET, PATCH, DELETE"})

	a.Router = mux.NewRouter()
	a.Router.StrictSlash(false)
	a.Router.Use(setBasicsMiddleware)

	resources.Initialize(a.Router, a.DB, databaseName)

	a.HttpHandler = handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(a.Router)
}

func (a *App) _Listen(addr string) error {
	return http.ListenAndServe(addr, a.HttpHandler)
}

// middleware para configurar os headers basicos
func setBasicsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
