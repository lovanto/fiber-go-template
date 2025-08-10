package routes_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/create-go-app/fiber-go-template/pkg/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestNotFoundRoute(t *testing.T) {
	app := fiber.New()

	routes.NotFoundRoute(app)

	req := httptest.NewRequest("GET", "/nonexistent", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)
	assert.Equal(t, true, body["error"])
	assert.Equal(t, "sorry, endpoint is not found", body["msg"])
}
