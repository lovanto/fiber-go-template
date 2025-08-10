package jwt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/create-go-app/fiber-go-template/app/models"
	"github.com/golang-jwt/jwt/v5"
)

// signTokenFunc is an overridable function used to sign JWTs. It allows tests to simulate signing errors.
var signTokenFunc = func(token *jwt.Token, secret []byte) (string, error) {
	return token.SignedString(secret)
}

// hashWriteFunc is an overridable function used to write data into a hash.Hash. Tests can swap it to inject errors.
var hashWriteFunc = func(h hash.Hash, data []byte) (int, error) {
	return h.Write(data)
}

func GenerateNewTokens(id string, credentials []string) (*models.Tokens, error) {
	accessToken, err := generateNewAccessToken(id, credentials)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateNewRefreshToken()
	if err != nil {
		return nil, err
	}

	return &models.Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

func generateNewAccessToken(id string, credentials []string) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY is not set")
	}
	minutesCount, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))

	claims := jwt.MapClaims{}

	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()
	claims["book:create"] = false
	claims["book:update"] = false
	claims["book:delete"] = false

	for _, credential := range credentials {
		claims[credential] = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := signTokenFunc(token, []byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func generateNewRefreshToken() (string, error) {
	hasher := sha256.New()
	base := os.Getenv("JWT_REFRESH_KEY")
	if base == "" {
		return "", fmt.Errorf("JWT_REFRESH_KEY is not set")
	}
	refresh := base + time.Now().String()
	_, err := hashWriteFunc(hasher, []byte(refresh))
	if err != nil {
		return "", err
	}

	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())
	t := hex.EncodeToString(hasher.Sum(nil)) + "." + expireTime

	return t, nil
}

func ParseRefreshToken(refreshToken string) (int64, error) {
	return strconv.ParseInt(strings.Split(refreshToken, ".")[1], 0, 64)
}
