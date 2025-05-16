package logger

import (
	"net/http"
	"net/url"
	"time"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/sirupsen/logrus"
)

func Log(r *http.Request, errStr string, logType int) error {
	_, err := global.LogDB.Exec(
		`INSERT INTO logs (type, timestamp, ip_address, method, path, agent, message) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`,

		logType,
		time.Now(),
		r.RemoteAddr,
		r.Method,
		r.URL.Path,
		r.UserAgent(),
		errStr,
	)

	if err != nil {
		return err
	}

	return nil
}

func Logr(r *http.Request, errStr string, logType int) {
	if r == nil {
		r = &http.Request{
			RemoteAddr: "localhost",
			Method:     "UNKNOWN",
			URL:        &url.URL{Path: "/"},
			Header:     make(http.Header),
		}
		r.Header.Set("User-Agent", "[SERVER]")
	}

	if logType == 0 {
		logrus.Info("INFO: " + errStr)
	} else {
		logrus.Error("ERROR: " + errStr)
	}
	Log(r, errStr, logType)
}
