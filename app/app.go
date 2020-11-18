package app

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/jmoiron/sqlx"
	"github.com/slowmanchan/gunplaScraper/models"

	// postgres db driver
	_ "github.com/lib/pq"
)

const (
	indexURL = "https://gunpla.fandom.com"
)

// App holds all major dependancies for the app such
// as http router and the db connection
type App struct {
	db *sqlx.DB
}

func New() *App {
	db, err := sqlx.Open("postgres", "dbname=gunpla sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return &App{
		db: db,
	}
}

func (a *App) Scrape() error {
	sites := []string{
		"/wiki/High_Grade_Universal_Century#1999",
		"/wiki/Real_Grade#2010",
		"/wiki/Real_Grade/Special_Editions",
		"/wiki/Real_Grade/Exclusives",
		"/wiki/Master_Grade",
		"/wiki/Master_Grade/Special_Editions",
		"/wiki/Master_Grade/Exclusives",
		"/wiki/Perfect_Grade#2010",
	}
	for _, v := range sites {
		err := a.scrapeSite(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) scrapeSite(site string) error {
	c := colly.NewCollector(
		colly.CacheDir("./gunpla_cache"),
	)
	dc := c.Clone()

	g := new(models.Gunpla)

	c.OnHTML(".article-table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, e *colly.HTMLElement) {
			if e.ChildAttr("a", "href") == "" {
				return
			}
			dc.Visit(indexURL + e.ChildAttr("a", "href"))
		})
	})

	dc.OnHTML(".mw-parser-output h2:first-of-type + ul", func(e *colly.HTMLElement) {
		g.Includes = e.Text
	})

	dc.OnHTML(".mw-parser-output > h2:nth-of-type(2) + h3 + ul", func(e *colly.HTMLElement) {
		g.Features = e.Text
	})

	dc.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
		g.BoxArtImageLink = e.ChildAttr(".image", "href")
		g.Title = e.ChildText(".pi-title")
		g.Subtitle = e.ChildText("[data-source=subtitle] > div")
		g.Classification = e.ChildText("[data-source=Classification] > div")
		g.LineupNo = e.ChildText("[data-source='Lineup no.'] > div")
		g.Scale = e.ChildText("[data-source=Scale] > div")
		g.Franchise = e.ChildText("[data-source=Franchise] > div")
		g.ReleaseDate = e.ChildText("[data-source='Release Date'] > div")
		g.JanISBN = e.ChildText("[data-source='JAN/ISBN'] > div")
		g.Run = e.ChildText("[data-source=Run] > div")
		g.Price = e.ChildText("[data-source=Price] > div")

		// e.ForEach(".wikia-gallery-item", func(i int, e *colly.HTMLElement) {
		// 	fmt.Println(e.ChildAttr("img", "data-src"))
		// })

	})

	// dc.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	dc.OnScraped(func(r *colly.Response) {
		if g.Title == "" {
			return
		}

		_, err := a.db.Exec(`
			 INSERT INTO gunpla
			 (box_art_image_link, title, subtitle, classification, lineup_no, scale, franchise, release_date, jan_isbn, run, price, features, includes)
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
			`, g.BoxArtImageLink, g.Title, g.Subtitle, g.Classification, g.LineupNo, g.Scale, g.Franchise, g.ReleaseDate, g.JanISBN, g.Run, g.Price, g.Features, g.Includes)
		if err != nil {
			log.Println(err)
		}
	})

	dc.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, "\nError:", err)
	})

	return c.Visit(indexURL + site)
}
