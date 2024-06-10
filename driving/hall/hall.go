package hall

import (
	// "encoding/json"
	"fmt"
	"github.com/zishang520/socket.io/socket"
	. "src/model"
	"src/utils"
	"strconv"
)

func Hall(io *socket.Server) (err error) {
	err = io.Of("/hall", nil).On("connection", func(clients ...any) {
		client := clients[0].(*socket.Socket)
		err = client.Emit("sendRooms", map[string]interface{}{
			"rooms": GlobalRooms,
		})
		if err != nil {
			utils.Logger.Error("Error in sendRooms " + err.Error())
			// fmt.Println("Error in sendRooms", err)
		}
		// data, _ := json.Marshal(GlobalRooms)
		// utils.Logger.Info("sendRooms " + string(data))

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
			MutexRooms.Lock()
			GlobalRooms["room"+strconv.Itoa(len(GlobalRooms)+1)] = make([]string, 0)
			MutexRooms.Unlock()
			err = client.Emit("sendRooms", map[string]interface{}{
				"rooms": GlobalRooms,
			})
		})
		if err != nil {
			fmt.Println("Error in createRooms", err)
		}
	})
	return
}
