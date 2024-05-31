package model

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
	UserID   int      `json:"user_id"`
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

type Event struct {
	Position Position `json:"position"`
	Rotation Rotation `json:"rotation"`
	SocketID string   `json:"id"`
	RoomID   string   `json:"roomID"`
	Event    string   `json:"event"`
}
