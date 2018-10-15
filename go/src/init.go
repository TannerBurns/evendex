package main

import (
	"log"
	"os"
)

func checkFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	initLogging(os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	db, err := connect()
	checkFatal(err)
	defer db.Close()

	err = db.Ping()
	checkFatal(err)

	Info.Printf("Connected!")

	// create tables
	iniTables(db)

}
