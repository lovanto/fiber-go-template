package routes

import (
	"github.com/create-go-app/fiber-go-template/app/controllers"
	"github.com/create-go-app/fiber-go-template/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func PrivateRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Post("/book", middleware.JWTProtected(), controllers.CreateBook)
	route.Post("/user/sign/out", middleware.JWTProtected(), controllers.UserSignOut)
	route.Post("/token/renew", middleware.JWTProtected(), controllers.RenewTokens)

	route.Put("/book", middleware.JWTProtected(), controllers.UpdateBook)

	route.Delete("/book", middleware.JWTProtected(), controllers.DeleteBook)
}
