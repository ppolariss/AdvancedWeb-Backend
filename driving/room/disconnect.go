package room

import (
	"github.com/zishang520/socket.io/socket"
	. "src/model"
	"src/utils"
)

func Disconnect(io *socket.Server, id socket.SocketId, roomID string) {
	var err error
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

	MutexRooms.Lock()
	delete(GlobalRooms, roomID)
	// removeIdx := -1
	// for idx, thisID := range GlobalRooms[roomID] {
	// 	if thisID == string(id) {
	// 		removeIdx = idx
	// 	}
	// }
	// if removeIdx > -1 {
	// 	GlobalRooms[roomID] = append(GlobalRooms[roomID][:removeIdx], GlobalRooms[roomID][removeIdx+1:]...)
	// }
	MutexRooms.Unlock()
	// io.Sockets().Sockets().Delete(id)
}

// err = io.Of(socket.Room(roomID), nil).Emit("offline", map[string]interface{}{
// 	"id": id,
// })
// io.Sockets().Sockets().Delete(id)
