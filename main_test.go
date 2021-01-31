package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/jvitoroc/b2w-star-wars/resources/common"
	"go.mongodb.org/mongo-driver/bson"
)

type TestPlanet struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Climate         string `json:"climate"`
	Terrain         string `json:"terrain"`
	FilmsAppearedIn int    `json:"filmsAppearedIn"`
}

type TestSinglePlanetResponse struct {
	Message string     `json:"message"`
	Planet  TestPlanet `json:"planet"`
}

type TestMultiplePlanetsResponse struct {
	Message string       `json:"message"`
	Planets []TestPlanet `json:"planets"`
}

type TestMatchedPlanetResponse struct {
	Message string       `json:"message"`
	Results []TestPlanet `json:"results"`
}

var databaseName string
var a App = App{}

var tatooine = TestPlanet{Name: "Tatooine", Climate: "temperate", Terrain: "grasslands, mountains", FilmsAppearedIn: 5}
var tatooineBytes = []byte(fmt.Sprintf(`{"name": "%s", "climate": "%s", "terrain": "%s"}`, tatooine.Name, tatooine.Climate, tatooine.Terrain))

func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	databaseName = os.Getenv("MONGODB_DBNAME_DEV")
	a.Initialize(os.Getenv("MONGODB_URI"), databaseName)

	code := m.Run()
	clearDatabase()
	a.DB.Disconnect(context.Background())
	os.Exit(code)
}

func TestCreatePlanetAndRetrieveIt(t *testing.T) {
	expectedPlanet := tatooine

	response := sendRequest("POST", "/planet/", tatooineBytes)
	if !checkResponseCode(t, http.StatusCreated, response.Code) {
		return
	}

	// comparar o planeta criado
	res := TestSinglePlanetResponse{}
	if !parseReponse(t, response, &res) {
		return
	}

	if !comparePlanet(t, expectedPlanet, res.Planet) {
		return
	}

	// pegar o planeta criado e comparar de novo
	response = sendRequest("GET", "/planet/"+res.Planet.ID, nil)
	if !checkResponseCode(t, http.StatusOK, response.Code) {
		return
	}

	if !parseReponse(t, response, &res) {
		return
	}

	if !comparePlanet(t, expectedPlanet, res.Planet) {
		return
	}
}

func TestTatooineFilmsAppearedIn(t *testing.T) {
	response := sendRequest("POST", "/planet/", tatooineBytes)
	if !checkResponseCode(t, http.StatusCreated, response.Code) {
		return
	}

	res := TestSinglePlanetResponse{}
	if !parseReponse(t, response, &res) {
		return
	}

	// tatooine apareceu em 5 filmes, mas se lançar mais filmes há a chance deste numero aumentar
	if res.Planet.FilmsAppearedIn < tatooine.FilmsAppearedIn {
		t.Errorf("Expected five or more films, but got %d.", res.Planet.FilmsAppearedIn)
	}
}

func TestBadRequestMissingFieldsCreatePlanet(t *testing.T) {
	response := sendRequest("POST", "/planet/", []byte(`{}`))
	if !checkResponseCode(t, http.StatusBadRequest, response.Code) {
		return
	}

	res := common.Error{}
	if !parseReponse(t, response, &res) {
		return
	}

	if res.Errors["name"] == "Name field is empty or missing." &&
		res.Errors["climate"] == "Climate field is empty or missing." &&
		res.Errors["terrain"] == "Terrain field is empty or missing." {
		return
	}

	t.Errorf("Server response listed missing fields incorrectly.")
}

func TestBadRequestJsonSyntaxCreatePlanet(t *testing.T) {
	response := sendRequest("POST", "/planet/", []byte(`INVALID JSON`))
	if !checkResponseCode(t, http.StatusBadRequest, response.Code) {
		return
	}

	res := common.Error{}
	if !parseReponse(t, response, &res) {
		return
	}

	if !(res.Message == common.EMINVALID && len(res.Detail) > 0) {
		t.Errorf("Server did not tell that there was something wrong with the request body.")
	}
}

