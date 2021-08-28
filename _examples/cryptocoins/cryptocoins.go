package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	fName := "cryptocoinmarketcap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Name", "Symbol", "Price (USD)", "Volume (USD)", "Market capacity (USD)", "Change (1h)", "Change (24h)", "Change (7d)"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML(".cmc-table__table-wrapper-outer tbody tr", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.ChildText("a.cmc-table__column-name--name"),
			e.ChildText("td.cmc-table__cell--sort-by__symbol div"),
			e.ChildText("td.cmc-table__cell--sort-by__price a"),
			e.ChildText("td.cmc-table__cell--sort-by__volume-24-h a"),
			e.ChildText("td.cmc-table__cell--sort-by__market-cap span.sc-1ow4cwt-1"),
			e.ChildText("td.cmc-table__cell--sort-by__percent-change-1-h div"),
			e.ChildText("td.cmc-table__cell--sort-by__percent-change-24-h div"),
			e.ChildText("td.cmc-table__cell--right cmc-table__cell--sort-by__percent-change-7-d div"),
		})
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
