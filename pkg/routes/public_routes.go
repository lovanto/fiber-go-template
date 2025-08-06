package routes

import (
	"github.com/create-go-app/fiber-go-template/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App) {
	route := a.Group("/api/v1")

	route.Get("/books", controllers.GetBooks)
	route.Get("/book/:id", controllers.GetBook)

	route.Post("/user/sign/up", controllers.UserSignUp)
	route.Post("/user/sign/in", controllers.UserSignIn)
}
