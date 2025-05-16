package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

func hashPassword(password string) string {
	if password == "" {
		return ""
	}

	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func GetRegister(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")

	if err != nil {
		logger.Logr(r, err.Error(), 1)

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: "",
		})

		if err := global.Tpl.ExecuteTemplate(w, "register.html", nil); err != nil {
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

	if err := global.Tpl.ExecuteTemplate(w, "register.html", nil); err != nil {
		logger.Logr(r, err.Error(), 1)
		return
	}
}

func PostRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	email := r.FormValue("email")
	uname := r.FormValue("uname")
	pass := hashPassword(r.FormValue("password"))
	lan := r.FormValue("lan")

	if email == "" || pass == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Check if user already exists
	var count int
	row := global.UserDB.QueryRow(
		"SELECT COUNT(*) FROM userdb WHERE email = $1 OR uname = $2",
		email, uname,
	)
	row.Scan(&count)
	if err := row.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)
		return
	}

	if count > 0 {
		w.Write([]byte("Email or Username already registered pick different"))
		return
	}

	// Insert user into database
	_, err := global.UserDB.Exec(
		"INSERT INTO userdb (email, passhash, uname, lang, rdate, ldate) VALUES ($1, $2, $3, $4, $5, $5)",
		email, pass, uname, lan, time.Now(), time.Now(),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
