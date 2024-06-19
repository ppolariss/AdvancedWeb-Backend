package room

import (
	"encoding/json"
	"fmt"
	"github.com/zishang520/socket.io/socket"
	"src/api"
	. "src/model"
	. "src/schemas"
	"src/utils"
	"time"
)

func Room(io *socket.Server) (err error) {
	err = io.Of("/room", nil).On("connection", func(clients ...any) {
		// only one client in this function
		if utils.AssertOnlyOne(clients) {
			//" + client.(*socket.Socket).Id() + "
			fmt.Println("Error: client connected, exceed 1")
		}
		client := clients[0].(*socket.Socket)
		id := client.Id()
		var roomID string
		var userID int
		//var chatID int

		fmt.Println("Info: client" + id + " connected")

		//io.Emit()
		err = client.Emit("online", map[string]interface{}{
			"id": id,
		})
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		err = client.On("disconnect", func(_ ...any) {
			fmt.Println("Info: client" + id + " disconnected")
			Disconnect(io, id, roomID)
		})
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		err = client.On("disconnection", func(_ ...any) {
			fmt.Println("Info: client" + id + " disconnection")
			Disconnect(io, id, roomID)
		})
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		err = client.On("init", func(requestData ...any) {
			if utils.AssertOnlyOne(requestData) {
				fmt.Println("Error: client" + id + " init, exceed 1")
			}
			requestDatum := requestData[0]
			var jsonData []byte
			var data Data
			jsonData, err = json.Marshal(requestDatum)
			if err != nil {
				utils.Logger.Error(err.Error())
				return
			}
			if err = json.Unmarshal(jsonData, &data); err != nil {
				utils.Logger.Error(err.Error())
				return
			}
			data.SocketID = string(id)
			roomID = data.RoomID
			userID = data.UserID
			var room = socket.Room(data.RoomID)
			fmt.Println("Info: client" + data.SocketID + " joined room " + data.RoomID)
			MutexRooms.Lock()
			if len(GlobalRooms[roomID]) == 0 {
				GlobalRooms[roomID] = make([]string, 0)
			}
			GlobalRooms[roomID] = append(GlobalRooms[roomID], string(id))
			MutexRooms.Unlock()
			client.Join(room)
			// for _, i := range client.Rooms().Keys() {
			// 	fmt.Println(i)
			// }
			client.SetData(data)
		})
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		err = client.On("update", func(requestData ...any) {
			if utils.AssertOnlyOne(requestData) {
				fmt.Println("Error: client" + id + " update, exceed 1")
			}
			requestDatum := requestData[0]
			var jsonData []byte
			var data Data
			jsonData, err = json.Marshal(requestDatum)
			if err != nil {
				utils.Logger.Error(err.Error())
				return
			}
			if err = json.Unmarshal(jsonData, &data); err != nil {
				fmt.Println("Error in Unmarshal:", err)
				return
			}
			if client.Data() == nil {
				return
			}
			oriData, ok := client.Data().(Data)
			if !ok {
				fmt.Println("Error: data not found, don't update")
				return
			}
			data.SocketID = string(id)
			data.UserID = oriData.UserID
			data.Colour = oriData.Colour
			data.Model = oriData.Model
			data.RoomID = oriData.RoomID
			client.SetData(data)
		})
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		err = client.On("chat", func(requestData ...any) {
			if utils.AssertOnlyOne(requestData) {
				fmt.Println("client" + id + " chat, exceed 1")
			}
			requestDatum := requestData[0]
			var jsonData []byte
			var chat Chat
			jsonData, err = json.Marshal(requestDatum)
			if err != nil {
				utils.Logger.Error(err.Error())
				return
			}
			if err = json.Unmarshal(jsonData, &chat); err != nil {
				utils.Logger.Error(err.Error())
				return
			}
			chat.ID = string(id)
			if chat.Type == "global" {
				err = client.Broadcast().Emit("message", chat)
			} else if chat.Type == "room" {
				fmt.Println("/room/" + (chat.RoomID))
				err = io.Of("/room", nil).To(socket.Room(roomID)).Emit("message", chat)
				// err = io.Of(socket.Room(chat.RoomID), nil).Emit("message", chat)
			} else if chat.Type == "private" {
				// TODO
				// err = io.Of("/room", nil).To(socket.Room(roomID)).Emit("message", chat)
				err = io.Of("/room", nil).To(socket.Room(chat.ToID)).Emit("message", chat)
			}
			if err != nil {
				utils.Logger.Error(err.Error())
				return
			}
			//if chatID == 0 {
			//chatID = AddChat(AddChatRequest{
			//	UserID: userID,
			//	RoomID: roomID,
			//})
			//}

			go api.AddRecord(AddRecordRequest{
				UserID:    userID,
				RoomID:    roomID,
				Type:      chat.Type,
				ToID:      0,
				Message:   chat.Message,
				CreatedAt: time.Now(),
			})

		})
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		err = client.On("event", func(requestData ...any) {
			if utils.AssertOnlyOne(requestData) {
				// utils.Logger.Error("client" + id + " event, exceed 1")
				fmt.Println("client" + id + " event, exceed 1")
			}
			requestDatum := requestData[0]
			var jsonData []byte
			var event Event
			jsonData, err = json.Marshal(requestDatum)
			if err != nil {
				utils.Logger.Error(err.Error())
				return
			}
			if err = json.Unmarshal(jsonData, &event); err != nil {
				utils.Logger.Error(err.Error())
				return
			}
			event.SocketID = string(id)
			fmt.Println("Info: client" + event.SocketID + " event " + event.Event + " in room " + event.RoomID)
			err = io.Of("/room", nil).To(socket.Room(roomID)).Emit("event", event)

			if err != nil {
				utils.Logger.Error(err.Error())
				return
			}
		})

		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}
	})
	return
}
