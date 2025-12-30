package unit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_api/internal/util"

	"github.com/stretchr/testify/assert"
)

func TestResponseWithSuccess(t *testing.T) {
	t.Run("should return success response with data", func(t *testing.T) {
		w := httptest.NewRecorder()
		data := map[string]string{"key": "value"}

		util.ResponseWithSuccess(w, http.StatusOK, "Success", data)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response util.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, "Success", response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("should return success response without data", func(t *testing.T) {
		w := httptest.NewRecorder()

		util.ResponseWithSuccess(w, http.StatusCreated, "Created", nil)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response util.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, "Created", response.Message)
		assert.Nil(t, response.Data)
	})
}

func TestResponseWithError(t *testing.T) {
	t.Run("should return error response with error details", func(t *testing.T) {
		w := httptest.NewRecorder()

		util.ResponseWithError(w, http.StatusBadRequest, "Bad Request", "Invalid input")

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

		var response util.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", response.Message)
		assert.Equal(t, "Invalid input", response.Error)
	})

	t.Run("should return error response for not found", func(t *testing.T) {
		w := httptest.NewRecorder()

		util.ResponseWithError(w, http.StatusNotFound, "Not Found", "Resource not found")

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response util.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, "Not Found", response.Message)
		assert.Equal(t, "Resource not found", response.Error)
	})

	t.Run("should return internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()

		util.ResponseWithError(w, http.StatusInternalServerError, "Internal Server Error", "")

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response util.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)

		assert.NoError(t, err)
		assert.Equal(t, "Internal Server Error", response.Message)
	})
}
