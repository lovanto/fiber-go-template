package configs_test

import (
	"os"
	"testing"
	"time"

	"github.com/create-go-app/fiber-go-template/pkg/configs"
	"github.com/stretchr/testify/assert"
)

func TestFiberConfig(t *testing.T) {
	// Case 1: Valid integer from environment variable
	os.Setenv("SERVER_READ_TIMEOUT", "5")
	cfg := configs.FiberConfig()
	assert.Equal(t, 5*time.Second, cfg.ReadTimeout, "Expected ReadTimeout to be 5 seconds")

	// Case 2: Invalid integer from environment variable
	os.Setenv("SERVER_READ_TIMEOUT", "not-a-number")
	cfg = configs.FiberConfig()
	assert.Equal(t, 0*time.Second, cfg.ReadTimeout, "Expected ReadTimeout to be 0 seconds when invalid number")
}
