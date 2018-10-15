package main

import (
	"encoding/json"
	"os"
)

func main() {
	initLogging(os.Stdout, os.Stdout, os.Stdout, os.Stderr, os.Stderr)
	db, err := connect()
	if err != nil {
		Fatal.Println(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		Fatal.Println(err)
	}

	Info.Printf("Connected!\n")

	comments, err := bulkComment(db)
	if err != nil {
		Fatal.Println(err)
	}

	Info.Printf("Checking ES status\n")
	resp := checkStatus()
	Info.Println(resp)

	for _, com := range comments {
		d, err := json.Marshal(com)
		if err != nil {
			Fatal.Println(err)
		}
		resp := postDoc("evendex", string(d))
		Info.Println(resp)
	}

}
