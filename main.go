package main

import (
	"log"
	"os"
	"time"

	"Tahlilchi.uz/developer"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")
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

	Router()
}
