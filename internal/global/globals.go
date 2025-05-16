package global

import (
	"database/sql"
	"net/smtp"
	"text/template"
	"time"

	"github.com/Turk1shGuy/torrent/internal/session"
)

type C0nf struct {
	Port           string `json:"port"`
	Pages          string `json:"pages"`
	UserDB         string `json:"userdb"`
	TorrentDB      string `json:"torrentdb"`
	CommentDB      string `json:"commentdb"`
	LogDB          string `json:"logdb"`
	AnnouncementDB string `json:"announcementsdb"`
}

var (
	Tpl *template.Template
	SM  = session.NewSessionManager(24*time.Hour, 20*time.Minute)

	Conf C0nf

	UserDB          *sql.DB
	TorrentDB       *sql.DB
	CommentDB       *sql.DB
	LogDB           *sql.DB
	AnnouncementsDB *sql.DB

	Auth smtp.Auth
)
