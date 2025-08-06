package main

import (
	"os"

	"github.com/create-go-app/fiber-go-template/pkg/configs"
	"github.com/create-go-app/fiber-go-template/pkg/middleware"
	"github.com/create-go-app/fiber-go-template/pkg/routes"
	"github.com/create-go-app/fiber-go-template/pkg/utils"

	"github.com/gofiber/fiber/v2"

	_ "github.com/create-go-app/fiber-go-template/docs"

	_ "github.com/joho/godotenv/autoload"
)

// @title API
// @version 1.0
// @description This is an auto-generated API Docs.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email your@mail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config := configs.FiberConfig()
	app := fiber.New(config)
	middleware.FiberMiddleware(app)

	routes.SwaggerRoute(app)
	routes.PublicRoutes(app)
	routes.PrivateRoutes(app)
	routes.NotFoundRoute(app)

	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
