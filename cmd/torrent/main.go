package main

import (
	"fmt"
	"net/http"

	"github.com/Turk1shGuy/torrent/internal/global"
	INIT "github.com/Turk1shGuy/torrent/internal/init"
	"github.com/Turk1shGuy/torrent/internal/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	err := INIT.INIT("./conf.json")
	if err != nil {
		logrus.Error(err.Error())
		select {}
	}

	logger.Logr(nil, fmt.Sprintf("SERVER STARTED ON PORT: %v", global.Conf.Port), 0)
	if err := http.ListenAndServe(":"+global.Conf.Port, nil); err != nil {
		logrus.Error(err.Error())
		select {}
	}
	select {}
}
