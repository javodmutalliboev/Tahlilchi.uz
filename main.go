package main

import (
	"fmt"
	"os"
	"time"

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

		s := gocron.NewScheduler(time.UTC)
		// Schedule a task to run every 5 seconds
		_, err := s.Every(5).Seconds().Do(task)
		if err != nil {
			fmt.Println("Error scheduling task:", err)
			return
		}
		// Start the scheduler
		s.StartAsync()

		Router()
	}()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func task() {
	fmt.Println("Task is being performed.")
}
