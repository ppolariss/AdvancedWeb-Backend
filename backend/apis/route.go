package apis

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/filesystem"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	// "net/http"
	"src/apis/auth"
	"src/apis/exam"
	"src/apis/message"
	"src/apis/user"
	"src/models"
)

func registerRoutes(app *fiber.App) {
	// app.Use(filesystem.New(filesystem.Config{
	// 	Root:         http.Dir("./html"),
	// 	Browse:       true,
	// 	Index:        "index.html",
	// 	NotFoundFile: "index.html",
	// 	MaxAge:       3600,
	// }))
	// app.Get("/*", func(c *fiber.Ctx) error {
	// 	return c.SendFile("html/index.html")
	// 	// return c.Redirect("/api")
	// 	//return c.SendString("Hello, World 👋!")
	// })
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
	message.RegisterRoutesWithoutAuthorization(groupWithoutAuthorization)

	group := app.Group("/api")
	group.Use(MiddlewareGetUser)
	user.RegisterRoutes(group)
	exam.RegisterRoutes(group)
	message.RegisterRoutes(group)
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
