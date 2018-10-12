package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// IndexedComment flattened for ES
type IndexedComment struct {
	EventID         int        `json:"event_id"`
	Name            string     `json:"name"`
	EventCreated    string     `json:"event_created"`
	ContentID       int        `json:"content_id"`
	ContentCreated  string     `json:"content_created"`
	Tag             string     `json:"tag"`
	ContentVersion  int        `json:"content_version"`
	Title           string     `json:"title"`
	Status          string     `json:"status"`
	CommentID       int        `json:"comment_id"`
	CommentVersion  int        `json:"comment_version"`
	Body            string     `json:"body"`
	CommentCreated  string     `json:"comment_created"`
	CommentModified string     `json:"comment_modified"`
	Labels          []LabelMap `json:"labels"`
}

type LabelMap struct {
	LabelID int    `json:"label_id"`
	Created string `json:"created"`
	Label   string `json:"label"`
}

type IndexBody struct {
	IndexCont IndexContent `json:"index"`
}

type IndexContent struct {
	Index string `json:"_index"`
	Type  string `json:"_type"`
	ID    int    `json:"_id"`
}

func IndexComments(db *sql.DB) (retErr error) {
	f, err := os.OpenFile("CommentsBulk.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	comments, err := bulkComment(db)
	if err != nil {
		retErr = err
		return
	}

	for ind, com := range comments {
		var c interface{}
		c = com

		indexBody := IndexBody{IndexContent{"evendex", "doc", ind}}
		d, err := json.Marshal(indexBody)
		if err != nil {
			retErr = err
			return
		}
		f.WriteString(string(d) + "\n")

		d, err = json.Marshal(c)
		if err != nil {
			retErr = err
			return
		}
		f.WriteString(string(d) + "\n")
	}
	return
}

func bulkComment(db *sql.DB) (comments []IndexedComment, retErr error) {
	counter := 0

	eventRows, err := db.Query("SELECT * FROM events ORDER BY event_id")
	if err != nil {
		retErr = err
		return
	}
	defer eventRows.Close()

	for eventRows.Next() {
		eventID := 0
		name := ""
		created := ""
		err := eventRows.Scan(&eventID, &name, &created)
		if err != nil {
			retErr = err
			return
		}

		contentRows, err := db.Query(fmt.Sprintf("SELECT * FROM content WHERE event_id='%d'", eventID))
		if err != nil {
			retErr = err
			return
		}
		defer contentRows.Close()

		for contentRows.Next() {
			contentID := 0
			contentCreated := ""
			status := ""
			contentVersion := 0
			title := ""
			tag := ""
			err := contentRows.Scan(&contentID, &eventID, &contentCreated, &status, &contentVersion, &title, &tag)
			if err != nil {
				retErr = err
				return
			}

			commentRows, err := db.Query(fmt.Sprintf("SELECT * FROM comment WHERE content_id='%d'", contentID))
			if err != nil {
				retErr = err
				return
			}
			defer commentRows.Close()

			for commentRows.Next() {
				commentID := 0
				commentVersion := 0
				commentCreated := ""
				commentModified := ""
				body := ""
				err := commentRows.Scan(&commentID, &contentID, &commentVersion, &commentCreated, &commentModified, &body)
				if err != nil {
					retErr = err
					return
				}

				labelRows, err := db.Query(fmt.Sprintf("SELECT * FROM labels WHERE comment_id='%d'", commentID))
				if err != nil {
					retErr = err
					return
				}
				defer labelRows.Close()
				var labels []LabelMap
				for labelRows.Next() {
					labelID := 0
					labelCreated := ""
					label := ""
					err := labelRows.Scan(&labelID, &contentID, &commentID, &labelCreated, &label)
					if err != nil {
						retErr = err
						return
					}
					labels = append(labels, LabelMap{labelID, labelCreated, label})
				}
				comments = append(comments, IndexedComment{eventID, name, created, contentID, contentCreated, tag, contentVersion, title, status, commentID, commentVersion, body, commentCreated, commentModified, labels})
			}
		}
		counter++
	}
	err = eventRows.Err()
	if err != nil {
		retErr = err
		return
	}

	return
}
