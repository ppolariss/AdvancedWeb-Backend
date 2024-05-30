package message

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func RegisterRoutesWithoutAuthorization(routes fiber.Router) {
	routes.Use("/ws/*", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	routes.Get("/ws/chat", websocket.New(MossChat))
	routes.Get("/ws/video", websocket.New(VideoChat))
	routes.Post("/records", AddRecords)
	go HandleMessages()
}

func RegisterRoutes(routes fiber.Router) {
	// chat
	routes.Get("/chats/:id", GetChat)
	routes.Get("/chats", ListChats)
	//routes.Post("/chats", AddChat)
	//routes.Put("/chats/:id", ModifyChat)
	routes.Delete("/chats/:id", DeleteChat)

	// record of chat
	routes.Get("/chats/:id/records", ListChatRecords)
	routes.Get("/chats/:id/records/me", ListMyChatRecords)

	// record
	routes.Get("/records", ListMyRecords)
	//routes.Post("/ws/chats/:id/infer", websocket.New(Infer))
	//routes.Post("/chats/:id/records", AddRecord)
	//routes.Post("/chats/:id/messages", AddMessage)
	//routes.Get("/ws/chats/:id/messages", AddMessageAsync)
	//routes.Put("/messages/:id", ModifyMessage)
	//routes.Delete("/messages/:id", DeleteMessage)
	//routes.Get("/messages/:id/screenshots", GenerateMessageScreenshot)

	//routes.Static("/screenshots", "./screenshots")
}
