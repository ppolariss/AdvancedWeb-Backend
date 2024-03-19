package bootstrap

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/opentreehole/go-common"
	"src/apis"
	"src/models"
)

func Init() (*fiber.App, error) {
	err := models.InitDB()
	if err != nil {
		return nil, err
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: common.ErrorHandler,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		//DisableStartupMessage: true,
	})
	registerMiddlewares(app)
	apis.RegisterRoutes(app)
	return app, nil
}

func registerMiddlewares(app *fiber.App) {
	//app.Use(recover.New(recover.Config{EnableStackTrace: true}))
	app.Use(common.MiddlewareGetUserID)
	//if config.Config.Mode != "bench" {
	//	app.Use(common.MiddlewareCustomLogger)
	//}
	app.Use(pprof.New())
}
