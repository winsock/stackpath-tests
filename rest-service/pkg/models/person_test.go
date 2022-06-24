package models

import (
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAllPeople(t *testing.T) {
	results := AllPeople()

	assert.Len(t, results, 5)
	for _, person := range results {
		assert.NotNil(t, person)
		assert.NotEmpty(t, person.ID)
		assert.NotEmpty(t, person.FirstName)
		assert.NotEmpty(t, person.LastName)
		assert.NotEmpty(t, person.PhoneNumber)
	}
}

func TestFindPersonByID(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		person, err := FindPersonByID(uuid.Must(uuid.FromString("135af595-aa86-4bb5-a8f7-df17e6148e63")))

		assert.Nil(t, err)
		assert.NotNil(t, person)
		assert.Equal(t, "John", person.FirstName)
		assert.Equal(t, "Doe", person.LastName)
		assert.Equal(t, "+1 (800) 555-1414", person.PhoneNumber)
	})
	t.Run("Not Found", func(t *testing.T) {
		person, err := FindPersonByID(uuid.Must(uuid.FromString("135af595-aa86-4bb5-a8f7-df17e6148e64")))

		assert.NotNil(t, err)
		assert.Nil(t, person)
	})
}

func TestFindPersonByName(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		results := FindPeopleByName("Jane", "Doe")

		assert.Len(t, results, 1)
		assert.Equal(t, "Jane", results[0].FirstName)
		assert.Equal(t, "Doe", results[0].LastName)
		assert.Equal(t, "+1 (800) 555-1313", results[0].PhoneNumber)
	})
	t.Run("Not Found", func(t *testing.T) {
		results := FindPeopleByName("Jack", "Doe")

		assert.Len(t, results, 0)
	})
	t.Run("Multiple Found", func(t *testing.T) {
		results := FindPeopleByName("John", "Doe")

		assert.Len(t, results, 2)
		assert.Equal(t, "John", results[0].FirstName)
		assert.Equal(t, "Doe", results[0].LastName)
		assert.Equal(t, "+1 (800) 555-1212", results[0].PhoneNumber)
		assert.Equal(t, "John", results[1].FirstName)
		assert.Equal(t, "Doe", results[1].LastName)
		assert.Equal(t, "+1 (800) 555-1414", results[1].PhoneNumber)
	})
}

func TestFindPeopleByPhoneNumber(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		results := FindPeopleByPhoneNumber("+1 (800) 555-1212")

		assert.Len(t, results, 1)
		assert.Equal(t, "+1 (800) 555-1212", results[0].PhoneNumber)
		assert.Equal(t, "John", results[0].FirstName)
		assert.Equal(t, "Doe", results[0].LastName)
	})
	t.Run("Not Found", func(t *testing.T) {
		results := FindPeopleByPhoneNumber("+1 (800) 555-1234")
		assert.Len(t, results, 0)
	})
	t.Run("Multiple Found", func(t *testing.T) {
		results := FindPeopleByPhoneNumber("+44 7700 900077")

		assert.Len(t, results, 2)
		assert.Equal(t, "Brian", results[0].FirstName)
		assert.Equal(t, "Smith", results[0].LastName)
		assert.Equal(t, "+44 7700 900077", results[0].PhoneNumber)
		assert.Equal(t, "Jenny", results[1].FirstName)
		assert.Equal(t, "Smith", results[1].LastName)
		assert.Equal(t, "+44 7700 900077", results[1].PhoneNumber)
	})
}
