package password_generator

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/create-go-app/fiber-go-template/pkg/utils/wrapper"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordGenerationAndComparison(t *testing.T) {
	// edge: empty password
	hashed := GeneratePassword("")
	if hashed == "" {
		t.Fatalf("hashed password should not be empty for empty input")
	}

	// edge: invalid bcrypt cost (simulate by passing a very large password)
	longPwd := make([]byte, 10000)
	for i := range longPwd {
		longPwd[i] = 'a'
	}
	hashedErr := GeneratePassword(string(longPwd))
	if hashedErr == "" {
		t.Fatalf("should return error string for bcrypt failure")
	}

	plain := "super-secret"
	hashed2 := GeneratePassword(plain)

	if err := bcrypt.CompareHashAndPassword([]byte(hashed2), []byte(plain)); err != nil {
		t.Fatalf("expected hashed password to match original, got error %v", err)
	}

	if !ComparePasswords(hashed2, plain) {
		t.Fatalf("ComparePasswords should return true for valid match")
	}

	if ComparePasswords(hashed2, "wrong") {
		t.Fatalf("ComparePasswords should return false for invalid match")
	}
}

func TestWrapperResponses(t *testing.T) {
	app := fiber.New()

	app.Get("/ok", func(c *fiber.Ctx) error {
		return wrapper.SuccessResponse(c, "done", fiber.Map{"foo": "bar"})
	})

	app.Get("/err", func(c *fiber.Ctx) error {
		return wrapper.ErrorResponse(c, fiber.StatusBadRequest, "oops", fiber.NewError(fiber.StatusBadRequest, "detail"))
	})

	// success path
	req := httptest.NewRequest("GET", "/ok", http.NoBody)
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// error path
	req2 := httptest.NewRequest("GET", "/err", http.NoBody)
	resp2, _ := app.Test(req2, -1)
	if resp2.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp2.StatusCode)
	}
}
