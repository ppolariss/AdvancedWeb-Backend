package model

import "time"

var GlobalRooms = make(map[string][]string)
var Ticker *time.Ticker
