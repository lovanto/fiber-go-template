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

func sendInterrupt(t *testing.T) {
	t.Helper()
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("failed to find process: %v", err)
	}
	_ = p.Signal(os.Interrupt)
}

func TestDefaultFuncsCoverage(t *testing.T) {
	app := fiber.New()

	go func() {
		time.Sleep(50 * time.Millisecond)
		_ = shutdownFunc(app)
	}()

	err := listenAppFunc(app, "127.0.0.1:0")
	if err != nil && !strings.Contains(err.Error(), "Server is not running") {
		t.Fatalf("unexpected listen error: %v", err)
	}
}

func TestStartServer_ListenFail(t *testing.T) {
	origListen, origShutdown, origBuild := listenAppFunc, shutdownFunc, buildURLFunc
	defer func() {
		listenAppFunc = origListen
		shutdownFunc = origShutdown
		buildURLFunc = origBuild
	}()

	buildURLFunc = func(_ string) (string, error) { return "dummy", nil }
	listenAppFunc = func(_ *fiber.App, _ string) error {
		return errors.New("listen fail in StartServer")
	}

	StartServer(fiber.New())
}

func TestGraceful_ShutdownError(t *testing.T) {
	origListen, origShutdown, origBuild := listenAppFunc, shutdownFunc, buildURLFunc
	defer func() {
		listenAppFunc = origListen
		shutdownFunc = origShutdown
		buildURLFunc = origBuild
	}()

	buildURLFunc = func(_ string) (string, error) { return "dummy", nil }
	listenAppFunc = func(_ *fiber.App, _ string) error { return nil }
	shutdownFunc = func(_ *fiber.App) error { return errors.New("shutdown fail") }

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartServerWithGracefulShutdown(fiber.New())
	}()
	time.Sleep(50 * time.Millisecond)
	sendInterrupt(t)
	wg.Wait()
}

func TestGraceful_ListenError(t *testing.T) {
	origListen, origShutdown, origBuild := listenAppFunc, shutdownFunc, buildURLFunc
	defer func() {
		listenAppFunc = origListen
		shutdownFunc = origShutdown
		buildURLFunc = origBuild
	}()

	buildURLFunc = func(_ string) (string, error) { return "dummy", nil }
	listenAppFunc = func(_ *fiber.App, _ string) error {
		return errors.New("listen fail in graceful")
	}
	shutdownFunc = func(_ *fiber.App) error { return nil }

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartServerWithGracefulShutdown(fiber.New())
	}()
	time.Sleep(50 * time.Millisecond)
	sendInterrupt(t)
	wg.Wait()
}

func TestGraceful_SuccessPath(t *testing.T) {
	origListen, origShutdown, origBuild := listenAppFunc, shutdownFunc, buildURLFunc
	defer func() {
		listenAppFunc = origListen
		shutdownFunc = origShutdown
		buildURLFunc = origBuild
	}()

	buildURLFunc = func(_ string) (string, error) { return "dummy", nil }
	listenAppFunc = func(_ *fiber.App, _ string) error {
		go func() {
			time.Sleep(20 * time.Millisecond)
			_ = shutdownFunc(fiber.New())
		}()
		return nil
	}
	shutdownFunc = func(_ *fiber.App) error { return nil }

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartServerWithGracefulShutdown(fiber.New())
	}()
	time.Sleep(50 * time.Millisecond)
	sendInterrupt(t)
	wg.Wait()
}
