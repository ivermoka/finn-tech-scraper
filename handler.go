package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	s := gocron.NewScheduler(time.Now().Location())

	s.Every(1).Day().At("12:00").Do(func() {
		fmt.Println("Running cron job at 12:00.")
		scraper()
	})

	s.StartBlocking()
}
