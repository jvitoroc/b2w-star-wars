package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jvitoroc/b2w-star-wars/resources/common"
	"github.com/jvitoroc/b2w-star-wars/resources/repo"
)

type PlanetRequestBody struct {
	Name    *string `json:"name"`
	Climate *string `json:"climate"`
	Terrain *string `json:"terrain"`
}

func initializePlanet(r *mux.Router) {
	sr := r.PathPrefix("/planet").Subrouter()
	sr.Handle("/", appHandler(createPlanetHandler)).Methods("POST")
	sr.Handle("/{id:[a-z0-9]+}", appHandler(getPlanetByIDHandler)).Methods("GET")
	sr.Handle("/", appHandler(getMatchedPlanetHandler)).Queries("search", "{search}").Methods("GET")
	sr.Handle("/", appHandler(getPlanetsHandler)).Methods("GET")
	sr.Handle("/{id:[a-z0-9]+}", appHandler(deletePlanetHandler)).Methods("DELETE")
}

func createPlanetHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	requestBody := PlanetRequestBody{}

	if err := extractPlanet(&requestBody, r); err != nil {
		return err
	}

	if err := validatePlanet(&requestBody); err != nil {
		return err
	}

	filmsAppearedIn, err := getFilmsAppearedIn(*requestBody.Name)
	if err != nil {
		return err
	}

	id, err := repo.CreatePlanet(*requestBody.Name, *requestBody.Climate, *requestBody.Terrain, filmsAppearedIn)
	if err != nil {
		return err
	}

	planet, err := repo.GetPlanetByID(id)
	if err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "The planet was successfully created.",
			"planet":  planet,
		},
		http.StatusCreated,
		w,
	)

	return nil
}

func getPlanetByIDHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	id, _ := extractParam("id", r)

	oid, err := stringToObjectID(id)
	if err != nil {
		return err
	}

	planet, err := repo.GetPlanetByID(*oid)
	if err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "The planet was successfully retrieved.",
			"planet":  planet,
		},
		http.StatusOK,
		w,
	)

	return nil
}

func getMatchedPlanetHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	results, err := repo.GetMatchedPlanets(map[string]string{
		"name": r.URL.Query().Get("search"),
	})

	if err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "The planets were successfully retrieved.",
			"results": results,
		},
		http.StatusOK,
		w,
	)

	return nil
}

func getPlanetsHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	planets, err := repo.GetAllPlanets()

	if err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "The planets were successfully retrieved.",
			"planets": planets,
		},
		http.StatusOK,
		w,
	)

	return nil
}

func deletePlanetHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	id, _ := extractParam("id", r)

	oid, err := stringToObjectID(id)
	if err != nil {
		return err
	}

	if err := repo.DeletePlanet(*oid); err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "The planet was successfully deleted.",
		},
		http.StatusOK,
		w,
	)

	return nil
}

func extractPlanet(planet *PlanetRequestBody, r *http.Request) *common.Error {
	err := json.NewDecoder(r.Body).Decode(planet)

	if err != nil {
		return common.CreateGenericBadRequestError(err)
	}

	return nil
}

func validatePlanet(planet *PlanetRequestBody) *common.Error {
	errors := map[string]string{}

	if planet.Name == nil || *planet.Name == "" {
		errors["name"] = "Name field is empty or missing."
	}

	if planet.Climate == nil || *planet.Climate == "" {
		errors["climate"] = "Climate field is empty or missing."
	}

	if planet.Terrain == nil || *planet.Terrain == "" {
		errors["terrain"] = "Terrain field is empty or missing."
	}

	if len(errors) == 0 {
		return nil
	} else {
		return common.CreateFormError(errors)
	}
}
