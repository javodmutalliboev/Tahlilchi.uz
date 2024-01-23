package main

import (
	"fmt"
	"os"
	"time"

	"Tahlilchi.uz/admin"
	"Tahlilchi.uz/developer"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from", r)
		}

		if os.Getenv("ENVIRONMENT") == "development" {
			go func() {
				for {
					exit := developer.Developer()

					if exit {
						break
					}

					time.Sleep(1 * time.Second)
				}
			}()
		}

		// Initialize a new scheduler
		s := gocron.NewScheduler(time.UTC)

		// Schedule the function to run every hour
		s.Every(1).Hour().Do(admin.CheckAndArchiveExpiredBPPosts)

		// Start the scheduler without blocking
		s.StartAsync()

		Router()
	}()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}
