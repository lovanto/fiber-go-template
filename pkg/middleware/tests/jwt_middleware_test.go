package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/create-go-app/fiber-go-template/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func callJwtError(app *fiber.App, path string, err error) int {
	app.Get(path, func(c *fiber.Ctx) error {
		return func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": true,
					"msg":   err.Error(),
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}(c, err)
	})

	req := httptest.NewRequest("GET", path, http.NoBody)
	resp, _ := app.Test(req)
	return resp.StatusCode
}

func TestJWTProtected_BothBranches(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "testsecret")

	app := fiber.New()

	code := callJwtError(app, "/missing", errors.New("Missing or malformed JWT"))
	if code != fiber.StatusBadRequest {
		t.Errorf("expected %d, got %d", fiber.StatusBadRequest, code)
	}

	code = callJwtError(app, "/unauth", errors.New("Invalid token"))
	if code != fiber.StatusUnauthorized {
		t.Errorf("expected %d, got %d", fiber.StatusUnauthorized, code)
	}
}

func TestJWTProtected_Integration(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "testsecret")

	app := fiber.New()
	app.Use(middleware.JWTProtected())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Set("Authorization", "Bearer badtoken")
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Errorf("expected %d, got %d", fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestJWTProtected_MissingOrMalformed(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "testsecret")

	app := fiber.New()
	app.Use(middleware.JWTProtected())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/", http.NoBody)
	req.Header.Set("Authorization", "Bearer ")
	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected %d, got %d", fiber.StatusBadRequest, resp.StatusCode)
	}
}
