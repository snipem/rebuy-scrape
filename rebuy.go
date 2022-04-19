package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Offer struct {
	Title    string
	Category string
	Rating   string
	Stock    string
	Price    float64
	URL      string
}

func main() {

	c := colly.NewCollector(colly.CacheDir("rebuy_cache"))

	query := os.Args[1]

	var offers []Offer

	c.OnRequest(func(request *colly.Request) {
		log.Printf("Visiting %s\n", request.URL)
	})

	c.OnHTML("a.page:last-child", func(e *colly.HTMLElement) {
		maxPage, _ := strconv.Atoi(e.Text)
		for i := 2; i <= maxPage; i++ { // Start loop at second page
			c.Visit(getPageURL(query, i))
		}
	})

	c.OnHTML("ry-product", func(e *colly.HTMLElement) {
		offers = append(offers, Offer{
			Title:    strings.TrimSpace(e.ChildText("div.title")),
			Category: e.ChildText("div.category"),
			Rating:   e.ChildAttr("div.rating", "title"),
			Stock:    e.ChildText("div.stock-count"),
			Price:    parsePrice(e.ChildText("span.price-font-size")),
			URL:      e.Request.AbsoluteURL(e.ChildAttr("a", "href")),
		})
	})

	// Start at first page
	c.Visit(getPageURL(query, 1))
	c.Wait()

	fmt.Println("Titel;Kategorie;Preis;URL")
	for i, _ := range offers {
		fmt.Printf("%s;%s;%.2f;\"%s\"\n", strings.ReplaceAll(offers[i].Title, ";", ""), offers[i].Category, offers[i].Price, offers[i].URL)
	}
}

func getPageURL(query string, i int) string {
	return query + fmt.Sprintf("&page=%d", i)
}

func parsePrice(text string) float64 {
	r, err := regexp.Compile("[0-9]*,[0-9]{2}")
	if err != nil {
		log.Fatal(err)
	}
	s := r.FindString(strings.TrimSpace(text))
	formattedS := strings.ReplaceAll(s, ",", ".")
	f, err := strconv.ParseFloat(formattedS, 64)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
