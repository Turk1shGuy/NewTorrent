package handler

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

type Comment struct {
	Cid        int
	Uid        int
	Uname      string
	Comment    string
	UploadTime string
}

type detail_user_torrent struct {
	Uid        int
	Name       string
	Cat        int
	Link       string
	Desc       string
	UploadTime string
	Tid        string

	Comments []Comment
}

type detail_torrent struct {
	Tid        int
	Cat        int
	Link       string
	UploadTime string
	Name       string
}

func Detail(w http.ResponseWriter, r *http.Request) {
	tid := r.URL.Query().Get("tid")
	uid := r.URL.Query().Get("uid")

	if tid != "" && uid != "" {

		http.Redirect(w, r, "/", http.StatusFound)
		return

	} else if tid != "" {
		row := global.TorrentDB.QueryRow(
			"SELECT uid, name, cat, link, desc, uploadtime FROM torrentdb WHERE tid = $1",
			tid,
		)

		var Data detail_user_torrent
		err := row.Scan(&Data.Uid, &Data.Name, &Data.Cat, &Data.Link, &Data.Desc, &Data.UploadTime)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Torrent not found", http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)
			return
		}

		Data.Desc = strings.Replace(Data.Desc, "\n", "<br>", -1)
		Data.Tid = tid

		//////////////////////////////////////////////////////////
		rows, err := global.CommentDB.Query("SELECT cid, uid, comment, uploadtime FROM commentdb WHERE tid = ?", tid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var comment Comment
			err = rows.Scan(
				&comment.Cid,
				&comment.Uid,
				&comment.Comment,
				&comment.UploadTime,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				logger.Logr(r, err.Error(), 1)
				return
			}
			Data.Comments = append(Data.Comments, comment)
		}

		for i, e := range Data.Comments {
			row := global.UserDB.QueryRow("SELECT uname FROM userdb WHERE uid = ?", e.Uid)
			var uname string
			err := row.Scan(&uname)
			if err != nil {
				if err == sql.ErrNoRows {
					// Handle the case where the user is not found
					// For example, you could set the username to "Unknown"
					Data.Comments[i].Uname = "Unknown"
				} else {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					logger.Logr(r, err.Error(), 1)
					return
				}
			} else {
				Data.Comments[i].Uname = uname
			}
		}

		//////////////////////////////////////////////////////////

		if err := global.Tpl.ExecuteTemplate(w, "detail_torrent.html", Data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)
			return
		}
	} else if uid != "" {

		if uid == "-1" {
			http.Error(w, "You can't see Anonymous torrents", http.StatusNotFound)
			return
		}

		var uname string
		if err := global.UserDB.QueryRow("SELECT uname FROM userdb WHERE uid = $1", uid).Scan(&uname); err != nil {
			http.Error(w, "User not found", http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)

			return
		}

		// if exists get torrents
		rows, err := global.TorrentDB.Query("SELECT tid, cat, link, uploadtime, name FROM torrentdb WHERE uid = $1", uid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)

			return
		}

		data := struct {
			Uname    string
			Torrents []detail_torrent
		}{
			Uname: uname,
		}

		for rows.Next() {
			var torrent detail_torrent
			err = rows.Scan(&torrent.Tid, &torrent.Cat, &torrent.Link, &torrent.UploadTime, &torrent.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				logger.Logr(r, err.Error(), 1)

				return
			}
			data.Torrents = append(data.Torrents, torrent)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)

			return
		}

		if err := global.Tpl.ExecuteTemplate(w, "detail_user.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Logr(r, err.Error(), 1)

			return
		}

	} else {

		http.Redirect(w, r, "/", http.StatusFound)
		return

	}
}
