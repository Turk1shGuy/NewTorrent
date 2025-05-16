package handler

import (
	"net/http"

	"github.com/Turk1shGuy/torrent/internal/global"
)

func Success(w http.ResponseWriter, r *http.Request) {

	if err := global.Tpl.ExecuteTemplate(w, "success.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
