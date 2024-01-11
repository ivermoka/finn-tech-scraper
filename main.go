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
		// Programming Languages
		"C", "C++", "C#", "Java", "Python", "JavaScript", "TypeScript", "Ruby", "Go", "PHP",
		"Swift", "Kotlin", "Objective-C", "Rust", "Scala", "Perl", "Dart",

		// Web Development
		"HTML", "CSS", "React", "Angular", "Vue.js", "Node.js", "Express.js", "Django", "Flask",
		"Spring Boot", "Laravel",

		// Mobile Development
		"Android", "iOS", "React Native", "Flutter", "SwiftUI",

		// Database Technologies
		"SQL", "MySQL", "PostgreSQL", "MongoDB", "Redis", "SQLite", "Firebase",

		// Data Science and Machine Learning
		"TensorFlow", "PyTorch", "Scikit-learn", "Pandas", "NumPy", "Apache Spark",

		// Cloud Platforms
		"AWS", "Azure", "Google Cloud", "Heroku", "DigitalOcean",

		// Containers and Orchestration
		"Docker", "Kubernetes", "OpenShift",

		// Version Control
		"Git", "GitHub", "GitLab",

		// Continuous Integration/Continuous Deployment
		"Jenkins", "Travis CI", "CircleCI",

		// Web Services and APIs
		"REST", "GraphQL", "SOAP",

		// Microservices and Serverless
		"Microservices", "Serverless", "AWS Lambda", "Azure Functions",

		// DevOps Practices
		"DevOps", "Infrastructure as Code (IaC)", "Configuration Management",

		// Web Frameworks
		"Flask", "Django", "Ruby on Rails", "Express.js", "Spring Boot", "Laravel",

		// Game Development
		"Unity", "Unreal Engine", "Godot",

		// Blockchain
		"Blockchain", "Ethereum", "Hyperledger",

		// Cybersecurity
		"Cybersecurity", "Ethical Hacking", "Penetration Testing",

		// Agile and Scrum
		"Agile", "Scrum", "Kanban",

		// Operating Systems
		"Linux", "Windows", "macOS",

		// Other Technologies
		"RESTful API", "SOAP", "WebSockets",
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