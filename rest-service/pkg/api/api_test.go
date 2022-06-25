package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/stackpath/backend-developer-tests/rest-service/pkg/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_SearchPeople(t *testing.T) {
	api := New()
	assert.NotNil(t, api)

	t.Run("List All", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people", nil)
		w := httptest.NewRecorder()

		api.SearchPeople(w, r, httprouter.ParamsFromContext(r.Context()))
		var result []*models.Person
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Nil(t, err)
		assert.Equal(t, models.AllPeople(), result)
	})
	t.Run("By Name", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people?first_name=Jenny&last_name=Smith", nil)
		w := httptest.NewRecorder()

		api.SearchPeople(w, r, httprouter.ParamsFromContext(r.Context()))
		var result []*models.Person
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Nil(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "+44 7700 900077", result[0].PhoneNumber)
		assert.Equal(t, "Jenny", result[0].FirstName)
		assert.Equal(t, "Smith", result[0].LastName)
	})
	t.Run("By Phone Number", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people?phone_number=%2B1%20%28800%29%20555-1414", nil)
		w := httptest.NewRecorder()

		api.SearchPeople(w, r, httprouter.ParamsFromContext(r.Context()))
		var result []*models.Person
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Nil(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "+1 (800) 555-1414", result[0].PhoneNumber)
		assert.Equal(t, "John", result[0].FirstName)
		assert.Equal(t, "Doe", result[0].LastName)
	})
	t.Run("No Results Name", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people?first_name=Bob&last_name=Smith", nil)
		w := httptest.NewRecorder()

		api.SearchPeople(w, r, httprouter.ParamsFromContext(r.Context()))
		var result []*models.Person
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		// Switch to StatusNotFound if uncommenting the 404 response in the API
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Nil(t, err)
		assert.Len(t, result, 0)
	})
	t.Run("No Results Phone Number", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people?phone_number=12345", nil)
		w := httptest.NewRecorder()

		api.SearchPeople(w, r, httprouter.ParamsFromContext(r.Context()))
		var result []*models.Person
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		// Switch to StatusNotFound if uncommenting the 404 response in the API
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Nil(t, err)
		assert.Len(t, result, 0)
	})
	t.Run("Invalid Search Name", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people?first_name", nil)
		w := httptest.NewRecorder()

		api.SearchPeople(w, r, httprouter.ParamsFromContext(r.Context()))
		var result models.Error
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Nil(t, err)
	})
	t.Run("Invalid Search Phone", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people?phone_number", nil)
		w := httptest.NewRecorder()

		api.SearchPeople(w, r, httprouter.ParamsFromContext(r.Context()))
		var result models.Error
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Nil(t, err)
	})
	t.Run("Invalid Search Bad Param", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people?test", nil)
		w := httptest.NewRecorder()

		api.SearchPeople(w, r, httprouter.ParamsFromContext(r.Context()))
		var result models.Error
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Nil(t, err)
	})
}

func TestAPI_GetPerson(t *testing.T) {
	api := New()
	assert.NotNil(t, api)

	t.Run("Exists", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people", nil)
		w := httptest.NewRecorder()

		api.GetPerson(w, r, []httprouter.Param{{Key: "id", Value: "df12ce76-767b-4bf0-bccb-816745df9e70"}})
		var result models.Person
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Nil(t, err)
		assert.Equal(t, "+44 7700 900077", result.PhoneNumber)
		assert.Equal(t, "Brian", result.FirstName)
		assert.Equal(t, "Smith", result.LastName)
	})
	t.Run("Does Not Exist", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people", nil)
		w := httptest.NewRecorder()

		api.GetPerson(w, r, []httprouter.Param{{Key: "id", Value: "df12ce76-767b-4bf0-bccb-816745df9e71"}})
		var result models.Error
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Nil(t, err)
	})
	t.Run("Invalid UUID", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/people", nil)
		w := httptest.NewRecorder()

		api.GetPerson(w, r, []httprouter.Param{{Key: "id", Value: "this-is-not-a-uuid"}})
		var result models.Error
		err := json.NewDecoder(w.Result().Body).Decode(&result)
		_ = w.Result().Body.Close()

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Nil(t, err)
	})
}
