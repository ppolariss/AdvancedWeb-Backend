package message

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"src/config"
	"src/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/opentreehole/go-common"

	//"github.com/gofiber/contrib/websocket"
	"os"
	"os/signal"
	. "src/models"
	"time"

	"github.com/gofiber/websocket/v2"
	client "github.com/gorilla/websocket"
)

// GetChat @GetChat
// @Router /api/chats/{id} [get]
// @Summary Get chat by ID
// @Description Get chat by ID
// @Tags Chat
// @Accept json
// @Produce json
// @Param id path int true "Chat ID"
// @Success 200 {object} Chat
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @Failure 404 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func GetChat(c *fiber.Ctx) (err error) {
	_, err = GetGeneralUser(c)
	if err != nil {
		return err
	}
	chatID, err := c.ParamsInt("id")
	if err != nil {
		return
	}
	var chat Chat
	err = DB.Preload("Record").Take(&chat, chatID).Error
	if err != nil {
		return
	}
	return c.JSON(chat)
}

// ListChats @ListChats
// @Router /api/chats [get]
// @Summary list my chats
// @Description list my chats
// @Tags Chat
// @Accept json
// @Produce json
// @Success 200 {object} Chats
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func ListChats(c *fiber.Ctx) (err error) {
	_, err = GetGeneralUser(c)
	if err != nil {
		return err
	}
	var chats Chats
	err = DB.Find(&chats).Error
	if err != nil {
		return
	}
	return c.JSON(chats)
}

