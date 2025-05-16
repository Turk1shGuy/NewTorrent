package handler

import (
	"fmt"
	"net/http"
	"net/smtp"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

type user struct {
	Uid   string
	Uname string
	Email string
}

func Admin(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, "Some one want to access admin panel", 1)

		return
	}
	if !global.SM.CheckSessionExists(cookie.Value) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, "Some one want to access admin panel", 1)

		return
	}
	if global.SM.GetUIDBySessionID(cookie.Value) != 0 { // Admin
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, http.StatusText(http.StatusNotFound), 1)
		return
	}

	// Authorized now
	var Users []user
	rows, err := global.UserDB.Query("SELECT uid, uname, email FROM userdb ORDER BY RANDOM() LIMIT 30")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		logger.Logr(r, err.Error(), 1)

		return
	}

	for rows.Next() {
		var u user
		if err = rows.Scan(&u.Uid, &u.Uname, &u.Email); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)

			return
		}

		Users = append(Users, u)
	}

	if err := global.Tpl.ExecuteTemplate(w, "admin.html", Users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)

		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, err.Error(), 1)
		return
	}

	if !global.SM.CheckSessionExists(cookie.Value) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, http.StatusText(http.StatusNotFound), 1)
		return
	}

	if global.SM.GetUIDBySessionID(cookie.Value) != 0 { // Admin
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, http.StatusText(http.StatusNotFound), 1)
		return
	}

	r.ParseForm()
	uid := r.FormValue("uid")

	_, err = global.UserDB.Exec("DELETE FROM userdb WHERE uid = ?", uid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, err.Error(), 1)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func SendMail(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, err.Error(), 1)
		return
	}

	if !global.SM.CheckSessionExists(cookie.Value) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, http.StatusText(http.StatusNotFound), 1)
		return
	}

	if global.SM.GetUIDBySessionID(cookie.Value) != 0 { // Admin
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, http.StatusText(http.StatusNotFound), 1)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	subject := r.FormValue("subject")
	body := r.FormValue("body")

	if subject == "" || body == "" {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	rows, err := global.UserDB.Query("SELECT email FROM userdb")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)
		return
	}
	defer rows.Close()

	var mail string
	for rows.Next() {
		if err = rows.Scan(&mail); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)
			return
		}

		// Set up the email headers
		headers := make(map[string]string)
		headers["From"] = "test@gmail.com"
		headers["To"] = mail
		headers["Subject"] = subject
		headers["Content-Type"] = "text/plain; charset=UTF-8"

		// Set up the email message
		msg := ""
		for k, v := range headers {
			msg += fmt.Sprintf("%s: %s\r\n", k, v)
		}
		msg += "\r\n" + body

		// Send the email
		err = smtp.SendMail(fmt.Sprintf("%s:%d", "smtp.google.com", 587), nil, "test@gmail.com", []string{mail}, []byte(msg))
		if err != nil {
			http.Error(w, "Failed to send email", http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)
			return
		}
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)
		return
	}

	http.Redirect(w, r, "/success", http.StatusSeeOther)
}
