package main

import (
	"encoding/json"
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/socket.io/socket"
)

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Rotation struct {
	W float64 `json:"w"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Data struct {
	Position Position `json:"position"`
	Rotation Rotation `json:"rotation"`
	SocketID string   `json:"id"`
	RoomID   string   `json:"roomID"`
	Model    string   `json:"model"`
	Colour   string   `json:"colour"`
}

//type Socket struct {
//	Data Data `json:"data"`
//}

type Chat struct {
	ID      string `json:"id"`
	RoomID  string `json:"roomID"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func assertOnlyOne(data []any) bool {
	for i := range data {
		if i != 0 {
			return true
		}
	}
	return false
}

func socketServer() (err error) {
	httpServer := types.CreateServer(nil)
	c := socket.DefaultServerOptions()
	c.SetAllowEIO3(true)
	c.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	globalRooms := make(map[string][]string)

	io := socket.NewServer(httpServer, c)
	err = io.Of("/hall", nil).On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		err = client.Emit("sendRooms", map[string]interface{}{
			"rooms": globalRooms,
		})
		if err != nil {
			fmt.Println("Error in sendRooms", err)
		}
		err = client.On("createRooms", func(_ ...any) {
			//for _, data := range roomsData {
			//	if data != nil {
			//		var jsonData []byte
			//		jsonData, err = json.Marshal(data)
			//		if err != nil {
			//			return
			//		}
			//		room := string(jsonData)
			//		if len(globalRooms[room]) == 0 {
			//			fmt.Println("Info: room " + room + " created")
			//			globalRooms[room] = make([]string, 0)
			//		}
			//	}
			//}
			globalRooms["room"+strconv.Itoa(len(globalRooms)+1)] = make([]string, 0)
			err = client.Emit("sendRooms", map[string]interface{}{
				"rooms": globalRooms,
			})
		})
		if err != nil {
			fmt.Println("Error in createRooms", err)
		}
	})
	if err != nil {
		return
	}
	err = io.Of("/room", nil).On("connection", func(clients ...any) {
		// only one client in this function
		if assertOnlyOne(clients) {
			//" + client.(*socket.Socket).Id() + "
			fmt.Println("Error: client connected, exceed 1")
		}
		client := clients[0].(*socket.Socket)
		id := client.Id()
		var roomID string
		fmt.Println("Info: client" + id + " connected")

		//io.Emit()
		err = client.Emit("online", map[string]interface{}{
			"id": id,
		})
		if err != nil {
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
			// io.Sockets().Sockets().Delete(id)
		})
		err = client.On("disconnection", func(_ ...any) {
			fmt.Println("Info: client" + id + " disconnection")
			if roomID == "" {
				return
			}
			err = io.Of("/room", nil).To(socket.Room(roomID)).Emit("offline", map[string]interface{}{
				"id": id,
			})
			io.Of("/room", nil).Sockets().Delete(id)
			// err = io.Of(socket.Room(roomID), nil).Emit("offline", map[string]interface{}{
			// 	"id": id,
			// })
			// io.Sockets().Sockets().Delete(id)
		})

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = client.On("init", func(requestData ...any) {
			if assertOnlyOne(requestData) {
				fmt.Println("Error: client" + id + " init, exceed 1")
			}
			requestDatum := requestData[0]
			var jsonData []byte
			var data Data
			jsonData, err = json.Marshal(requestDatum)
			if err != nil {
				return
			}
			if err = json.Unmarshal(jsonData, &data); err != nil {
				fmt.Println("Error:", err)
				return
			}
			data.SocketID = string(id)
			roomID = data.RoomID
			var room = socket.Room(data.RoomID)
			fmt.Println("Info: client" + data.SocketID + " joined room " + data.RoomID)
			globalRooms[roomID] = make([]string, 0)
			globalRooms[roomID] = append(globalRooms[roomID], string(id))
			client.Join(room)
			// for _, i := range client.Rooms().Keys() {
			// 	fmt.Println(i)
			// }
			client.SetData(data)
		})

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		err = client.On("update", func(requestData ...any) {
			if assertOnlyOne(requestData) {
				fmt.Println("Error: client" + id + " update, exceed 1")
			}
			requestDatum := requestData[0]
			var jsonData []byte
			var data Data
			jsonData, err = json.Marshal(requestDatum)
			if err != nil {
				return
			}
			if err = json.Unmarshal(jsonData, &data); err != nil {
				fmt.Println("Error in Ummarshal:", err)
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
			data.Colour = oriData.Colour
			data.Model = oriData.Model
			data.RoomID = oriData.RoomID
			client.SetData(data)
		})

		if err != nil {
			fmt.Println(err)
			return
		}

		err = client.On("chat", func(requestData ...any) {
			if assertOnlyOne(requestData) {
				fmt.Println("client" + id + " chat, exceed 1")
			}
			requestDatum := requestData[0]
			var jsonData []byte
			var chat Chat
			jsonData, err = json.Marshal(requestDatum)
			if err != nil {
				return
			}
			if err = json.Unmarshal(jsonData, &chat); err != nil {
				fmt.Println("Error:", err)
				return
			}
			chat.ID = string(id)
			if chat.Type == "global" {
				err = client.Broadcast().Emit("message", chat)
			} else if chat.Type == "room" {
				err = io.Of("/room", nil).To(socket.Room(chat.RoomID)).Emit("message", chat)
				// err = io.Of(socket.Room(chat.RoomID), nil).Emit("message", chat)
			} else if chat.Type == "private" {
				// TODO
				err = io.Of("/room", nil).To(socket.Room(chat.RoomID)).Emit("message", chat)
				// err = io.Of(socket.Room(chat.ID), nil).Emit("message", chat)
			}
			if err != nil {
				fmt.Println(err)
				return
			}
		})
	})
	if err != nil {
		return
	}
	httpServer.Listen(":3000", func() {
		fmt.Println("Listening on 3000")
	})

	ticker := time.NewTicker(40 * time.Millisecond)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			socketMap := io.Of("/room", nil).Sockets()
			roomMap := make(map[string][]Data)
			for _, socketID := range socketMap.Keys() {
				//fmt.Println("update" + socketID)
				thisSocket, ok := socketMap.Load(socketID)
				if !ok {
					fmt.Println("Error: socket not found, don't send update")
					continue
				}
				data, ok := thisSocket.Data().(Data)
				if !ok {
					fmt.Println("Error: data not found, don't send update")
					continue
				}
				if data.Model == "" {
					fmt.Println("Error: model not found, don't send update")
					continue
				}
				if roomMap[data.RoomID] == nil {
					roomMap[data.RoomID] = make([]Data, 0)
				}
				roomMap[data.RoomID] = append(roomMap[data.RoomID], data)
			}
			if len(roomMap) > 0 {
				for roomID, roomData := range roomMap {
					// fmt.Println("/room/"+roomID)
					err = io.Of("/room", nil).To(socket.Room(roomID)).Emit("update", roomData)
					if err != nil {
						fmt.Println(err)
						continue
					}
				}
			}
		}
	}()

	exit := make(chan struct{})
	SignalC := make(chan os.Signal)

	signal.Notify(SignalC, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range SignalC {
			switch s {
			case os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close(exit)
				return
			}
		}
	}()

	<-exit
	err = httpServer.Close(nil)
	return
}

func main() {
	//err := ginServer()
	err := socketServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	os.Exit(0)
}
