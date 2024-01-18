package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Ad struct {
	tech                        []string
}

func scraper() {
	Connect()
	defer db.Close()
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

	c.OnHTML("article.sf-search-ad-legendary", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		if link != "" {
			c.Visit(e.Request.AbsoluteURL(link))
		}
	})

	c.OnHTML("div.u-word-break", func(e *colly.HTMLElement) {
		ad := Ad{}

		content := e.ChildText("div.import-decoration")
		ad.tech = grabTech(content)
		ads = append(ads, ad)
	})
	c.OnHTML("div[data-testid='aggregated-ad-object']", func(e *colly.HTMLElement) {
		ad := Ad{}

		content := e.ChildText("section.mt-28")
		ad.tech = grabTech(content)
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

	for _, ad := range ads {

		for _, tech := range ad.tech {
			techCount[tech]++
		}
	}
	cleanDB()
	for tech, count := range techCount {
		uploadToDB(tech, count)
	}
}

func grabTech(content string) []string {
	techPatterns := []string{
		// Programming Languages
		"C", "C++", "C#", "Java", "Python", "JavaScript", "TypeScript", "Ruby", "Go", "PHP",
		"Swift", "Kotlin", "Rust", "Scala", "Perl", "Dart", "OCaml", "Zig",

		// Web Development
		"HTML", "CSS", "React", "Angular", "Vue.js", "Node.js", "Express.js", "Django", "Flask",
		"Spring Boot", "Laravel", ".NET",

		// Mobile Development
		"React Native", "Flutter", "SwiftUI",

		// Database Technologies
		"SQL", "MySQL", "PostgreSQL", "MongoDB", "Redis", "SQLite", "Firebase",

		// Data Science and Machine Learning
		"TensorFlow", "PyTorch", "Scikit-learn", "Pandas", "NumPy", "Apache Spark",

		// Cloud Platforms
		"AWS", "Azure", "Google Cloud", "Heroku", "DigitalOcean", "Linode",

		// Containers and Orchestration
		"Docker", "Kubernetes", "OpenShift", "k8s", 

		// Version Control
		"Git", "GitHub", "GitLab",

		// Continuous Integration/Continuous Deployment
		"Jenkins", "Travis CI", "CircleCI",

		// Web Services and APIs
		"REST", "GraphQL", "SOAP",

		// Microservices and Serverless
		"AWS Lambda", "Azure Functions",

		// DevOps Practices
		"DevOps",

		// Web Frameworks
		"Flask", "Django", "Ruby on Rails", "Express.js", "Spring Boot", "Laravel",

		// Game Development
		"Unity", "Unreal Engine", "Godot",

		// Cybersecurity
		"Cybersecurity", "Ethical Hacking", "Penetration Testing",

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
