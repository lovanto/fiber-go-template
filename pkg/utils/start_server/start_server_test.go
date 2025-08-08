package start_server

import (
	"errors"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestDefaultFuncsCoverage(t *testing.T) {
	app := fiber.New()

	// Use a goroutine to shutdown the server immediately after starting
	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = shutdownFunc(app) // Calls the real a.Shutdown()
	}()

	// This will call the real a.Listen() from listenAppFunc with ephemeral port
	err := listenAppFunc(app, "127.0.0.1:0")
	if err != nil && !strings.Contains(err.Error(), "Server is not running") {
		t.Fatalf("unexpected listen error: %v", err)
	}
}

func TestStartServer_AllBranches(t *testing.T) {
	// Backup originals
	origListen := listenAppFunc
	origShutdown := shutdownFunc
	origBuild := buildURLFunc
	defer func() {
		listenAppFunc = origListen
		shutdownFunc = origShutdown
		buildURLFunc = origBuild
	}()

	// Always return a dummy URL
	buildURLFunc = func(_ string) (string, error) { return "dummy", nil }

	// ===== Branch 1: StartServer listen error =====
	listenAppFunc = func(_ *fiber.App, _ string) error {
		return errors.New("listen fail in StartServer")
	}
	StartServer(fiber.New()) // hit log.Printf error branch in StartServer

	// ===== Branch 2: StartServerWithGracefulShutdown shutdown error =====
	listenAppFunc = func(_ *fiber.App, _ string) error { return nil }
	shutdownFunc = func(_ *fiber.App) error {
		return errors.New("shutdown fail")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartServerWithGracefulShutdown(fiber.New())
	}()
	time.Sleep(50 * time.Millisecond) // ensure goroutine is ready
	p, _ := os.FindProcess(os.Getpid())
	_ = p.Signal(os.Interrupt)
	wg.Wait() // hit log.Printf shutdown error branch

	// ===== Branch 3: StartServerWithGracefulShutdown listen error =====
	listenAppFunc = func(_ *fiber.App, _ string) error {
		return errors.New("listen fail in graceful")
	}
	shutdownFunc = func(_ *fiber.App) error { return nil }
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartServerWithGracefulShutdown(fiber.New())
	}()
	time.Sleep(50 * time.Millisecond)
	p2, _ := os.FindProcess(os.Getpid())
	_ = p2.Signal(os.Interrupt)
	wg.Wait() // hit log.Printf listen error branch in graceful

	// ===== Branch 4: Full success path (cover <-idleConnsClosed) =====
	listenAppFunc = func(_ *fiber.App, _ string) error {
		// simulate running until shutdown is called
		go func() {
			time.Sleep(20 * time.Millisecond)
			_ = shutdownFunc(fiber.New())
		}()
		return nil
	}
	shutdownFunc = func(_ *fiber.App) error { return nil }
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartServerWithGracefulShutdown(fiber.New())
	}()
	time.Sleep(50 * time.Millisecond)
	p3, _ := os.FindProcess(os.Getpid())
	_ = p3.Signal(os.Interrupt)
	wg.Wait() // cover success path
}
