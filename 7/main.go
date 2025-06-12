package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

func main() {
	password := flag.String("p", "", "password")
	flag.Parse()
	connStr := "user=postgres password=" + *password + " dbname=phone_normalizer sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL PRIMARY KEY,
		number VARCHAR(20) NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
	phoneNumbers := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, phoneNumber := range phoneNumbers {
		db.Exec(`INSERT INTO phone_numbers(number) VALUES ($1)`, phoneNumber)
	}
	rows, err := db.Query("SELECT id, number FROM phone_numbers")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int
		var number string
		if err := rows.Scan(&id, &number); err != nil {
			log.Fatal(err)
		}
		normalizedNumber := normalizePhoneNumber(number)
		db.Exec(`UPDATE phone_numbers SET number = $1 WHERE id = $2`, normalizedNumber, id)

	}
	_, err = db.Exec(`
		DELETE FROM phone_numbers 
		WHERE id NOT IN (
			SELECT min(id) 
			FROM phone_numbers 
			GROUP BY number
		)`)
	if err != nil {
		log.Fatal(err)
	}
	rows, err = db.Query("SELECT number FROM phone_numbers ORDER BY number")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Unique phone numbers:")
	for rows.Next() {
		var number string
		if err := rows.Scan(&number); err != nil {
			log.Fatal(err)
		}
		fmt.Println(number)
	}
}

func normalizePhoneNumber(number string) string {
	var sb strings.Builder
	for _, val := range number {
		if val >= '0' && val <= '9' {
			sb.WriteString(string(val))
		}
	}
	return sb.String()
}
