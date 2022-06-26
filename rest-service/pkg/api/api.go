package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/stackpath/backend-developer-tests/rest-service/pkg/models"
	"log"
	"net/http"
	"time"
)

type API struct {
}

func New() *API {
	return &API{}
}

func (_ *API) RequestLogger(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		log.Printf("[%s] %s: %s\n", r.RemoteAddr, r.Method, r.URL)
		handler(w, r, ps)
	}
}

// writeJsonResponse Writes a JSON response with the specified status code
func (_ *API) writeJsonResponse(w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonEncoder := json.NewEncoder(w)
	if err := jsonEncoder.Encode(response); err != nil {
		log.Printf("Error writting response %s\n", err.Error())
	}
}

func (api *API) writeErrorResponse(w http.ResponseWriter, message string, code int) {
	api.writeJsonResponse(w, models.Error{Message: message, Timestamp: time.Now()}, code)
}
