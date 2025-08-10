package configs_test

import (
	"os"
	"testing"
	"time"

	"github.com/create-go-app/fiber-go-template/pkg/configs"
	"github.com/stretchr/testify/assert"
)

func TestFiberConfig_InvalidReadTimeout(t *testing.T) {
	os.Setenv("SERVER_READ_TIMEOUT", "not-a-number")
	cfg := configs.FiberConfig()
	assert.Equal(t, 0*time.Second, cfg.ReadTimeout, "Expected ReadTimeout to be 0 seconds when invalid number")
}

func TestFiberConfig_InvalidWriteTimeout(t *testing.T) {
	os.Setenv("SERVER_WRITE_TIMEOUT", "not-a-number")
	cfg := configs.FiberConfig()
	assert.Equal(t, 0*time.Second, cfg.WriteTimeout, "Expected WriteTimeout to be 0 seconds when invalid number")
}
