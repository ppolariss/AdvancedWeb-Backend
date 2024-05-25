package message

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"src/config"

	//"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/websocket/v2"
	client "github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	. "src/models"
	"time"
)

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

func ListRecords(c *fiber.Ctx) (err error) {
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

func MossChat(c *websocket.Conn) {
	var err error
	defer func() {
		if err != nil {
			//Logger.Error(
			//	"client websocket return with error",
			//	zap.Error(err),
			//)
			response := AIResponse{Status: -1, Output: err.Error()}
			//if httpError, ok := err.(*HttpError); ok {
			//	response.StatusCode = httpError.Code
			//}
			err = c.WriteJSON(response)
			if err != nil {
				log.Println("write err error: ", err)
			}
			_ = c.Close()
		}
	}()
	_, requestMess, err := c.ReadMessage()
	if err != nil {
		log.Println("Error reading message from client:", err)
		return
	}
	requestMessage := string(requestMess)
	//fmt.Println("Received message from client:", requestMessage)
	log.Println("Received message from client:", requestMessage)

	// 连接到目标 WebSocket 服务器
	mossConn, _, err := client.DefaultDialer.Dial(config.Config.MossUrl, http.Header{
		"MOSS_API_KEY": []string{config.Config.MossApiKey},
	})
	if err != nil {
		log.Fatal(config.Config.MossUrl+" Error connecting to WebSocket server:", err)
	}
	defer func(conn *client.Conn) {
		err = conn.Close()
		if err != nil {
			fmt.Println(err)
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
				log.Println("Error reading message:", err)
				return
			}
			fmt.Printf("Received response from server: %s\n", responseMessage)
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
				return
			}
		}
	}()

	request := map[string]string{"request": requestMessage}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		log.Fatal("Error marshalling request:", err)
	}

	err = mossConn.WriteMessage(client.TextMessage, requestBytes)
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
	fmt.Printf("Sent message to server: %s\n", requestMessage)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Received interrupt signal, shutting down...")

			// Cleanly close the WebSocket connection by sending a close message
			err := mossConn.WriteMessage(client.CloseMessage, client.FormatCloseMessage(client.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error sending close message:", err)
				return
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}

	//log.Println("Server started on :8080")
	//log.Fatal(app.Listen(":8090"))
}
