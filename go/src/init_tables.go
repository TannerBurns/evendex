package main

import (
	"fmt"
	"log"
)

func checkFatal(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {

	db, err := connect()
	checkFatal(err)
	defer db.Close()

	err = db.Ping()
	checkFatal(err)

	fmt.Println("Connected!")

	// create tables
	iniTables(db)

	// create test event
	//eventID := createEvent(db, "Test Event 2")
	//fmt.Println(eventID)

	// create test topic for test event
	//contentID, tag := createContent(db, 1, "Test Topic 1")
	//fmt.Println(contentID)
	//fmt.Println(tag)

	// create test comment for Test Topic 1
	//commentID := createComment(db, 1, "New comment TEST")
	//fmt.Println(commentID)

	// update test comment for Test topic 1
	//contentID = updateComment(db, commentID, "UPDATING Commenting for Test Topic 1!! UPDATE")
	//fmt.Println(contentID)

	// test adding label to a comment
	//labelID := createLabel(db, 1, 1, "TEST")
	//fmt.Println(labelID)

	// test delete label
	//deleteLabel(db, labelID)

	// test get paginated events
	//resp := getEvents(db, 0)
	//log.Print(string(resp))

	// test get single event
	//resp := getEvent(db, 1)
	//log.Print(string(resp))

	// test get a single content
	//resp := getContent(db, 1)
	//log.Print(string(resp))

}
