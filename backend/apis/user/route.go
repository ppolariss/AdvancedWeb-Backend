package user

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	//app.Get("/users/all", GetAllUsers)
	app.Get("/users/data", GetUserInfo)
	app.Put("/users", UpdateUser)
	app.Put("/users/password", ChangePassword)
}
