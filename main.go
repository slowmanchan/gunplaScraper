package main

import (
	"log"

	"github.com/slowmanchan/gunplaScraper/app"
)

func main() {
	a := app.New()
	if err := a.Scrape(); err != nil {
		log.Println(err)
	}
}
