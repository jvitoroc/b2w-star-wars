package handlers

import (
	"github.com/gorilla/mux"
)

func Initialize(r *mux.Router) {
	initializePlanet(r)
}
