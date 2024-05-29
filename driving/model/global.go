package model

import (
	"sync"
	"time"
)

var GlobalRooms = make(map[string][]string)
var Ticker *time.Ticker
var MutexRooms sync.Mutex
