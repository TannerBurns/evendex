package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

// Connect to database
func connect() (db *sql.DB, err error) {
	var (
		host     string
		database string
		user     string
		password string
	)
	dat, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		Fatal.Println(err)
	}
	sections := strings.Split(string(dat), "\n\n")
	for i := range sections {
		lines := strings.Split(sections[i], "\n")
		header := lines[0]
		conf := lines[1:]
		if strings.Contains(header, "[postgresql]") {
			for i2 := range conf {
				if strings.Contains(conf[i2], "host") {
					host = strings.Split(conf[i2], "=")[1]
				}
				if strings.Contains(conf[i2], "database") {
					database = strings.Split(conf[i2], "=")[1]
				}
				if strings.Contains(conf[i2], "user") {
					user = strings.Split(conf[i2], "=")[1]
				}
				if strings.Contains(conf[i2], "password") {
					password = strings.Split(conf[i2], "=")[1]
				}
			}
		}
	}

	db, err = sql.Open("postgres", fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, database, user, password))
	return
}

// create tables for evendex database
func iniTables(db *sql.DB) {
	eventTable := `
		CREATE TABLE events (
			event_id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			created VARCHAR(255) NOT NULL
		)
	`

	contentTable := `
		CREATE TABLE content (
			content_id SERIAL PRIMARY KEY,
			event_id INTEGER NOT NULL,
			FOREIGN KEY (event_id) REFERENCES events (event_id) ON DELETE CASCADE,
			created VARCHAR(255),
			status VARCHAR(255),
			version INTEGER NOT NULL,
			title TEXT,
			tag TEXT		
		)
	`

	commentTable := `
		CREATE TABLE comment (
			comment_id SERIAL PRIMARY KEY,
			content_id INTEGER NOT NULL,
			FOREIGN KEY (content_id) REFERENCES content (content_id) ON DELETE CASCADE,
			version INTEGER NOT NULL,
			created VARCHAR(255) NOT NULL,
			modified VARCHAR(255) NOT NULL,
			body TEXT
		)
	`

	labelsTable := `
		CREATE TABLE labels (
			label_id SERIAL PRIMARY KEY,
			content_id INTEGER NOT NULL,
			FOREIGN KEY (content_id) REFERENCES content (content_id) ON DELETE CASCADE,
			comment_id INTEGER NOT NULL,
			FOREIGN KEY (comment_id) REFERENCES comment (comment_id) ON DELETE CASCADE,
			created VARCHAR(255) NOT NULL,
			label TEXT
		)
	`

	tables := []string{eventTable, contentTable, commentTable, labelsTable}

	for _, q := range tables {
		_, err := db.Exec(q)
		if err != nil {
			Fatal.Println(err)
		}
	}
	Info.Println("Done creating tables!")

}
