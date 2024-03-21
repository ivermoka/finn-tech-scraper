package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	fmt.Println("Program started.")
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	s := gocron.NewScheduler(location)

	s.Every(1).Day().At("12:00").Do(func() {
		fmt.Println("Running cron job at ", time.Now().UTC().UnixMilli())
		scraper()
	})

	s.StartBlocking()
}
