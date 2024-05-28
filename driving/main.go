package main

import (
	"fmt"
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/socket.io/socket"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"src/hall"
	"src/model"
	"src/room"
	"src/utils"
	"syscall"
	"time"
)

func socketServer() (err error) {
	httpServer := types.CreateServer(nil)
	c := socket.DefaultServerOptions()
	c.SetAllowEIO3(true)
	c.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})

	io := socket.NewServer(httpServer, c)

	err = hall.Hall(io)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	err = room.Room(io)
	if err != nil {
		utils.Logger.Error(err.Error())
		return
	}

	httpServer.Listen(":3000", func() {
		fmt.Println("Listening on 3000")
	})

	model.Ticker = time.NewTicker(40 * time.Millisecond)
	defer model.Ticker.Stop()
	go room.Update(io)

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
		utils.Logger.Error(err.Error())
		return
	}
	os.Exit(0)
}
