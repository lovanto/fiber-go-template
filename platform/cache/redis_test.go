package cache

import (
	"errors"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockURLBuilder is a mock implementation of URLBuilder for testing
func mockURLBuilder(service string) (string, error) {
	if os.Getenv("REDIS_URL") == "" {
		return "", errors.New("REDIS_URL is not set")
	}
	return os.Getenv("REDIS_URL"), nil
}

// errorURLBuilder is a mock URLBuilder that always returns an error
func errorURLBuilder(service string) (string, error) {
	return "", errors.New("connection error")
}

func TestRedisConnection(t *testing.T) {
	// First test the RedisConnection function specifically
	t.Run("test RedisConnection with default builder", func(t *testing.T) {
		// Setup environment
		os.Setenv("REDIS_URL", "redis://localhost:6379")
		os.Setenv("REDIS_DB_NUMBER", "0")
		os.Setenv("REDIS_PASSWORD", "")
		defer func() {
			os.Unsetenv("REDIS_URL")
			os.Unsetenv("REDIS_DB_NUMBER")
			os.Unsetenv("REDIS_PASSWORD")
		}()

		// Call the function under test
		client, err := RedisConnection()

		// Assertions
		require.NoError(t, err, "RedisConnection should not return an error")
		require.NotNil(t, client, "Client should not be nil")
		assert.IsType(t, &redis.Client{}, client, "Returned client has wrong type")

		// Cleanup
		if client != nil {
			err := client.Close()
			assert.NoError(t, err, "Failed to close Redis client")
		}
	})

	tests := []struct {
		name           string
		setup          func()
		urlBuilder     func(service string) (string, error)
		expectError    bool
		expectNilClient bool
		cleanup        func()
	}{
		{
			name: "successful connection with default builder",
			setup: func() {
				os.Setenv("REDIS_URL", "redis://localhost:6379")
				os.Setenv("REDIS_DB_NUMBER", "0")
				os.Setenv("REDIS_PASSWORD", "")
			},
			urlBuilder:     mockURLBuilder,
			expectError:    false,
			expectNilClient: false,
			cleanup: func() {
				os.Unsetenv("REDIS_URL")
				os.Unsetenv("REDIS_DB_NUMBER")
				os.Unsetenv("REDIS_PASSWORD")
			},
		},
		{
			name: "successful connection with mock builder",
			setup: func() {
				os.Setenv("REDIS_URL", "redis://mock:6379")
				os.Setenv("REDIS_DB_NUMBER", "1")
			},
			urlBuilder:     mockURLBuilder,
			expectError:    false,
			expectNilClient: false,
			cleanup: func() {
				os.Unsetenv("REDIS_URL")
				os.Unsetenv("REDIS_DB_NUMBER")
			},
		},
		{
			name: "invalid db number",
			setup: func() {
				os.Setenv("REDIS_URL", "redis://localhost:6379")
				os.Setenv("REDIS_DB_NUMBER", "invalid")
			},
			urlBuilder:     mockURLBuilder,
			expectError:    false, // strconv.Atoi will use 0 on error
			expectNilClient: false,
			cleanup: func() {
				os.Unsetenv("REDIS_URL")
				os.Unsetenv("REDIS_DB_NUMBER")
			},
		},
		{
			name: "missing redis url",
			setup: func() {
				os.Unsetenv("REDIS_URL")
			},
			urlBuilder:     mockURLBuilder,
			expectError:    true,
			expectNilClient: true,
			cleanup:        func() {},
		},
		{
			name: "url builder error",
			setup: func() {
				// No setup needed for this test case
			},
			urlBuilder:     errorURLBuilder,
			expectError:    true,
			expectNilClient: true,
			cleanup:        func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			if tt.setup != nil {
				tt.setup()
			}

			// Run the function under test with the specified URL builder
			var client *redis.Client
			var err error

			if tt.urlBuilder == nil {
				// Default case, use the RedisConnection function
				client, err = RedisConnection()
			} else {
				// Use the NewRedisConnection with the specified URL builder
				client, err = NewRedisConnection(tt.urlBuilder)
			}

			// Assertions
			if tt.expectError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				require.NoError(t, err, "Unexpected error: %v", err)
			}

			if tt.expectNilClient {
				assert.Nil(t, client, "Expected client to be nil")
			} else {
				assert.NotNil(t, client, "Client should not be nil")
				assert.IsType(t, &redis.Client{}, client, "Returned client has wrong type")

				// Verify client can be closed
				if client != nil {
					err := client.Close()
					assert.NoError(t, err, "Failed to close Redis client")
				}
			}

			// Cleanup
			if tt.cleanup != nil {
				tt.cleanup()
			}
		})
	}
}

// TestMain is not needed here as we're not running any global setup/teardown
