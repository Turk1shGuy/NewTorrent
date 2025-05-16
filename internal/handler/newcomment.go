package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

func newCid() (string, error) {
	for {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
		var builder strings.Builder

		for i := 0; i < 11; i++ {
			digit := rnd.Intn(10)
			builder.WriteByte('0' + byte(digit))
		}

		cid := builder.String()

		var count int
		err := global.CommentDB.QueryRow("SELECT COUNT(*) FROM commentdb WHERE cid = ?", cid).Scan(&count)
		if err != nil {
			return "", err
		}

		if count == 0 {
			return cid, nil
		}
	}
}

func checkTid(tid string) (bool, error) {
	var num int
	err := global.TorrentDB.QueryRow("SELECT COUNT(*) FROM torrentdb WHERE tid = ?", tid).Scan(&num)
	if err != nil {
		return false, err
	}

	return num > 0, nil
}

func NewComment(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: "",
		})
		return
	}

	r.ParseForm()
	uid := global.SM.GetUIDBySessionID(cookie.Value)
	tid := r.URL.Query().Get("tid")
	comment := strings.TrimSpace(r.FormValue("comment"))

	if tid == "" || comment == "" {
		http.Error(w, "Fill the values", http.StatusBadRequest)
		logger.Logr(r, "Fill the values", 1)
		return
	}

	if len(comment) > 1000 {
		http.Error(w, "Comment should be less than 1000 characters", http.StatusBadRequest)
		return
	}

	cid, err := newCid()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)
		return
	}

	exists, err := checkTid(tid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)
		return
	}

	if !exists {
		http.Error(w, "Invalid tid", http.StatusBadRequest)
		return
	}

	_, err = global.CommentDB.Exec(
		"INSERT INTO commentdb (cid, tid, uid, comment, uploadtime) VALUES (?, ?, ?, ?, ?)",
		cid, tid, uid, comment, time.Now().Format(time.RFC3339),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Logr(r, err.Error(), 1)
		return
	}

	logger.Logr(r, "Comment created successfully", 0)

	http.Redirect(w, r, fmt.Sprintf("/detail?tid=%s", tid), http.StatusSeeOther)

}
