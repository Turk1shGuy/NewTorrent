package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	_ "github.com/mattn/go-sqlite3"
)

const tpb = "https://tpb.party"
const dbPath = "./dbs/torrent.db"

var top11 []string = []string{
	"The Shawshank Redemption",
	"The Godfather",
	"The Dark Knight",
	"The Godfather Part II",
	"12 Angry Men",
	"The Lord of the Rings: The Return of the King",
	"Schindler's List",
	"Pulp Fiction",
	"The Lord of the Rings: The Fellowship of the Ring",
	"Fight Club",
	"Shin Godzilla",
}

var db *sql.DB

func main() {
	var err error
	if db, err = sql.Open("sqlite3", dbPath); err != nil {
		log.Println(err.Error())
		return
	}
	defer db.Close()

	c := colly.NewCollector(
		colly.AllowURLRevisit(),
	)

	c.OnHTML("table#searchResult tr", func(e *colly.HTMLElement) {
		magnetPageLink := e.ChildAttr("td:nth-child(2) a", "href")
		if magnetPageLink != "" {
			if err := c.Visit(magnetPageLink); err != nil {
				log.Println(err.Error())
				return
			}
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if href := e.Attr("href"); len(href) > 0 && strings.HasPrefix(href, "magnet:") {
			tid, err := generateTid()
			if err != nil {
				log.Println(err.Error())
				return
			}

			parts := strings.Split(e.Request.URL.Path, "/")
			name := parts[len(parts)-1]

			// Get the name from the #title div
			nameSelector := e.DOM.Find("div#title")
			if nameSelector.Length() > 0 {
				name = nameSelector.Text()
			}

			desc := ""
			// Get the description from the .nfo pre
			descSelector := e.DOM.Find("div.nfo pre")
			if descSelector.Length() > 0 {
				desc = descSelector.Text()
			} else {
				// If description is not found in the current element, try to find it in the parent element
				parent := e.DOM.Parents().Find("div.nfo pre")
				if parent.Length() > 0 {
					desc = parent.Text()
				}
			}

			// Insert data into database
			_, err = db.Exec("INSERT INTO torrentdb (tid, uid, name, cat, link, desc, uploadtime) VALUES ($1, $2, $3, $4, $5, $6, $7)",
				tid, -1, name, 201, href, desc, time.Now())
			if err != nil {
				log.Println(err.Error())
				return
			}

			fmt.Printf("\nFound magnet link:%s\n", href)
		}
	})

	var wg sync.WaitGroup
	for _, movie := range top11 {
		wg.Add(1)
		go func(movie string) {
			defer wg.Done()
			movieQuery := strings.ReplaceAll(movie, " ", "+")
			err := c.Visit(tpb + "/search/" + movieQuery + "/0/99/0")
			if err != nil {
				log.Fatal(err)
			}
			// Add a delay to avoid overwhelming the server
			time.Sleep(1 * time.Second)
		}(movie)
	}

	wg.Wait()
}

func generateTid() (int64, error) {
	rand.Seed(time.Now().UnixNano())

	var res int64
	var count int

	for attempts := 0; attempts < 10; attempts++ {
		res = rand.Int63n(900000000000) + 100000000000

		row := db.QueryRow("SELECT COUNT(*) FROM torrentdb WHERE tid = $1", res)
		if err := row.Scan(&count); err != nil {
			return -1, err
		}

		if count == 0 {
			return res, nil
		}
	}

	return -1, fmt.Errorf("failed to generate a unique tid after %d attempts", 10)
}
