package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/create-go-app/fiber-go-template/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func TestFiberMiddleware(t *testing.T) {
	app := fiber.New()
	middleware.FiberMiddleware(app)

	// Add a simple route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Create a request
	req := httptest.NewRequest("GET", "/", http.NoBody)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("failed to test app: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}
}
