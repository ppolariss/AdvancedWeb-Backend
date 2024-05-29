package room

import (
	"fmt"
	"github.com/zishang520/socket.io/socket"
	. "src/model"
)

func Update(io *socket.Server) {
	var err error
	for range Ticker.C {
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
}
