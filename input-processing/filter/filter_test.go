package filter

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestFilter_Process(t *testing.T) {
	filter := New("error", 4096)
	t.Run("No Match", func(t *testing.T) {
		source := bytes.NewReader([]byte("hello world"))
		var dest bytes.Buffer
		err := filter.Process(&dest, source)

		assert.Nil(t, err)
		assert.Equal(t, 0, dest.Len())
	})
	t.Run("Match", func(t *testing.T) {
		source := bytes.NewReader([]byte("error: hello world"))
		var dest bytes.Buffer
		err := filter.Process(&dest, source)

		assert.Nil(t, err)
		assert.Equal(t, "error: hello world", dest.String())
	})
	t.Run("Match Second Line", func(t *testing.T) {
		source := bytes.NewReader([]byte("not a match\nerror: hello"))
		var dest bytes.Buffer
		err := filter.Process(&dest, source)

		assert.Nil(t, err)
		assert.Equal(t, "error: hello", dest.String())
	})
	t.Run("First Line Only Matches", func(t *testing.T) {
		source := bytes.NewReader([]byte("error a match\nhello!\n"))
		var dest bytes.Buffer
		err := filter.Process(&dest, source)

		assert.Nil(t, err)
		assert.Equal(t, "error a match\n", dest.String())
	})
	t.Run("Multiple Matches", func(t *testing.T) {
		source := bytes.NewReader([]byte("error a match\nhello error!\n"))
		var dest bytes.Buffer
		err := filter.Process(&dest, source)

		assert.Nil(t, err)
		assert.Equal(t, "error a match\nhello error!\n", dest.String())
	})
	t.Run("Large Single Line not Matching", func(t *testing.T) {
		randomData := make([]byte, 100000000) // 100MB
		_, _ = rand.Read(randomData)
		// Base64 encoded to ensure no random newlines were added
		randomData = []byte(base64.StdEncoding.EncodeToString(randomData))
		source := bytes.NewReader(randomData)
		var dest bytes.Buffer
		err := filter.Process(&dest, source)

		assert.Nil(t, err)
		assert.Equal(t, 0, dest.Len())
	})
	t.Run("Large Single Line Matching", func(t *testing.T) {
		randomData := make([]byte, 100000000) // 100MB
		_, _ = rand.Read(randomData)
		// Base64 encoded to ensure no random newlines were added
		randomData = []byte(base64.StdEncoding.EncodeToString(randomData) + "error")
		source := bytes.NewReader(randomData)
		var dest bytes.Buffer
		err := filter.Process(&dest, source)

		assert.Nil(t, err)
		assert.Equal(t, len(randomData), dest.Len())
	})
	t.Run("Error Writing to Output", func(t *testing.T) {
		source := bytes.NewReader([]byte("error: hello world"))
		err := filter.Process(&ErrorReaderWriter{}, source)

		assert.NotNil(t, err)
	})
	t.Run("Error Reading from Input", func(t *testing.T) {
		var dest bytes.Buffer
		err := filter.Process(&dest, &ErrorReaderWriter{})

		assert.NotNil(t, err)
	})
	t.Run("Almost Match", func(t *testing.T) {
		source := bytes.NewReader([]byte("err.or: Not a match"))
		var dest bytes.Buffer
		err := filter.Process(&dest, source)

		assert.Nil(t, err)
		assert.Equal(t, 0, dest.Len())
	})
	t.Run("Almost Match", func(t *testing.T) {
		source := bytes.NewReader([]byte("err.or: Not a match"))
		var dest bytes.Buffer
		err := filter.Process(&dest, source)

		assert.Nil(t, err)
		assert.Equal(t, 0, dest.Len())
	})
}

// Used to test that we handle errors writing
type ErrorReaderWriter struct {
}

func (_ *ErrorReaderWriter) Read(_ []byte) (n int, err error) {
	return 0, errors.New("failed to read")
}

func (_ *ErrorReaderWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New("failed to write")
}