// DeleteChat @DeleteChat
// @Router /api/chats/{id} [delete]
// @Summary Delete chat by ID
// @Description Delete chat by ID
// @Tags Chat
// @Accept json
// @Produce json
// @Param id path int true "Chat ID"
// @Success 204 {object} nil
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @Failure 404 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func DeleteChat(c *fiber.Ctx) (err error) {
	_, err = GetGeneralUser(c)
	if err != nil {
		return err
	}
	chatID, err := c.ParamsInt("id")
	if err != nil {
		return
	}
	err = DB.Delete(&Chat{}, chatID).Error
	if err != nil {
		return
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ListChatRecords @ListChatRecords
// @Router /api/chats/{id}/records [get]
// @Summary list records by chat ID
// @Description list records by chat ID
// @Tags ChatRecord
// @Accept json
// @Produce json
// @Param id path int true "Chat ID"
// @Success 200 {object} Records
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func ListChatRecords(c *fiber.Ctx) (err error) {
	// I don't want to do Authorization here
	_, err = GetGeneralUser(c)
	if err != nil {
		return err
	}
	chatID, err := c.ParamsInt("id")
	if err != nil {
		return
	}
	var records Records
	err = DB.Find(&records, "chat_id = ?", chatID).Error
	//err = DB.Transaction(func(tx *gorm.DB) (err error) {
	//	err = tx.Take(&chat, chatID).Error
	//	if err != nil {
	//		return
	//	}
	//	err = tx.Model(&Record{}).Where("chat_id = ? AND user_id = ?", chatID, tmpUser.ID).Find(&chat.Records).Error
	//})
	if err != nil {
		return
	}
	return c.JSON(records)
}

// ListMyChatRecords @ListMyChatRecords
// @Router /api/chats/{id}/records/me [get]
// @Summary list my records in a chat
// @Description list my records in a chat
// @Tags ChatRecord
// @Accept json
// @Produce json
// @Success 200 {object} Records
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func ListMyChatRecords(c *fiber.Ctx) (err error) {
	user, err := GetGeneralUser(c)
	if err != nil {
		return err
	}
	var records Records
	err = DB.Find(&records, "user_id = ?", user.ID).Error
	if err != nil {
		return
	}
	return c.JSON(records)
}

func addRecord(addRecordsRequest AddRecordsRequest) (err error) {
	record := Record{
		CreatedAt: addRecordsRequest.CreatedAt,
		UserID:    addRecordsRequest.UserID,
		RoomID:    addRecordsRequest.RoomID,
		Type:      addRecordsRequest.Type,
		ToID:      addRecordsRequest.ToID,
		Message:   addRecordsRequest.Message,
	}
	if err = DB.Create(&record).Error; err != nil {
		utils.Logger.Error(err.Error())
	}
	return
}

// AddRecords @AddRecords
// @Router /api/records [post]
// @Summary Add records
// @Description Add records
// @Tags ChatRecord
// @Accept json
// @Produce json
// @Param json body AddRecordsRequest true "json"
// @Success 200 {object} nil
// @Failure 400 {object} common.HttpError
// @Failure 403 {object} common.HttpError
func AddRecords(c *fiber.Ctx) (err error) {
	// private
	var addRecordsRequest AddRecordsRequest
	if err = c.BodyParser(&addRecordsRequest); err != nil {
		return common.BadRequest("Invalid request body")
	}
	return addRecord(addRecordsRequest)
}

// ListMyRecords @ListMyRecords
// @Router /api/records [get]
// @Summary list my records
// @Description list my records
// @Tags ChatRecord
// @Accept json
// @Produce json
// @Success 200 {object} Records
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func ListMyRecords(c *fiber.Ctx) (err error) {
	user, err := GetGeneralUser(c)
	if err != nil {
		return err
	}
	var records Records
	err = DB.Find(&records, "user_id = ?", user.ID).Error
	if err != nil {
		return
	}
	var toRecords Records
	err = DB.Find(&toRecords, "to_id = ?", user.ID).Error
	if err != nil {
		return
	}
	records = append(records, toRecords...)
	sort.Slice(records, func(i, j int) bool {
		return records[i].ID > records[j].ID
	})
	return c.JSON(records)
}

//func Infer(c *websocket.Conn) (err error) {
//	var (
//		chatID int
//		user   *User
//		chat   Chat
//	)
//	if chatID, err = strconv.Atoi(c.Params("id")); err != nil {
//		return common.BadRequest("invalid chat_id")
//	}
//
//	if user, err = LoadUserFromWs(c); err != nil {
//		return
//	}
//
//}

// MossChat @MossChat
// @Router /api/ws/moss [get]
// @Summary Moss Chat
// @Description Moss Chat
// @Tags Chat
// @Accept json
// @Produce json
// @Success 200 {object} AIResponse
// @Failure 400 {object} common.HttpError
// @Failure 401 {object} common.HttpError
// @Failure 403 {object} common.HttpError
// @Failure 404 {object} common.HttpError
// @param Authorization header string true "Bearer和token空格拼接"
func MossChat(c *websocket.Conn) {
	var err error
	userID, err := LoadUserFromWs(c)
	if err != nil {
		_ = c.WriteJSON(common.Unauthorized("Unauthorized"))
		return
	}
	defer func() {
		if err != nil {
			utils.Logger.Error(
				"client websocket return with error",
				//zap.Error(err),
			)
			response := AIResponse{Status: -1, Output: err.Error()}
			//if httpError, ok := err.(*HttpError); ok {
			//	response.StatusCode = httpError.Code
			//}

			if err = c.WriteJSON(response); err != nil {
				utils.Logger.Error("write err error: ", err)
			}
			if err = c.Close(); err != nil {
				utils.Logger.Error("close err error: ", err)
			}
		}
	}()
	_, requestMess, err := c.ReadMessage()
	if err != nil {
		utils.Logger.Error("Error reading message from client:", err)
		return
	}
	requestMessage := string(requestMess)
	//log.Println("Received message from client:", requestMessage)
	utils.Logger.Info("Received message from client:" + requestMessage)
	log.Println("Received message from client:", requestMessage)
	go addRecord(AddRecordsRequest{
		CreatedAt: time.Now(),
		UserID:    userID,
		RoomID:    "moss",
		Type:      "toMoss",
		ToID:      0,
		Message:   requestMessage,
	})

	// 连接到目标 WebSocket 服务器
	mossConn, _, err := client.DefaultDialer.Dial(config.Config.MossUrl, http.Header{
		"MOSS_API_KEY": []string{config.Config.MossApiKey},
	})
	if err != nil {
		utils.Logger.Error(config.Config.MossUrl+" Error connecting to WebSocket server:", err)
	}
	defer func(conn *client.Conn) {
		if err = conn.Close(); err != nil {
			utils.Logger.Error("Error closing connection:", err)
		}
	}(mossConn)

	// Handle interrupt signal to gracefully shut down the client
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, responseMessage, err := mossConn.ReadMessage()
			if err != nil {
				utils.Logger.Error("Error reading message:", err)
				return
			}
			log.Printf("Received response from server: %s\n", responseMessage)
			var mossResponse MossResponse
			err = json.Unmarshal(responseMessage, &mossResponse)
			if err != nil {
				continue
			}
			if mossResponse.Status == 0 {
				err = c.WriteMessage(websocket.TextMessage, []byte(mossResponse.Output))
				if err != nil {
					continue
				}
				go addRecord(AddRecordsRequest{
					CreatedAt: time.Now(),
					UserID:    0,
					RoomID:    "moss",
					Type:      "fromMoss",
					ToID:      userID,
					Message:   mossResponse.Output,
				})
				return
			}
		}
	}()

	request := map[string]string{"request": requestMessage}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		utils.Logger.Error("Error marshalling request:", err)
	}

	err = mossConn.WriteMessage(client.TextMessage, requestBytes)
	if err != nil {
		utils.Logger.Error("Error sending message:", err)
		return
	}
	log.Printf("Sent message to server: %s\n", requestMessage)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Received interrupt signal, shutting down...")

			// Cleanly close the WebSocket connection by sending a close message
			err = mossConn.WriteMessage(client.CloseMessage, client.FormatCloseMessage(client.CloseNormalClosure, ""))
			if err != nil {
				utils.Logger.Error("Error sending close message:", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

}

func VideoChat(c *websocket.Conn) {
	// TODO: userID and dbLog
	var err error
	defer func() {
		mu.Lock()
		delete(clients, c)
		mu.Unlock()
		err = c.Close()
		if err != nil {
			utils.Logger.Error(err.Error())
		}
	}()

	mu.Lock()
	clients[c] = true
	log.Println("clients", clients)
	log.Println("length of clients", len(clients))
	mu.Unlock()

	log.Println("connection")

	for {
		_, requestMessage, err := c.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		log.Println(string(requestMessage))
		//var msg VideoMessage
		//if err := c.ReadJSON(&msg);
		//log.Println("message")
		//log.Printf("Received message => %s", msg.Data)
		log.Printf("Clients count => %d", len(clients))
		broadcast <- struct {
			C       *websocket.Conn
			Message []byte
		}{C: c, Message: requestMessage}
	}
}

func HandleMessages() {
	var err error
	for {
		msg := <-broadcast
		log.Println(msg.C)
		log.Println(string(msg.Message))
		mu.Lock()
		// afterSelf := false
		for c := range clients {
			if c == msg.C {
				// afterSelf = true
				continue
			}
			if peers[msg.C] == nil {
				peers[msg.C] = c
				peers[c] = msg.C
			} else if peers[msg.C] != c {
				continue
			}
			// if !afterSelf {
			// 	continue
			// }
			err = c.WriteMessage(websocket.TextMessage, msg.Message)
			if err != nil {
				log.Println(websocket.IsCloseError(err, client.CloseNormalClosure))
				log.Println(websocket.IsUnexpectedCloseError(err, client.CloseNormalClosure))
				utils.Logger.Error("error: %v", err)
				err = c.Close()
				if err != nil {
					utils.Logger.Error(err.Error())
				}
				delete(clients, c)
			}
		}
		mu.Unlock()
	}
}

//app.Get("/ws", func(c *fiber.Ctx) error {
//	if websocket.IsWebSocketUpgrade(c) {
//		c.Locals("allowed", true)
//		return c.Next()
//	}
//	return fiber.ErrUpgradeRequired
//})
