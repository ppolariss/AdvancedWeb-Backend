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
			//io.Emit("offline", map[string]interface{}{
			//	"socketid": id,
			//})
			//thisSocket, ok := io.Sockets().Sockets().Load(id)
			//if !ok {
			//	fmt.Println("socket not found")
			//	return
			//}
			if roomID == "" {
				return
			}
			err = io.Of("/room", nil).To(socket.Room(roomID)).Emit("offline", map[string]interface{}{
				"id": id,
			})
			if err != nil {
				utils.Logger.Error(err.Error())
				return
			}
			////fmt.Println(client.Rooms().Len())
			//for _, room := range client.Rooms().Keys() {
			//	fmt.Println(room)
			//	err = io.Of(room, nil).Emit("offline", map[string]interface{}{
			//		//err = client.Broadcast().Emit("offline", map[string]interface{}{
			//		"id": client.Id(),
			//		// "action":   "disconnect",
			//	})
			//}
			io.Of("/room", nil).Sockets().Delete(id)
			removeIdx := -1
			for idx, thisID := range GlobalRooms[roomID] {
				if thisID == string(id) {
					removeIdx = idx
				}
			}
			if removeIdx > -1 {
				GlobalRooms[roomID] = append(GlobalRooms[roomID][:removeIdx], GlobalRooms[roomID][removeIdx+1:]...)
			}
			// io.Sockets().Sockets().Delete(id)
		})
		if err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		err = client.On("disconnection", func(_ ...any) {
			fmt.Println("Info: client" + id + " disconnection")
			if roomID == "" {
				return
			}
			err = io.Of("/room", nil).To(socket.Room(roomID)).Emit("offline", map[string]interface{}{
				"id": id,
			})
			if err != nil {
				utils.Logger.Error(err.Error())
				return
			}
			io.Of("/room", nil).Sockets().Delete(id)
			// err = io.Of(socket.Room(roomID), nil).Emit("offline", map[string]interface{}{
			// 	"id": id,
			// })
			// io.Sockets().Sockets().Delete(id)
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
			if len(GlobalRooms[roomID]) == 0 {
				GlobalRooms[roomID] = make([]string, 0)
			}
			GlobalRooms[roomID] = append(GlobalRooms[roomID], string(id))
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
				err = io.Of("/room", nil).To(socket.Room(roomID)).Emit("message", chat)
				// err = io.Of(socket.Room(chat.ID), nil).Emit("message", chat)
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
	})
	return
}
