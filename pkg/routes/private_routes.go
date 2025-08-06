package routes

import (
	"github.com/create-go-app/fiber-go-template/app/controllers"
	"github.com/create-go-app/fiber-go-template/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for POST method:
	route.Post("/book", middleware.JWTProtected(), controllers.CreateBook)
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut)
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)

	// Routes for PUT method:
	route.Put("/book", middleware.JWTProtected(), controllers.UpdateBook)

	// Routes for DELETE method:
	route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook)
}
