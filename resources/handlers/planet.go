package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jvitoroc/b2w-star-wars/resources/common"
	"github.com/jvitoroc/b2w-star-wars/resources/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	sr.Handle("/", appHandler(getPlanetHandler)).Queries("search", "{search}").Methods("GET")
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

	var id primitive.ObjectID
	var err *common.Error
	var filmsAppearedIn int

	if filmsAppearedIn, err = getFilmsAppearedIn(*requestBody.Name); err != nil {
		return err
	}

	if id, err = repo.CreatePlanet(*requestBody.Name, *requestBody.Climate, *requestBody.Terrain, filmsAppearedIn); err != nil {
		return err
	}

	var planet *repo.Planet

	if planet, err = repo.GetPlanetByID(id); err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "Planet successfully created.",
			"planet":  planet,
		},
		http.StatusCreated,
		w,
	)

	return nil
}

func getPlanetByIDHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	id, _ := extractParam("id", r)

	var oid primitive.ObjectID
	var _err error

	if oid, _err = primitive.ObjectIDFromHex(id); _err != nil {
		return common.CreateGenericBadRequestError(_err)
	}

	var planet *repo.Planet
	var err *common.Error

	planet, err = repo.GetPlanetByID(oid)

	if err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "Planet successfully retrieved.",
			"planet":  planet,
		},
		http.StatusOK,
		w,
	)

	return nil
}

func getPlanetHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	var results []*repo.Planet
	var err *common.Error

	results, err = repo.GetMatchedPlanets(map[string]string{
		"name": r.URL.Query().Get("search"),
	})

	if err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "Planets successfully retrieved.",
			"results": results,
		},
		http.StatusOK,
		w,
	)

	return nil
}

func getPlanetsHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	planets, err := repo.GetPlanets()

	if err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "Planets successfully retrieved.",
			"planets": planets,
		},
		http.StatusOK,
		w,
	)

	return nil
}

func deletePlanetHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	id, _ := extractParam("id", r)

	var oid primitive.ObjectID
	var err error

	if oid, err = primitive.ObjectIDFromHex(id); err != nil {
		return common.CreateGenericBadRequestError(err)
	}

	if err := repo.DeletePlanet(oid); err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "Planet successfully deleted.",
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
