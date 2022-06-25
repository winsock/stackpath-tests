package api

import (
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/stackpath/backend-developer-tests/rest-service/pkg/models"
	"log"
	"net/http"
)

func (api *API) SearchPeople(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	phoneNumber := r.FormValue("phone_number")

	var results []*models.Person
	if len(r.Form) == 0 {
		results = models.AllPeople()
	} else if len(firstName) > 0 && len(lastName) > 0 {
		results = models.FindPeopleByName(firstName, lastName)
	} else if len(phoneNumber) > 0 {
		results = models.FindPeopleByPhoneNumber(phoneNumber)
	} else {
		api.writeErrorResponse(w, "Invalid search parameters provided, must provide either a first and last name or a phone number", http.StatusBadRequest)
		return
	}

	if len(results) == 0 {
		api.writeJsonResponse(w, results, http.StatusNotFound)
		return
	}

	api.writeJsonResponse(w, results, http.StatusOK)
}

func (api *API) GetPerson(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	id, err := uuid.FromString(ps.ByName("id"))
	if err != nil {
		log.Printf("Error parsing provided id, %s\n", err.Error())
		api.writeErrorResponse(w, "Invalid ID provided", http.StatusBadRequest)
		return
	}

	person, err := models.FindPersonByID(id)
	if err != nil {
		api.writeErrorResponse(w, "Person with the provided ID was not found.", http.StatusNotFound)
		return
	}

	api.writeJsonResponse(w, person, http.StatusOK)
}
