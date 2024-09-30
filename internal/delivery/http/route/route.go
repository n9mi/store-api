package route

import (
	"store-api/internal/delivery/http/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Router struct {
	App             *fiber.App
	ControllerSetup *controller.ControllerSetup
}

func NewRouter(app *fiber.App, controller *controller.ControllerSetup) *Router {
	return &Router{
		App:             app,
		ControllerSetup: controller,
	}
}

func (c *Router) Setup() {
	route := c.App.Group("/api/v1")

	route.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	c.SetupAuthRoute(route)
}

func (c *Router) SetupAuthRoute(route fiber.Router) {
	auth := route.Group("/auth")
	auth.Post("/register", c.ControllerSetup.AuthController.Register)
	auth.Post("/login", c.ControllerSetup.AuthController.Login)
}
