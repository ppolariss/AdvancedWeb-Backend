package apis

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"src/apis/auth"
	"src/apis/user"
	"src/models"
)

func registerRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api")
		//return c.SendString("Hello, World 👋!")
	})
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	app.Get("/docs/*", fiberSwagger.WrapHandler)
}

// RegisterRoutes registers the necessary routes to the app
func RegisterRoutes(app *fiber.App) {
	registerRoutes(app)

	groupWithoutAuthorization := app.Group("/api")
	auth.RegisterRoutes(groupWithoutAuthorization)

	group := app.Group("/api")
	group.Use(MiddlewareGetUser)
	user.RegisterRoutes(group)
}

func MiddlewareGetUser(c *fiber.Ctx) error {
	userObject, err := models.GetGeneralUser(c)
	if err != nil {
		return err
	}
	c.Locals("user", userObject)
	//if config.Config.AdminOnly {
	//	if !userObject.IsAdmin {
	//		return common.Forbidden()
	//	}
	//}
	return c.Next()
}
