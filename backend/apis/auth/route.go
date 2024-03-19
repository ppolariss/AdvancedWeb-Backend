package auth

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app fiber.Router) {
	//app.Get("/users/me", GetCurrentUser)
	app.Post("/login", Login)
	app.Post("/register", Register)
}
