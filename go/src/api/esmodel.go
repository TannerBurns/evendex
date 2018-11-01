package main

import (
	"encoding/json"
	"fmt"
)

// IndexedComment flattened for ES
type IndexedComment struct {
	EventID         int    `json:"event_id"`
	Name            string `json:"name"`
	EventCreated    string `json:"event_created"`
	ContentID       int    `json:"content_id"`
	ContentCreated  string `json:"content_created"`
	ContentModified string `json:"content_modified"`
	Tag             string `json:"tag"`
	ContentVersion  int    `json:"content_version"`
	Title           string `json:"title"`
	Status          string `json:"status"`
	CommentID       int    `json:"comment_id"`
	CommentVersion  int    `json:"comment_version"`
	Body            string `json:"body"`
	CommentCreated  string `json:"comment_created"`
	CommentModified string `json:"comment_modified"`
}

func addDoc(eventID int, contentID int, commentID int) {
	var (
		eventCreated    string
		name            string
		contentCreated  string
		contentModified string
		tag             string
		contentVersion  int
		title           string
		status          string
		commentVersion  int
		body            string
		commentCreated  string
		commentModified string
	)

	db, err := connect()
	if err != nil {
		Fatal.Println(err)
	}
	defer db.Close()

	// get all event data
	row, err := db.Query("SELECT created, name FROM events WHERE event_id=$1", eventID)
	if err != nil {
		Error.Printf("ERROR - esmodel - addDoc - Failed to query event - %s", err)
		return
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&eventCreated, &name)
		if err != nil {
			Error.Printf("ERROR - esmodel - addDoc - Failed to scan event - %s", err)
			return
		}
	}

	// get all content data
	row, err = db.Query("SELECT created, modified, status, version, title, tag FROM content WHERE content_id=$1", contentID)
	if err != nil {
		Error.Printf("ERROR - esmodel - addDoc - Failed to query content - %s", err)
		return
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&contentCreated, &contentModified, &status, &contentVersion, &title, &tag)
		if err != nil {
			Error.Printf("ERROR - esmodel - addDoc - Failed to scan content - %s", err)
			return
		}
	}

	// get all comment data
	row, err = db.Query("SELECT version, created, modified, body FROM comment WHERE comment_id=$1", commentID)
	if err != nil {
		Error.Printf("ERROR - esmodel - addDoc - Failed to query comment - %s", err)
		return
	}
	defer row.Close()
	for row.Next() {
		err = row.Scan(&commentVersion, &commentCreated, &commentModified, &body)
		if err != nil {
			Error.Printf("ERROR - esmodel - addDoc - Failed to scan comment - %s", err)
			return
		}
	}

	indComment := IndexedComment{eventID, name, eventCreated, contentID, contentCreated, contentModified, tag, contentVersion, title, status, commentID, commentVersion, body, commentCreated, commentModified}

	dat, err := json.Marshal(indComment)
	if err != nil {
		Error.Printf("ERROR - esmodel - addDoc - Failed to Marshal comment - %s", err)
		return
	}

	resp := postDoc("evendex", string(dat))
	var respMap map[string]interface{}
	err = json.Unmarshal([]byte(resp), &respMap)
	if err != nil {
		Error.Printf("ERROR - esmodel - addDoc - Failed to unMarshal response - %s", err)
		return
	}

	esID := respMap["_id"]
	indexID := 0
	query := `INSERT INTO indexes (es_id, event_id, content_id, comment_id) VALUES($1, $2, $3, $4) RETURNING indexes_id`
	err = db.QueryRow(query, esID, eventID, contentID, commentID).Scan(&indexID)
	if err != nil {
		Error.Printf("ERROR - esmodel - addDoc - Failed to insert into indexes table - %s", err)
		return
	}

}

func updateEsStatus(contentID int, status string) {
	db, err := connect()
	if err != nil {
		Fatal.Println(err)
	}
	defer db.Close()

	esquery := []byte(fmt.Sprintf(`{"script": {"source": "ctx._source.status='%s';", "lang": "painless"}}`, status))
	query := fmt.Sprintf("SELECT es_id FROM indexes WHERE content_id=%d", contentID)

	rows, err := db.Query(query)
	if err != nil {
		Error.Printf("ERROR - esmodel - updateEsStatus - Failed to update- %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		esID := ""
		rows.Scan(&esID)
		_ = updateIndex("evendex", esID, esquery)
	}
}

func updateEsContentModified(contentID int, modified string) {
	db, err := connect()
	if err != nil {
		Fatal.Println(err)
	}
	defer db.Close()

	esquery := []byte(fmt.Sprintf(`{"script": {"source": "ctx._source.content_modified='%s';", "lang": "painless"}}`, modified))
	query := fmt.Sprintf("SELECT es_id FROM indexes WHERE content_id=%d", contentID)

	rows, err := db.Query(query)
	if err != nil {
		Error.Printf("ERROR - esmodel - updateEsStatus - Failed to update- %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		esID := ""
		rows.Scan(&esID)
		_ = updateIndex("evendex", esID, esquery)
	}
}

func updateEsContentVersion(contentID int, version int) {
	db, err := connect()
	if err != nil {
		Fatal.Println(err)
	}
	defer db.Close()

	esquery := []byte(fmt.Sprintf(`{"script": {"source": "ctx._source.content_version=%d;", "lang": "painless"}}`, version))
	query := fmt.Sprintf("SELECT es_id FROM indexes WHERE content_id=%d", contentID)

	rows, err := db.Query(query)
	if err != nil {
		Error.Printf("ERROR - esmodel - updateEsStatus - Failed to update- %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		esID := ""
		rows.Scan(&esID)
		_ = updateIndex("evendex", esID, esquery)
	}
}

func updateEsCommentVersion(commentID int, version int) {
	db, err := connect()
	if err != nil {
		Fatal.Println(err)
	}
	defer db.Close()

	esquery := []byte(fmt.Sprintf(`{"script": {"source": "ctx._source.comment_version=%d;", "lang": "painless"}}`, version))
	query := fmt.Sprintf("SELECT es_id FROM indexes WHERE comment_id=%d", commentID)

	rows, err := db.Query(query)
	if err != nil {
		Error.Printf("ERROR - esmodel - updateEsStatus - Failed to update- %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		esID := ""
		rows.Scan(&esID)
		_ = updateIndex("evendex", esID, esquery)
	}
}

func updateEsCommentModified(commentID int, modified string) {
	db, err := connect()
	if err != nil {
		Fatal.Println(err)
	}
	defer db.Close()

	esquery := []byte(fmt.Sprintf(`{"script": {"source": "ctx._source.comment_modified='%s';", "lang": "painless"}}`, modified))
	query := fmt.Sprintf("SELECT es_id FROM indexes WHERE comment_id=%d", commentID)

	rows, err := db.Query(query)
	if err != nil {
		Error.Printf("ERROR - esmodel - updateEsStatus - Failed to update- %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		esID := ""
		rows.Scan(&esID)
		_ = updateIndex("evendex", esID, esquery)
	}
}
