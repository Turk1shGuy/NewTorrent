package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

func GetLogin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if err != nil {
		logger.Logr(r, err.Error(), 1)

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: "",
		})

		if err := global.Tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
			logger.Logr(r, err.Error(), 1)
			return
		}
		return
	}

	if cookie != nil && cookie.Value != "" {
		if global.SM.CheckSessionExists(cookie.Value) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	if err := global.Tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		logger.Logr(r, err.Error(), 1)
		return
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.FormValue("email")
	pass := hashPassword(r.FormValue("password"))

	if email == "" || pass == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// QUERY SECTION
	var uid int
	{
		row := global.UserDB.QueryRow(
			"SELECT uid FROM userdb WHERE email = $1 AND passhash = $2",
			email, pass,
		)
		err := row.Scan(&uid)
		if err != nil {
			if err == sql.ErrNoRows {
				println(pass)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				logger.Logr(r, err.Error(), 1)
				return
			}
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)
			return
		}
	}

	session_id, err := global.SM.CreateOrUpdateSession(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session",
		Value:   session_id,
		Expires: time.Now().Add(24 * time.Hour),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
