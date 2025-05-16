package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/Turk1shGuy/torrent/internal/global"
	"github.com/Turk1shGuy/torrent/internal/logger"
)

func Ico(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile(global.Conf.Pages + "/static/logo.png")
	if err != nil {
		w.Write(nil)
		logger.Logr(r, err.Error(), 1)
		return
	}

	w.Write(file)
}

func Sitemap(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadFile(global.Conf.Pages + "/static/sitemap.xml")
	if err != nil {
		w.Write(nil)
		return
	}

	w.Write(file)
}

func Robots(w http.ResponseWriter, r *http.Request) {
	agent := r.UserAgent()

	if !strings.Contains(agent, "Googlebot") || !strings.Contains(agent, "Yandex") ||
		!strings.Contains(agent, "Bingbot") || !strings.Contains(agent, "DuckDuckBot") {

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	file, err := os.ReadFile(global.Conf.Pages + "/static/robots.txt")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Write(file)
}
