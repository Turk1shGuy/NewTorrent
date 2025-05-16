package handler

import (
	"net/http"
	"sort"
	"time"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

type announcement struct {
	DateTime string
	Text     string
}

func MakeAnnouncement(w http.ResponseWriter, r *http.Request) {
	logger.Logr(r, "New announcement", 0)

	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, "Session cookie not found", 1)
		return
	}

	if !global.SM.CheckSessionExists(cookie.Value) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		logger.Logr(r, "Invalid session", 1)
		return
	}

	if global.SM.GetUIDBySessionID(cookie.Value) != 0 {
		http.Error(w, "You are not authorized to make announcements", http.StatusForbidden)
		logger.Logr(r, "User is not an admin", 1)
		return
	}

	r.ParseForm()
	text := r.FormValue("text")
	if text == "" || len(text) < 10 {
		http.Error(w, "Announcement text must be at least 10 characters long", http.StatusBadRequest)
		return
	}

	_, err = global.AnnouncementsDB.Exec(
		"INSERT INTO announcements (_time, _text) VALUES(?, ?);",
		time.Now().Format("2006-01-02 15:04:05.000"),
		text,
	)

	if err != nil {
		http.Error(w, "Failed to make announcement", http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)
		return
	}

	http.Redirect(w, r, "/announcements", http.StatusSeeOther)
}

func Announcements(w http.ResponseWriter, r *http.Request) {
	logger.Logr(r, "New request", 0)

	rows, err := global.AnnouncementsDB.Query("SELECT * FROM announcements")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)

		return
	}

	var data []announcement
	for rows.Next() {
		var d announcement
		if err = rows.Scan(&d.DateTime, &d.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)

			return
		}

		data = append(data, d)
	}

	// Reverse it
	sort.Slice(data, func(i, j int) bool {
		return data[i].DateTime > data[j].DateTime
	})

	if err := global.Tpl.ExecuteTemplate(w, "announcements.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)

		return
	}
}
