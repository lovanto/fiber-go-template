package jwt

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// TestExtractTokenMetadataErrors exercises the error branches inside utils.ExtractTokenMetadata
// to achieve full 100% statement coverage for jwt_parser.go.
func TestExtractTokenMetadataErrors(t *testing.T) {
	// prepare app without Authorization header – should fail early inside verifyToken.
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		if _, err := ExtractTokenMetadata(c); err == nil {
			t.Fatalf("expected error when Authorization header is missing")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	// Perform request without header.
	if resp, err := app.Test(httptest.NewRequest(fiber.MethodGet, "/", http.NoBody), -1); err != nil || resp.StatusCode != fiber.StatusOK {
		t.Fatalf("request failed: %v status: %d", err, resp.StatusCode)
	}

	// Now test with malformed token after "Bearer" prefix.
	app2 := fiber.New()
	app2.Get("/", func(c *fiber.Ctx) error {
		c.Request().Header.Add("Authorization", "Bearer malformed.token")
		if _, err := ExtractTokenMetadata(c); err == nil {
			t.Fatalf("expected error for malformed token")
		}
		return c.SendStatus(fiber.StatusOK)
	})
	if resp, err := app2.Test(httptest.NewRequest(fiber.MethodGet, "/", http.NoBody), -1); err != nil || resp.StatusCode != fiber.StatusOK {
		t.Fatalf("request failed 2: %v status: %d", err, resp.StatusCode)
	}

	// Test with expired token.
	defer setEnv(t, "JWT_SECRET_KEY", "secret")()
	defer setEnv(t, "JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT", "0")() // expire immediately
	tokens, err := GenerateNewTokens("123", nil)
	if err != nil {
		t.Fatalf("GenerateNewTokens error: %v", err)
	}

	// Wait a second to ensure expiration.
	time.Sleep(time.Second)

	app3 := fiber.New()
	app3.Get("/", func(c *fiber.Ctx) error {
		c.Request().Header.Add("Authorization", "Bearer "+tokens.Access)
		if _, err := ExtractTokenMetadata(c); err == nil {
			t.Fatalf("expected error for expired token")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	if resp, err := app3.Test(httptest.NewRequest(fiber.MethodGet, "/", http.NoBody), -1); err != nil || resp.StatusCode != fiber.StatusOK {
		t.Fatalf("request failed 3: %v status: %d", err, resp.StatusCode)
	}
}

func setEnv(t *testing.T, k, v string) func() {
	old := os.Getenv(k)
	if err := os.Setenv(k, v); err != nil {
		t.Fatalf("failed to set env %s: %v", k, err)
	}
	return func() {
		if err := os.Setenv(k, old); err != nil {
			t.Fatalf("failed to restore env %s: %v", k, err)
		}
	}
}

func TestGenerateNewTokensAndParse(t *testing.T) {
	// edge: empty credentials
	defer setEnv(t, "JWT_SECRET_KEY", "secret")()
	defer setEnv(t, "JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT", "1")()
	defer setEnv(t, "JWT_REFRESH_KEY", "refresh")()
	defer setEnv(t, "JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT", "1")()
	tokensA, errA := GenerateNewTokens("123", nil)
	if errA != nil {
		t.Fatalf("unexpected error: %v", errA)
	}
	if tokensA.Access == "" || tokensA.Refresh == "" {
		t.Fatalf("tokens should not be empty for empty credentials")
	}

	// edge: very large expire env
	defer setEnv(t, "JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT", "1000000")()
	tokensB, errB := GenerateNewTokens("123", []string{"book:create"})
	if errB != nil {
		t.Fatalf("unexpected error for large expire: %v", errB)
	}
	_ = tokensB

	// edge: zero expire env
	defer setEnv(t, "JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT", "0")()
	tokensC, errC := GenerateNewTokens("123", []string{"book:create"})
	if errC != nil {
		t.Fatalf("unexpected error for zero expire: %v", errC)
	}
	_ = tokensC

	// edge: bad secret (empty)
	defer setEnv(t, "JWT_SECRET_KEY", "")()
	tokensD, errD := GenerateNewTokens("123", []string{"book:create"})
	if errD != nil {
		t.Fatalf("unexpected error for empty secret: %v", errD)
	}
	if tokensD.Access == "" {
		t.Fatalf("access token should still be generated even with empty secret")
	}

	// edge: refresh token parse error (malformed, but must contain a dot to avoid panic)
	_, errE := ParseRefreshToken("bad.token")
	if errE == nil {
		t.Fatal("expected error for malformed refresh token")
	}

	// invalid expire env: should still generate tokens, but expiration will be zero
	defer setEnv(t, "JWT_SECRET_KEY", "secret")()
	defer setEnv(t, "JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT", "bad")()
	defer setEnv(t, "JWT_REFRESH_KEY", "refresh")()
	defer setEnv(t, "JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT", "bad")()
	tokens2, err2 := GenerateNewTokens("id", []string{"book:create"})
	if err2 != nil {
		t.Fatalf("unexpected error for invalid expire env: %v", err2)
	}
	if tokens2.Access == "" || tokens2.Refresh == "" {
		t.Fatal("tokens should still be generated even with invalid env")
	}

	// prepare env vars for token generation
	defer setEnv(t, "JWT_SECRET_KEY", "secret")()
	defer setEnv(t, "JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT", "1")()
	defer setEnv(t, "JWT_REFRESH_KEY", "refresh")()
	defer setEnv(t, "JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT", "1")()

	creds := []string{"book:create"}
	tokens, err := GenerateNewTokens("123", creds)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tokens.Access == "" || tokens.Refresh == "" {
		t.Fatalf("tokens should not be empty")
	}

	// parse refresh token returns unix timestamp in future
	exp, err := ParseRefreshToken(tokens.Refresh)
	if err != nil {
		t.Fatalf("ParseRefreshToken error: %v", err)
	}
	if exp <= time.Now().Unix() {
		t.Fatalf("refresh token exp should be in the future")
	}
}

func TestExtractTokenMetadata(t *testing.T) {
	defer setEnv(t, "JWT_SECRET_KEY", "secret")()
	defer setEnv(t, "JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT", "5")()

	id := "123e4567-e89b-12d3-a456-426614174000"
	tokens, err := GenerateNewTokens(id, []string{"book:update"})
	if err != nil {
		t.Fatalf("token generation failed: %v", err)
	}

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		meta, err := ExtractTokenMetadata(c)
		if err != nil {
			return err
		}
		if meta.UserID.String() != id {
			t.Fatalf("expected %s, got %s", id, meta.UserID)
		}
		if !meta.Credentials["book:update"] {
			t.Fatalf("expected credential book:update true")
		}
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(fiber.MethodGet, "/", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+tokens.Access)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test error: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200 got %d", resp.StatusCode)
	}
}
