package start_server

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

// This test exercises both StartServer and StartServerWithGracefulShutdown to achieve full
// statement coverage. We deliberately provide empty SERVER_HOST and SERVER_PORT so that the
// underlying fiber.Listen call returns an error immediately (ensuring the error branch is
// taken). For the graceful variant we additionally send an os.Interrupt signal so that the
// idleConnsClosed channel is closed and the function returns without hanging.
func TestStartServerFunctions(t *testing.T) {
	// unset server env so connection_url_builder returns invalid address and Listen errors.
	os.Unsetenv("SERVER_HOST")
	os.Unsetenv("SERVER_PORT")

	// Non-graceful variant – run in background without waiting (may listen forever but shouldn't block test).
	app := fiber.New()
	go StartServer(app)

	// Graceful variant – track with WaitGroup so we can wait for completion.
	app2 := fiber.New()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartServerWithGracefulShutdown(app2)
	}()

	// Give servers some time to start and enter listen loop.
	time.Sleep(50 * time.Millisecond)

	// send interrupt to self which the goroutine listens for.
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("FindProcess: %v", err)
	}
	if err := p.Signal(os.Interrupt); err != nil {
		t.Fatalf("Signal: %v", err)
	}

	wg.Wait()

	// Scenario 2: valid addr so Listen succeeds and <-idleConnsClosed line is reached.
	defer os.Unsetenv("SERVER_HOST")
	defer os.Unsetenv("SERVER_PORT")
	if err := os.Setenv("SERVER_HOST", "127.0.0.1"); err != nil {
		t.Fatalf("set env host: %v", err)
	}
	if err := os.Setenv("SERVER_PORT", "0"); err != nil {
		t.Fatalf("set env port: %v", err)
	}

	app3 := fiber.New()
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		StartServerWithGracefulShutdown(app3)
	}()

	time.Sleep(50 * time.Millisecond)
	p2, _ := os.FindProcess(os.Getpid())
	_ = p2.Signal(os.Interrupt)
	wg2.Wait()
}
