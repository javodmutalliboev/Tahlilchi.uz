package developer

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func AddAdmin() {
	port, _ := strconv.Atoi(os.Getenv("DBPORT"))
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DBHOST"), port, os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBNAME"))

	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please, enter the name, email, role, password of the new admin:")
	var name, email, role, password string

	name, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	name = strings.TrimSpace(name)

	email, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	email = strings.TrimSpace(email)

	role, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	role = strings.TrimSpace(role)

	password, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	password = strings.TrimSpace(password)

	fmt.Printf("Read line: %s, %s, %s, %s\n", name, email, role, password)

	password, err = HashPassword(password)
	if err != nil {
		panic(err)
	}

	insertStmt := `INSERT INTO public.admins (name, email, role, password, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(insertStmt, name, email, role, password, time.Now())
	if err != nil {
		panic(err)
	}

	fmt.Println("The new admin has been added successfully!")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
