package route

import (
	"store-api/internal/delivery/http/controller"
	"store-api/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Router struct {
	App         *fiber.App
	Controllers *controller.Controllers
	Middlewares *middleware.Middlewares
}

func NewRouter(app *fiber.App, controllers *controller.Controllers, middlewares *middleware.Middlewares) *Router {
	return &Router{
		App:         app,
		Controllers: controllers,
		Middlewares: middlewares,
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
	c.SetupCustomerRoute(route)
}

func (c *Router) SetupAuthRoute(route fiber.Router) {
	auth := route.Group("/auth")
	auth.Post("/register", c.Controllers.AuthController.Register)
	auth.Post("/login", c.Controllers.AuthController.Login)
}

func (c *Router) SetupCustomerRoute(route fiber.Router) {
	customer := route.Group("/customer")
	customer.Use(c.Middlewares.AuthMiddleware)
	customer.Get("/products", c.Controllers.ProductController.GetAll)
}
