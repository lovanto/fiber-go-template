package cache

import (
	"errors"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockURLBuilder(service string) (string, error) {
	if os.Getenv("REDIS_URL") == "" {
		return "", errors.New("REDIS_URL is not set")
	}
	return os.Getenv("REDIS_URL"), nil
}

func errorURLBuilder(service string) (string, error) {
	return "", errors.New("connection error")
}

func setEnv(vars map[string]string) {
	for k, v := range vars {
		os.Setenv(k, v)
	}
}

func unsetEnv(keys ...string) {
	for _, k := range keys {
		os.Unsetenv(k)
	}
}

func runConnectionTest(t *testing.T, builder func(string) (string, error), expectError, expectNilClient bool) {
	t.Helper()

	var client *redis.Client
	var err error

	if builder == nil {
		client, err = RedisConnection()
	} else {
		client, err = NewRedisConnection(builder)
	}

	if expectError {
		assert.Error(t, err, "Expected error but got none")
	} else {
		require.NoError(t, err, "Unexpected error: %v", err)
	}

	if expectNilClient {
		assert.Nil(t, client, "Expected client to be nil")
		return
	}

	require.NotNil(t, client, "Client should not be nil")
	assert.IsType(t, &redis.Client{}, client)

	if client != nil {
		err := client.Close()
		assert.NoError(t, err, "Failed to close Redis client")
	}
}

func TestRedisConnection(t *testing.T) {
	t.Run("default RedisConnection", func(t *testing.T) {
		setEnv(map[string]string{
			"REDIS_URL":       "redis://localhost:6379",
			"REDIS_DB_NUMBER": "0",
			"REDIS_PASSWORD":  "",
		})
		defer unsetEnv("REDIS_URL", "REDIS_DB_NUMBER", "REDIS_PASSWORD")

		runConnectionTest(t, nil, false, false)
	})

	tests := []struct {
		name            string
		setup           func()
		urlBuilder      func(string) (string, error)
		expectError     bool
		expectNilClient bool
		cleanup         func()
	}{
		{
			name: "successful connection with default builder",
			setup: func() {
				setEnv(map[string]string{
					"REDIS_URL":       "redis://localhost:6379",
					"REDIS_DB_NUMBER": "0",
					"REDIS_PASSWORD":  "",
				})
			},
			urlBuilder:      mockURLBuilder,
			expectError:     false,
			expectNilClient: false,
			cleanup: func() {
				unsetEnv("REDIS_URL", "REDIS_DB_NUMBER", "REDIS_PASSWORD")
			},
		},
		{
			name: "successful connection with mock builder",
			setup: func() {
				setEnv(map[string]string{
					"REDIS_URL":       "redis://mock:6379",
					"REDIS_DB_NUMBER": "1",
				})
			},
			urlBuilder:      mockURLBuilder,
			expectError:     false,
			expectNilClient: false,
			cleanup: func() {
				unsetEnv("REDIS_URL", "REDIS_DB_NUMBER")
			},
		},
		{
			name: "invalid db number",
			setup: func() {
				setEnv(map[string]string{
					"REDIS_URL":       "redis://localhost:6379",
					"REDIS_DB_NUMBER": "invalid",
				})
			},
			urlBuilder:      mockURLBuilder,
			expectError:     false, // fallback to 0
			expectNilClient: false,
			cleanup: func() {
				unsetEnv("REDIS_URL", "REDIS_DB_NUMBER")
			},
		},
		{
			name:            "missing redis url",
			setup:           func() { unsetEnv("REDIS_URL") },
			urlBuilder:      mockURLBuilder,
			expectError:     true,
			expectNilClient: true,
			cleanup:         func() {},
		},
		{
			name:            "url builder error",
			setup:           func() {},
			urlBuilder:      errorURLBuilder,
			expectError:     true,
			expectNilClient: true,
			cleanup:         func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			runConnectionTest(t, tt.urlBuilder, tt.expectError, tt.expectNilClient)
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}
