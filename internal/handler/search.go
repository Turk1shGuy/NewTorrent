package handler

import (
	"database/sql"
	"net/http"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

type SearchResults struct {
	Tid        int
	Uid        int
	Uname      string
	Name       string
	Cat        int
	Link       string
	UploadTime string
}

func Search(w http.ResponseWriter, r *http.Request) {
	logger.Logr(r, "New request", 0)

	query := r.URL.Query().Get("q")
	cat := r.URL.Query().Get("cat")

	var rows *sql.Rows
	var err error

	if query != "" && cat != "" {
		rows, err = global.TorrentDB.Query(
			"SELECT tid, uid, name, cat, link, uploadtime FROM torrentdb WHERE name LIKE $1 AND cat = $2",
			"%"+query+"%",
			cat,
		)
	} else if query != "" {
		rows, err = global.TorrentDB.Query(
			"SELECT tid, uid, name, cat, link, uploadtime FROM torrentdb WHERE name LIKE $1",
			"%"+query+"%",
		)
	} else if cat != "" {
		rows, err = global.TorrentDB.Query(
			"SELECT tid, uid, name, cat, link, uploadtime FROM torrentdb WHERE cat = $1",
			cat,
		)
	} else {
		// NOTE HERE: SELECT RANDOM TORRENT FROM DATABASE

		if err := global.Tpl.ExecuteTemplate(w, "search.html", nil); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)
			return
		}

		return
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 0)
		return
	}
	defer rows.Close()

	var res []SearchResults
	for rows.Next() {
		var rez SearchResults

		if err := rows.Scan(&rez.Tid, &rez.Uid, &rez.Name, &rez.Cat, &rez.Link, &rez.UploadTime); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)
			return
		}

		res = append(res, rez)
	}

	// Problem is here - Original re values are not changing
	for i, re := range res {
		if re.Uid == -1 {
			re.Uname = "Anonymous"
			res[i] = re
			continue
		}

		if err = global.UserDB.QueryRow("SELECT uname FROM userdb WHERE uid = $1", re.Uid).Scan(&re.Uname); err != nil {
			if err == sql.ErrNoRows {
				re.Uname = "Unknown"
			} else {

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				logger.Logr(r, err.Error(), 1)
				return
			}
		}

		res[i] = re
	}

	if err := global.Tpl.ExecuteTemplate(w, "search.html", res); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)
		return
	}
}
