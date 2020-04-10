package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return setRoutes(r, nil)
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestGoHelloWorld(t *testing.T) {
	router := setupTestRouter()

	w := performRequest(router, "GET", "/")

	// Convert the JSON response to a map
	var response map[string]string
	json.Unmarshal([]byte(w.Body.String()), &response)

	assert.EqualValues(
		t,
		map[string]string{"data": "hello world"},
		response,
	)

}
