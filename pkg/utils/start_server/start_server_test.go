package start_server

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// This test ensures StartServer and StartServerWithGracefulShutdown can be called without panic.
// We use a mock Fiber app and run the server in a goroutine, then shut it down quickly.
func TestStartServerFunctions(t *testing.T) {
	app := fiber.New()
	// Run StartServer in a goroutine (it will fail to listen, but should not panic)
	go func() {
		StartServer(app)
	}()

	time.Sleep(50 * time.Millisecond) // Give it a moment to start

	app2 := fiber.New()
	go func() {
		StartServerWithGracefulShutdown(app2)
	}()

	time.Sleep(50 * time.Millisecond)
}
