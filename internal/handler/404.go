package handler

import (
	"net/http"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

func P404(w http.ResponseWriter, r *http.Request) {
	if err := global.Tpl.ExecuteTemplate(w, "404.html", nil); err != nil {
		http.Error(w, "Some thing went wrong", http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)

		return
	}
}
