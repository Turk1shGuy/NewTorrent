package handler

import (
	"net/http"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

func Catagories(w http.ResponseWriter, r *http.Request) {
	if err := global.Tpl.ExecuteTemplate(w, "catagories.html", nil); err != nil {
		logger.Logr(r, err.Error(), 1)
		return
	}
}
