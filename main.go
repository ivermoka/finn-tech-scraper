package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Ad struct {
	url, title, company, location, tech, content string
}

func main() {
	var ads []Ad
	c := colly.NewCollector()
	
	c.OnError(func(r *colly.Response, err error) {
        log.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError: ", err)
    })

	// c.OnRequest(func (r *colly.Request)  {
	// 	fmt.Println("visiting: ", r.URL)
	// })

	c.OnHTML("article.sf-search-ad-legendary", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		if link != "" {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnHTML("div.u-word-break", func(e *colly.HTMLElement) {
		ad := Ad{}

		ad.url = e.Request.URL.String()
		ad.title = e.ChildText(".u-t2")
		ad.company = e.DOM.Find("dl.definition-list>:nth-child(2)").First().Text()
		ad.location = e.ChildText("dl.definition-list--inline>dd:nth-child(4)") // FUNKER IKKE
		// ad.content = e.ChildText("div.import-decoration")
		fmt.Println(ad)

		ads = append(ads, ad)
	})

	c.Visit("https://www.finn.no/job/fulltime/search.html?occupation=0.23")

	file, err := os.Create("products.csv")
	if err != nil {
		log.Fatalln("Error with creation of CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"url",
		"title",
		"company",
		"location",
	}

	err = writer.Write(headers)
	if err != nil {
		log.Fatalln("Error writing CSV headers", err)
	}

	for _, ad := range ads {
		record := []string {
			ad.url,
			ad.title,
			ad.company,
			ad.location,
			// ad.content,
		}
		err = writer.Write(record)
		if err != nil {
			log.Fatalln("Error writing csv file.", err)
		}
	}

	defer writer.Flush()
}