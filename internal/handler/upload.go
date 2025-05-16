package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

func generateTid() (int64, error) {
	rand.Seed(time.Now().UnixNano())

	var res int64
	var count int

	for attempts := 0; attempts < 10; attempts++ {
		res = rand.Int63n(900000000000) + 100000000000

		row := global.TorrentDB.QueryRow("SELECT COUNT(*) FROM torrentdb WHERE tid = $1", res)
		if err := row.Scan(&count); err != nil {
			return -1, err
		}

		if count == 0 {
			return res, nil
		}
	}

	return -1, fmt.Errorf("failed to generate a unique tid after %d attempts", 10)
}

func GetUploadTorrent(w http.ResponseWriter, r *http.Request) {
	if err := global.Tpl.ExecuteTemplate(w, "upload.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		logger.Logr(r, err.Error(), 1)

		return
	}
}

func PostUploadTorrent(w http.ResponseWriter, r *http.Request) {
	var uid int
	cookie, err := r.Cookie("session")
	if err != nil {
		uid = -1
	} else {
		if global.SM.CheckSessionExists(cookie.Value) {
			uid = global.SM.GetUIDBySessionID(cookie.Value)
		} else {
			uid = -1
		}
	}

	r.ParseForm()

	tid, err := generateTid()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)

		return
	}

	name := r.FormValue("name")
	cat := r.FormValue("cat")
	link := r.FormValue("link")
	desc := r.FormValue("desc")
	uploadtime := time.Now().Format("2006-01-02 15:04:05")

	if name == "" || cat == "" || link == "" || desc == "" || len(link) > 3000 {
		http.Error(w, "Please fill the fields OR magnet link should be less than 3000 chars", http.StatusInternalServerError)
		return
	}

	catInt, err := strconv.Atoi(cat)
	if err != nil {
		http.Error(w, "Invalid category", http.StatusInternalServerError)
		return
	}

	// Validate magnet link using a regular expression
	if !regexp.MustCompile(`^magnet:\?xt=urn:btih:[0-9a-fA-F]{40,}.*$`).MatchString(link) {
		http.Error(w, "Invalid magnet link", http.StatusInternalServerError)
		return
	}

	// Validate description
	if len(desc) > 30000 {
		http.Error(w, "Description is too long", http.StatusInternalServerError)
		return
	}

	// Register torrent
	_, err = global.TorrentDB.Exec(
		"INSERT INTO torrentdb (tid, uid, name, cat, link, desc, uploadtime) VALUES($1, $2, $3, $4, $5, $6, $7)",
		tid, uid, name, catInt, link, desc, uploadtime,
	)
	if err != nil {
		http.Error(w, "Failed to register torrent", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/success", http.StatusSeeOther)
}
