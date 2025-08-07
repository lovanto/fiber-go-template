package jwt

import (
	"errors"
	"hash"
	"net/http"
	"net/http/httptest"
	"os"

	gojwt "github.com/golang-jwt/jwt/v5"

	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setDefaultEnv() {
	_ = os.Setenv("JWT_SECRET_KEY", "test-secret")
	_ = os.Setenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT", "15")
	_ = os.Setenv("JWT_REFRESH_KEY", "test-refresh")
	_ = os.Setenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT", "1")
}

func TestGenerateNewTokens(t *testing.T) {
	setDefaultEnv()

	tokens, err := GenerateNewTokens("123e4567-e89b-12d3-a456-426614174000", []string{"book:create", "book:update"})
	assert.NoError(t, err)
	assert.NotEmpty(t, tokens.Access)
	assert.NotEmpty(t, tokens.Refresh)
}

func TestExtractTokenMetadata(t *testing.T) {
	setDefaultEnv()

	tokens, err := GenerateNewTokens("123e4567-e89b-12d3-a456-426614174000", []string{"book:create"})
	assert.NoError(t, err)

	app := fiber.New()

	var meta *TokenMetadata

	app.Get("/", func(c *fiber.Ctx) error {
		var err error
		meta, err = ExtractTokenMetadata(c)
		if err != nil {
			return err
		}
		return nil
	})

	req := httptest.NewRequest("GET", "http://localhost/", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+tokens.Access)

	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	assert.NotNil(t, meta)
	assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", meta.UserID.String())
	assert.True(t, meta.Credentials["book:create"])
	assert.False(t, meta.Credentials["book:delete"])
	assert.Greater(t, meta.Expires, time.Now().Unix())
}

func TestParseRefreshToken_InvalidNumeric(t *testing.T) {
	_, err := ParseRefreshToken("abc.xyz")
	assert.Error(t, err)
}

func TestGenerateNewTokens_MissingSecret(t *testing.T) {
	// unset secret to trigger error path
	_ = os.Unsetenv("JWT_SECRET_KEY")
	setDefaultEnv() // sets other vars but overridden secret removed
	_ = os.Unsetenv("JWT_SECRET_KEY")

	_, err := GenerateNewTokens("id", nil)
	assert.Error(t, err)
}

func TestGenerateNewTokens_SignError(t *testing.T) {
	setDefaultEnv()
	// override signTokenFunc to force error
	orig := signTokenFunc
	signTokenFunc = func(token *gojwt.Token, secret []byte) (string, error) {
		return "", errors.New("sign error")
	}
	defer func() { signTokenFunc = orig }()

	_, err := GenerateNewTokens("id", nil)
	assert.Error(t, err)
}

func TestGenerateNewRefreshToken_HashError(t *testing.T) {
	setDefaultEnv()
	// override hash error
	orig := hashWriteFunc
	hashWriteFunc = func(h hash.Hash, data []byte) (int, error) {
		return 0, errors.New("hash error")
	}
	defer func() { hashWriteFunc = orig }()

	_, err := GenerateNewTokens("id", nil)
	assert.Error(t, err)
}

func TestGenerateNewTokens_MissingRefreshKey(t *testing.T) {
	setDefaultEnv()
	_ = os.Unsetenv("JWT_REFRESH_KEY")

	_, err := GenerateNewTokens("id", nil)
	assert.Error(t, err)
}

func TestExtractTokenMetadata_NoAuthHeader(t *testing.T) {
	setDefaultEnv()

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		_, err := ExtractTokenMetadata(c)
		assert.Error(t, err)
		return nil
	})

	req := httptest.NewRequest("GET", "http://localhost/", http.NoBody)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestExtractTokenMetadata_InvalidToken(t *testing.T) {
	setDefaultEnv()
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		_, err := ExtractTokenMetadata(c)
		assert.Error(t, err)
		return nil
	})

	req := httptest.NewRequest("GET", "http://localhost/", http.NoBody)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestExtractTokenMetadata_ExpiredToken(t *testing.T) {
	setDefaultEnv()
	// create expired token
	claims := gojwt.MapClaims{
		"id":  "123e4567-e89b-12d3-a456-426614174000",
		"exp": time.Now().Add(-time.Hour).Unix(),
	}
	tokenStr, _ := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	app := fiber.New()
	var meta *TokenMetadata
	app.Get("/", func(c *fiber.Ctx) error {
		meta, _ = ExtractTokenMetadata(c)
		return nil
	})

	req := httptest.NewRequest("GET", "http://localhost/", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	_, err := app.Test(req)
	assert.NoError(t, err)
	assert.Nil(t, meta)
}

func TestExtractTokenMetadata_InvalidSignature(t *testing.T) {
	setDefaultEnv()
	// sign token with different secret
	claims := gojwt.MapClaims{
		"id":  "123e4567-e89b-12d3-a456-426614174000",
		"exp": time.Now().Add(time.Minute).Unix(),
	}
	tokenStr, _ := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims).SignedString([]byte("other-secret"))

	app := fiber.New()
	var meta *TokenMetadata
	app.Get("/", func(c *fiber.Ctx) error {
		meta, _ = ExtractTokenMetadata(c)
		return nil
	})

	req := httptest.NewRequest("GET", "http://localhost/", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	_, err := app.Test(req)
	assert.NoError(t, err)
	assert.Nil(t, meta)
}

func TestExtractTokenMetadata_InvalidUserID(t *testing.T) {
	setDefaultEnv()
	// create token with invalid uuid in id claim
	claims := gojwt.MapClaims{
		"id":  "invalid",
		"exp": time.Now().Add(time.Minute).Unix(),
	}
	tokenStr, _ := gojwt.NewWithClaims(gojwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		_, err := ExtractTokenMetadata(c)
		assert.Error(t, err)
		return nil
	})

	req := httptest.NewRequest("GET", "http://localhost/", http.NoBody)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestParseRefreshToken(t *testing.T) {
	setDefaultEnv()

	tokens, err := GenerateNewTokens("123e4567-e89b-12d3-a456-426614174000", nil)
	assert.NoError(t, err)

	exp, err := ParseRefreshToken(tokens.Refresh)
	assert.NoError(t, err)

	// The expiration timestamp should be within the next few hours
	now := time.Now().Unix()
	oneHourLater := now + int64(3600)
	assert.Greater(t, exp, now)
	assert.LessOrEqual(t, exp, oneHourLater+int64(3600)) // allow small buffer
}

func TestExtractTokenMetadata_FallbackBranch(t *testing.T) {
	// Save and restore original func
	orig := verifyTokenFunc
	verifyTokenFunc = func(c *fiber.Ctx) (*gojwt.Token, error) {
		return &gojwt.Token{Valid: false, Claims: gojwt.MapClaims{}}, nil
	}
	defer func() { verifyTokenFunc = orig }()

	app := fiber.New()
	var meta *TokenMetadata
	var err error
	app.Get("/", func(c *fiber.Ctx) error {
		meta, err = ExtractTokenMetadata(c)
		return nil
	})

	req := httptest.NewRequest("GET", "http://localhost/", http.NoBody)
	_, e := app.Test(req)
	assert.NoError(t, e)
	// Function should return nil metadata and nil error (fallback branch reached)
	assert.Nil(t, meta)
	assert.Nil(t, err)
}
