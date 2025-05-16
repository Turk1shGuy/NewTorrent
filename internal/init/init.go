package INIT

import (
	"database/sql"
	"net/http"
	"net/smtp"
	"text/template"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/handler"

	_ "github.com/mattn/go-sqlite3"
)

func INIT(conf_str string) error {
	global.Conf = readConf(conf_str)

	var err error
	global.Tpl, err = template.ParseGlob(global.Conf.Pages + "/*.html")
	if err != nil {
		return err
	}

	// Open databases
	if global.UserDB, err = sql.Open("sqlite3", global.Conf.UserDB); err != nil {
		return err
	}
	if global.TorrentDB, err = sql.Open("sqlite3", global.Conf.TorrentDB); err != nil {
		return err
	}
	if global.CommentDB, err = sql.Open("sqlite3", global.Conf.CommentDB); err != nil {
		return err
	}
	if global.LogDB, err = sql.Open("sqlite3", global.Conf.LogDB); err != nil {
		return err
	}
	if global.AnnouncementsDB, err = sql.Open("sqlite3", global.Conf.AnnouncementDB); err != nil {
		return err
	}

	// EMail server
	global.Auth = smtp.PlainAuth("", "test@gmail.com", "passwd1234", "smtp.google.com")

	// Handlers
	http.HandleFunc("GET /search", handler.Search)
	http.HandleFunc("GET /catagories", handler.Catagories)

	http.HandleFunc("GET /login", handler.GetLogin)
	http.HandleFunc("POST /login", handler.PostLogin)

	http.HandleFunc("POST /logout", handler.Logout)

	http.HandleFunc("GET /register", handler.GetRegister)
	http.HandleFunc("POST /register", handler.PostRegister)

	http.HandleFunc("POST /makecomment", handler.NewComment)

	http.HandleFunc("GET /announcements", handler.Announcements)
	http.HandleFunc("POST /makeannouncements", handler.MakeAnnouncement)

	http.HandleFunc("GET /admin", handler.Admin)
	http.HandleFunc("POST /deleteuser", handler.DeleteUser)
	http.HandleFunc("POST /sendmail", handler.SendMail)

	http.HandleFunc("GET /about", handler.About)
	http.HandleFunc("GET /detail", handler.Detail)

	http.HandleFunc("GET /favicon.png", handler.Ico)
	http.HandleFunc("GET /favicon.ico", handler.Ico)
	http.HandleFunc("GET /sitemap.xml", handler.Sitemap)
	http.HandleFunc("GET /robots.txt", handler.Robots)

	http.HandleFunc("GET /404", handler.P404)

	http.HandleFunc("GET /upload", handler.GetUploadTorrent)
	http.HandleFunc("POST /upload", handler.PostUploadTorrent)

	http.HandleFunc("GET /success", handler.Success)

	http.HandleFunc("GET /", handler.Index)

	return nil
}
