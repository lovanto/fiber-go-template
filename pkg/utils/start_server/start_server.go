package start_server

import (
    "log"
    "os"
    "os/signal"

    "github.com/create-go-app/fiber-go-template/pkg/utils/connection_url_builder"
    "github.com/gofiber/fiber/v2"
)

// hooks to allow stubbing in tests
var listenAppFunc = func(a *fiber.App, addr string) error { return a.Listen(addr) }
var shutdownFunc = func(a *fiber.App) error { return a.Shutdown() }
var buildURLFunc = connection_url_builder.ConnectionURLBuilder

func StartServerWithGracefulShutdown(a *fiber.App) {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

        if err := shutdownFunc(a); err != nil {
			log.Printf("Oops... Server is not shutting down! Reason: %v", err)
		}

		close(idleConnsClosed)
	}()

	fiberConnURL, _ := buildURLFunc("fiber")

	if err := listenAppFunc(a, fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}

func StartServer(a *fiber.App) {
	fiberConnURL, _ := buildURLFunc("fiber")
	if err := listenAppFunc(a, fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
