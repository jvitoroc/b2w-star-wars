package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jvitoroc/b2w-star-wars/resources/common"
)

type SWAPIPlanet struct {
	Name  string   `json:"name"`
	Films []string `json:"films"`
}

type SWAPISearchResult struct {
	Results []SWAPIPlanet `json:"results"`
}

type appHandler func(http.ResponseWriter, *http.Request) *common.Error

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

func getFilmsAppearedIn(planetName string) (int, *common.Error) {
	planetName = strings.Trim(strings.ToLower(planetName), "\n\r ")
	resp, err := http.Get("https://swapi.dev/api/planets?search=" + planetName)

	if err != nil {
		return 0, common.CreateGenericInternalError(err)
	}

	var result SWAPISearchResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, common.CreateGenericInternalError(err)
	}

	if resCount := len(result.Results); resCount > 0 {
		for i := 0; i < resCount; i++ {
			if strings.Trim(strings.ToLower(result.Results[i].Name), " ") == planetName {
				return len(result.Results[i].Films), nil
			}
		}
	}

	return 0, nil
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
