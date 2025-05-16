package handler

import (
	"net/http"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	logger.Logr(r, "New request", 0)

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if global.SM.CheckSessionExists(cookie.Value) {
		global.SM.DeleteSessionBySessionID(cookie.Value)
		http.SetCookie(w, &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
