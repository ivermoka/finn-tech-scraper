package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)
type Ad struct {
	url, title, company, origin string
	tech []string
}

func main() {
	var ads []Ad
	techCount := make(map[string]int)
	c := colly.NewCollector()
	
	
	c.OnError(func(r *colly.Response, err error) {
        log.Println("Request URL: ", r.Request.URL, " failed with response: ", r, "\nError: ", err)
    })
	fmt.Println("No errors. Starting scrape. ")
	
	go func() {
		startTime := time.Now()
		for {
			elapsed := time.Since(startTime)
			fmt.Printf("\rRunning scrape. Elapsed time: %v", elapsed.Round(time.Second)) 
			time.Sleep(time.Second)                                      
		}
	}()

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
		content := e.ChildText("div.import-decoration")
		ad.tech = grabTech(content)
		ad.origin = "finn.no"
		// fmt.Println(ad)
		ads = append(ads, ad)
	})
	c.OnHTML("div[data-testid='aggregated-ad-object']", func(e *colly.HTMLElement) {
		ad := Ad{}

		ad.url = e.Request.URL.String()
		ad.title = e.ChildText("h1.mb-32")
		ad.company = e.DOM.Find("dl.space-y-8>:nth-child(2)").First().Text() 
		content := e.ChildText("section.mt-28")
		ad.tech = grabTech(content)
		ad.origin = "NAV"
		// fmt.Println("FROM NAV: ", ad)
		ads = append(ads, ad)
	})
	c.OnHTML("a.button--icon-right", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
        if nextPage != "" {
            e.Request.Visit(nextPage)
        }
	})

	if err := c.Visit("https://www.finn.no/job/fulltime/search.html?occupation=0.23"); err != nil {
        log.Fatal(err)
    }

	fmt.Println("\nScraping finished. Creating csv file.")

	file, err := os.Create("ads.csv")
	if err != nil {
		log.Fatalln("Error with creation of CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"url",
		"title",
		"company",
		"origin",
		"tech",
	}

	err = writer.Write(headers)
	if err != nil {
		log.Fatalln("Error writing CSV headers", err)
	}


	for _, ad := range ads {
		for _, tech := range ad.tech {
			techCount[tech]++
		}
		record := []string {
			ad.url,
			ad.title,
			ad.company,
			ad.origin,
			strings.Join(ad.tech, ", "),
			// ad.location,
			// ad.content,
		}
		err = writer.Write(record)
		if err != nil {
			log.Fatalln("Error writing csv file.", err)
		}
	}

	defer writer.Flush()
	// allAds, err := json.Marshal(ads)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// err = os.WriteFile("ads.json", allAds, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("Creating counts.json file.")
	
	counts, err := json.Marshal(techCount)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("counts.json", counts, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println("Technology Counts:")
	// for tech, count := range techCount {
	// 	fmt.Printf("%s: %d\n", tech, count)
	// }
}

func grabTech(content string) []string {
	techPatterns := []string{
		`SQL`, `Looker`, `PowerBI`, `Tableau`, `BigQuery`, `dbt`, `Python`, `Java`, `JavaScript`,
		`C\+\+`, `C#`, `HTML`, `CSS`, `React`, `Angular`, `Node\.js`, `Vue\.js`, `Swift`, `Kotlin`,
		`Ruby`, `Go(lang)?`, `PHP`, `Scala`, `Perl`, `Rust`, `TypeScript`, `Dart`, `Objective-C`,
		`TensorFlow`, `PyTorch`, `Keras`, `Scikit-learn`, `Pandas`, `NumPy`, `Spark`, `Hadoop`,
		`AWS`, `Azure`, `Google Cloud`, `Docker`, `Kubernetes`, `Git`, `Jenkins`, `CI/CD`,
		`Machine Learning`, `Deep Learning`, `Artificial Intelligence`, `Data Science`,
		`Blockchain`, `Cybersecurity`, `DevOps`, `Agile`, `Scrum`, `REST`, `GraphQL`, `API`,
		`Microservices`, `Serverless`, `NoSQL`, `MongoDB`, `Redis`, `PostgreSQL`, `MySQL`, `SQLite`,
		`Linux`, `Windows`, `macOS`, `iOS`, `Android`, `React Native`, `Flutter`, `Unity`,
	}

	techPattern := fmt.Sprintf(`\b(%s)\b`, strings.Join(techPatterns, `|`))
	r := regexp.MustCompile(techPattern)
	matches := r.FindAllString(content, -1)

	uniqueMatches := make(map[string]bool)
	for _, match := range matches {
		uniqueMatches[match] = true
	}
	var techList []string
	for tech := range uniqueMatches {
		techList = append(techList, tech)
	}
	return techList
}