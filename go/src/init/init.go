package main

import (
	"io/ioutil"
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

	Info.Printf("DB - Connected!")

	// create tables
	iniTables(db)

	Info.Printf("Creating ES Index")
	dat, err := ioutil.ReadFile("src/init/mapping/evendex_map.json")
	if err != nil {
		log.Fatal(err)
	}
	resp := createIndex("evendex", string(dat))
	Info.Printf(resp)

}
