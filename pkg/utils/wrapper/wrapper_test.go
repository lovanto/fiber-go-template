package wrapper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSuccessResponse(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return SuccessResponse(c, "ok", map[string]string{"foo": "bar"})
	})
	req := httptest.NewRequest("GET", "/", http.NoBody)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestErrorResponse(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return ErrorResponse(c, 400, "bad", assert.AnError)
	})
	req := httptest.NewRequest("GET", "/", http.NoBody)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
