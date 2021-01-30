package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jvitoroc/b2w-star-wars/resources/common"
)

type appHandler func(http.ResponseWriter, *http.Request) *common.Error

func Initialize(r *mux.Router) {
	initializePlanet(r)
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		respondWithError(*err, w)
	}
}

func respond(data interface{}, statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondWithMessage(message string, statusCode int, w http.ResponseWriter) {
	respond(map[string]string{"message": message}, statusCode, w)
}

func respondWithError(err common.Error, w http.ResponseWriter) {
	respond(err, err.Code, w)
}

func extractParam(param string, r *http.Request) (string, bool) {
	value, ok := mux.Vars(r)[param]
	return value, ok
}

func extractParamInt(param string, r *http.Request) (int, bool) {
	if value, ok := extractParam(param, r); ok {
		if value, err := strconv.Atoi(value); err == nil {
			return value, true
		}
	}

	return 0, false
}
