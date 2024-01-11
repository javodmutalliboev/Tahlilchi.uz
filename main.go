package main

import (
	"fmt"
	"os"
	"time"

	"Tahlilchi.uz/developer"
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

		Router()
	}()

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}
