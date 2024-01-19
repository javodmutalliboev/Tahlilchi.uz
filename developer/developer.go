package developer

import (
	"fmt"
	"log"

	"Tahlilchi.uz/telegramBot"
)

func Developer() bool {
	var exit bool

	log.Println("Hello, developer.")
	fmt.Println("1. Exit")
	fmt.Println("2. Add admin")
	fmt.Println("3. Get telegram bot chat id")
	fmt.Println("What do you want to do? Please enter 1 or 2 or 3:")
	var decision int
	_, err := fmt.Scan(&decision)

	if err != nil {
		log.Fatal(err)
	}

	if decision == 1 {
		exit = true
		return exit
	} else if decision == 2 {
		// AddAdmin
		AddAdmin()
	} else if decision == 3 {
		telegramBot.ChatID()
	}

	exit = true
	return exit
}