func TestGetNonExistingPlanet(t *testing.T) {
	clearDatabase()
	response := sendRequest("GET", "/planet/6016c8a5e18d9b3786d7eaf4", nil)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetMatchedPlanets(t *testing.T) {
	clearDatabase()
	expectedPlanet := tatooine

	response := sendRequest("POST", "/planet/", tatooineBytes)
	if !checkResponseCode(t, http.StatusCreated, response.Code) {
		return
	}

	response = sendRequest("POST", "/planet/", []byte(`{"name":"Alderaan", "climate": "temperate", "terrain": "grasslands, mountains"}`))
	if !checkResponseCode(t, http.StatusCreated, response.Code) {
		return
	}

	response = sendRequest("GET", "/planet/?search=tatoo", nil)
	if !checkResponseCode(t, http.StatusOK, response.Code) {
		return
	}

	res := TestMatchedPlanetResponse{}
	if !parseReponse(t, response, &res) {
		return
	}

	if len(res.Results) != 1 {
		t.Errorf("Expected one single result, but got %d.", len(res.Results))
		return
	}

	if !comparePlanet(t, expectedPlanet, res.Results[0]) {
		return
	}
}

func TestGetAllPlanets(t *testing.T) {
	clearDatabase()

	response := sendRequest("POST", "/planet/", tatooineBytes)
	if !checkResponseCode(t, http.StatusCreated, response.Code) {
		return
	}

	response = sendRequest("POST", "/planet/", []byte(`{"name":"Alderaan", "climate": "temperate", "terrain": "grasslands, mountains"}`))
	if !checkResponseCode(t, http.StatusCreated, response.Code) {
		return
	}

	response = sendRequest("GET", "/planet/", nil)
	if !checkResponseCode(t, http.StatusOK, response.Code) {
		return
	}

	res := TestMultiplePlanetsResponse{}
	if !parseReponse(t, response, &res) {
		return
	}

	if len(res.Planets) != 2 {
		t.Errorf("Expected two planets, but got %d.", len(res.Planets))
		return
	}
}

func TestGetAllPlanetsAndCompare(t *testing.T) {
	clearDatabase()
	expectedPlanet := tatooine

	response := sendRequest("POST", "/planet/", tatooineBytes)
	if !checkResponseCode(t, http.StatusCreated, response.Code) {
		return
	}

	response = sendRequest("GET", "/planet/", nil)
	if !checkResponseCode(t, http.StatusOK, response.Code) {
		return
	}

	res := TestMultiplePlanetsResponse{}
	if !parseReponse(t, response, &res) {
		return
	}

	if len(res.Planets) != 1 {
		t.Errorf("Expected one planet, but got %d.", len(res.Planets))
		return
	}

	if !comparePlanet(t, expectedPlanet, res.Planets[0]) {
		return
	}
}

func TestDeletePlanet(t *testing.T) {
	clearDatabase()

	response := sendRequest("POST", "/planet/", []byte(`{"name":"Tatooine", "climate": "temperate", "terrain": "grasslands, mountains"}`))
	if !checkResponseCode(t, http.StatusCreated, response.Code) {
		return
	}

	res := TestSinglePlanetResponse{}
	if !parseReponse(t, response, &res) {
		return
	}

	response = sendRequest("DELETE", "/planet/"+res.Planet.ID, nil)
	if !checkResponseCode(t, http.StatusOK, response.Code) {
		return
	}

	response = sendRequest("GET", "/planet/"+res.Planet.ID, nil)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestDeleteNonExistingPlanet(t *testing.T) {
	clearDatabase()

	response := sendRequest("DELETE", "/planet/6016c8a5e18d9b3786d7eaf4", nil)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func clearDatabase() {
	if _, err := a.DB.Database(databaseName).Collection("planets").DeleteMany(context.TODO(), bson.D{}); err != nil {
		panic(err)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func sendRequest(method string, url string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	return executeRequest(req)
}

func comparePlanet(t *testing.T, expected, got TestPlanet) bool {
	if expected.Name != got.Name ||
		expected.Climate != got.Climate ||
		expected.Terrain != got.Terrain ||
		expected.FilmsAppearedIn > got.FilmsAppearedIn {
		t.Errorf(
			"Expected a planet like\n{name: '%s'; climate: '%s'; terrain: '%s'; filmsAppearedIn: %d or greater},\nbut got {name: '%s'; climate: '%s'; terrain: '%s'; filmsAppearedIn: %d}.",
			expected.Name, expected.Climate, expected.Terrain, expected.FilmsAppearedIn,
			got.Name, got.Climate, got.Terrain, got.FilmsAppearedIn,
		)
		return false
	}

	return true
}

func parseReponse(t *testing.T, response *httptest.ResponseRecorder, v interface{}) bool {
	if err := json.Unmarshal(response.Body.Bytes(), v); err != nil {
		t.Errorf("Could not parse response: %s", err.Error())
		return false
	}

	return true
}

func checkResponseCode(t *testing.T, expected int, actual int) bool {
	if expected != actual {
		t.Errorf("Expected response code %d, but got %d.", expected, actual)
		return false
	}

	return true
}
