package main

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type CommentResponse struct {
	Comment Comment `json:"comment"`
}

type LabelResponse struct {
	Label Label `json:"label"`
}

// Response sing Event template
type ResponseEvent struct {
	Event Event `json:"event"`
}

// Response Events template
type ResponseEvents struct {
	Count  int     `json:"count"`
	Events []Event `json:"events"`
}

// Response Events Paginated template
type ResponseEventsPage struct {
	Offset interface{} `json:"offset"`
	Count  int         `json:"count"`
	Events []Event     `json:"events"`
}

// Event holder
type Event struct {
	EventID int       `json:"event_id"`
	Name    string    `json:"name"`
	Created string    `json:"created"`
	Content []Content `json:"content"`
}

// Context of an event
type Content struct {
	ContentID int       `json:"content_id"`
	Created   string    `json:"created"`
	Version   int       `json:"version"`
	Status    string    `json:"status"`
	Title     string    `json:"title"`
	Tag       string    `json:"tag"`
	Comments  []Comment `json:"comments"`
}

// Comment of Context
type Comment struct {
	CommentID int     `json:"comment_id"`
	Version   int     `json:"version"`
	Created   string  `json:"created"`
	Modified  string  `json:"modified"`
	Body      string  `json:"body"`
	Labels    []Label `json:"labels"`
}

// Label of Comments
type Label struct {
	LabelID int    `json:"label_id"`
	Created string `json:"created"`
	Value   string `json:"label"`
}

/*
	create functions:
		create single event - createEvent(db, name) returns created eventID
		create content (topics) for a single event - createContent(db, eventID, title) returns created contentID and unique tag
			_create tag - createTag() returns unique tag for content
		create a comment inside of content - createComment(db, contentID or tag, body) returns created commentID
		create a label for a comment - createLabel(db, contentID or tag, commentID, label) returns created labelID
*/
func createEvent(db *sql.DB, name string) (event interface{}, retErr error) {
	eventID := 0
	nowISO := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO events (name, created) VALUES($1, $2) RETURNING event_id;`

	err := db.QueryRow(query, name, nowISO).Scan(&eventID)
	if err != nil {
		retErr = err
		return
	}

	event, retErr = getEvent(db, eventID)
	return
}

func createContent(db *sql.DB, eventID int, title string) (resp interface{}, retErr error) {
	contentID := 0
	nowISO := time.Now().UTC().Format(time.RFC3339)
	query := `INSERT INTO content (event_id, created, status, version, title, tag) VALUES($1, $2, $3, $4, $5, $6) RETURNING content_id;`

	tag := string(createTag())
	err := db.QueryRow(query, eventID, nowISO, "Created", 1, title, tag).Scan(&contentID)
	if err != nil {
		retErr = err
		return
	}

	resp, retErr = getContent(db, eventID, contentID)
	return
}

func createTag() (tag []rune) {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	rand.Seed(time.Now().UnixNano())

	tag = make([]rune, 7)
	for i := range tag {
		if i == 3 {
			tag[i] = '-'
			continue
		}
		tag[i] = runes[rand.Intn(len(runes))]
	}
	return
}

func createComment(db *sql.DB, eventID int, arg interface{}, body string) (resp interface{}, retErr error) {
	contentID := 0
	commentID := 0
	nowISO := time.Now().UTC().Format(time.RFC3339)
	switch a := arg.(type) {
	case string:
		for _, c := range a {
			if _, err := strconv.Atoi(string(c)); err == nil {
				Fatal.Println("ERROR - CREATE COMMENT - string supplied for tag but integer found in tag. tag consist of non digits")
			}
		}
		query := fmt.Sprintf("SELECT (content_id) FROM content WHERE tag='%s';", a)
		err := db.QueryRow(query).Scan(&contentID)
		if err != nil {
			retErr = err
			return
		}

		query = `INSERT INTO comment (content_id, version, created, modified, body) VALUES($1, $2, $3, $4, $5) RETURNING comment_id;`

		err = db.QueryRow(query, contentID, 1, nowISO, nowISO, body).Scan(&commentID)
		if err != nil {
			retErr = err
			return
		}
	case int:
		query := `INSERT INTO comment (content_id, version, created, modified, body) VALUES($1, $2, $3, $4, $5) RETURNING comment_id;`
		err := db.QueryRow(query, a, 1, nowISO, nowISO, body).Scan(&commentID)
		if err != nil {
			retErr = err
			return
		}
		contentID = a
	default:
		retErr = errors.New("ERROR - CREATE COMMENT - no valid input for var 'arg'; contentID (int) or tag (string)")
		return
	}
	resp, retErr = getContent(db, eventID, contentID)
	increaseContentVersion(db, contentID)
	return
}

func createLabel(db *sql.DB, arg interface{}, commentID int, label string) (labelID int, retErr error) {
	contentID := 0
	nowISO := time.Now().UTC().Format(time.RFC3339)
	switch a := arg.(type) {
	case string:
		for _, c := range a {
			if _, err := strconv.Atoi(string(c)); err == nil {
				Fatal.Println("ERROR - CREATE COMMENT - string supplied for tag but integer found in tag. tag consist of non digits")
			}
		}
		query := fmt.Sprintf("SELECT (content_id) FROM content WHERE tag='%s';", a)
		err := db.QueryRow(query).Scan(&contentID)
		if err != nil {
			retErr = err
			return
		}

		query = `INSERT INTO labels (content_id, comment_id, created, label) VALUES($1, $2, $3, $4) RETURNING label_id`
		err = db.QueryRow(query, contentID, commentID, nowISO, label).Scan(&labelID)
		if err != nil {
			retErr = err
			return
		}

	case int:
		query := `INSERT INTO labels (content_id, comment_id, created, label) VALUES($1, $2, $3, $4) RETURNING label_id`
		contentID = a
		err := db.QueryRow(query, contentID, commentID, nowISO, label).Scan(&labelID)
		if err != nil {
			retErr = err
			return
		}
	default:
		retErr = errors.New("ERROR - CREATE LABEL - no valid input for var 'arg'; contentID (int) or tag (string)")
		return
	}
	increaseContentVersion(db, contentID)
	increaseCommentVersion(db, commentID)
	return
}

/*
	update functions:
		update a comment - updateComment(db, commentID, bodt) returns contentID
*/
func updateComment(db *sql.DB, commentID int, body string) (contentID int) {
	nowISO := time.Now().UTC().Format(time.RFC3339)
	query := `UPDATE comment SET body=$1, modified=$2 WHERE comment_id=$3 RETURNING content_id`

	err := db.QueryRow(query, body, nowISO, commentID).Scan(&contentID)
	if err != nil {
		Fatal.Println(err)
	}
	increaseCommentVersion(db, commentID)
	increaseContentVersion(db, contentID)
	return
}

/*
	delete functions:
		delete an entire event and all connected content - deleteEvent(db, eventID)
		delete entire content and all connected content - deleteContent(db, contentID)
		delete a comment and all labels associated - deleteComment(db, commentID)s
		delete a label from a comment - deleteLabel(db, labelID)
*/
func deleteLabel(db *sql.DB, labelID int) (retErr error) {
	row := ""

	query := fmt.Sprintf("DELETE FROM labels WHERE label_id='%d' RETURNING (content_id, comment_id)", labelID)
	err := db.QueryRow(query).Scan(&row)
	if err == nil {
		contentID, err := strconv.Atoi(string(row[1]))
		if err != nil {
			retErr = err
			return
		}
		commentID, err := strconv.Atoi(string(row[3]))
		if err != nil {
			retErr = err
			return
		}
		increaseContentVersion(db, contentID)
		increaseCommentVersion(db, commentID)
	}
	retErr = err
	return
}

func deleteEvent(db *sql.DB, eventID int) (err error) {
	query := fmt.Sprintf("DELETE FROM events WHERE event_id='%d' RETURNING event_id", eventID)
	_, err = db.Query(query)
	return
}

func deleteContent(db *sql.DB, contentID int) (err error) {
	query := fmt.Sprintf("DELETE FROM content WHERE content_id='%d' RETURNING event_id", contentID)
	_, err = db.Query(query)
	return
}

func deleteComment(db *sql.DB, commentID int) (err error) {
	ret := ""
	query := fmt.Sprintf("DELETE FROM comment WHERE comment_id='%d' RETURNING content_id", commentID)
	err = db.QueryRow(query).Scan(&ret)
	contentID, err := strconv.Atoi(string(ret))
	if err != nil {
		Error.Println(err)
		return
	}
	increaseContentVersion(db, contentID)
	return
}

/*
	get functions:
		get single event - getEvent(db, eventID) returns json resp in byte array
		get pages of 50 events - getEvents(db) returns json resp in byte array
		get a single content of an event - getContent(db, contentID or tag) returns json resp in byte array
*/
func getEvent(db *sql.DB, eventID int) (response interface{}, retErr error) {
	var content []Content

	name := ""
	created := ""

	contentRows, err := db.Query(fmt.Sprintf("SELECT * FROM content WHERE event_id='%d' ORDER BY content_id ASC;", eventID))
	if err != nil {
		retErr = err
		return
	}
	defer contentRows.Close()

	for contentRows.Next() {
		var comments []Comment
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
			var labels []Label
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

			for labelRows.Next() {
				labelID := 0
				labelCreated := ""
				label := ""
				err := labelRows.Scan(&labelID, &contentID, &commentID, &labelCreated, &label)
				if err != nil {
					retErr = err
					return
				}

				labels = append(labels, Label{labelID, labelCreated, label})
			}

			comments = append(comments, Comment{commentID, commentVersion, commentCreated, commentModified, body, labels})
		}

		content = append(content, Content{contentID, contentCreated, contentVersion, status, title, tag, comments})
	}

	err = db.QueryRow(fmt.Sprintf("SELECT * FROM events WHERE event_id=%d;", eventID)).Scan(&eventID, &name, &created)
	if err != nil {
		retErr = err
		return
	}

	response = ResponseEvent{Event{eventID, name, created, content}}
	return
}

func getEvents(db *sql.DB, offset int) (response interface{}, retErr error) {
	var events []Event
	counter := 0

	eventRows, err := db.Query(fmt.Sprintf("SELECT * FROM events ORDER BY event_id ASC LIMIT 50 OFFSET %d;", offset))
	if err != nil {
		retErr = err
		return
	}
	defer eventRows.Close()

	for eventRows.Next() {
		var content []Content

		eventID := 0
		name := ""
		created := ""
		err := eventRows.Scan(&eventID, &name, &created)
		if err != nil {
			retErr = err
			return
		}

		contentRows, err := db.Query(fmt.Sprintf("SELECT * FROM content WHERE event_id='%d' ORDER BY content_id ASC;", eventID))
		if err != nil {
			retErr = err
			return
		}
		defer contentRows.Close()

		for contentRows.Next() {
			var comments []Comment
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
				var labels []Label
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

				for labelRows.Next() {
					labelID := 0
					labelCreated := ""
					label := ""
					err := labelRows.Scan(&labelID, &contentID, &commentID, &labelCreated, &label)
					if err != nil {
						retErr = err
						return
					}

					labels = append(labels, Label{labelID, labelCreated, label})
				}

				comments = append(comments, Comment{commentID, commentVersion, commentCreated, commentModified, body, labels})
			}

			content = append(content, Content{contentID, contentCreated, contentVersion, status, title, tag, comments})
		}

		events = append(events, Event{eventID, name, created, content})
		counter++
	}
	err = eventRows.Err()
	if err != nil {
		retErr = err
		return
	}

	if counter < 50 {
		response = ResponseEvents{counter, events}
	} else {
		response = ResponseEventsPage{(offset + 50), counter, events}
	}
	return
}

func getContent(db *sql.DB, eventID int, arg interface{}) (response interface{}, retErr error) {
	var comments []Comment
	name := ""
	created := ""
	contentID := 0
	contentCreated := ""
	status := ""
	contentVersion := 0
	title := ""
	tag := ""

	switch a := arg.(type) {
	case string:
		for _, c := range a {
			if _, err := strconv.Atoi(string(c)); err == nil {
				Fatal.Println("ERROR - GET CONTENT - string supplied for tag but integer found in tag. tag consist of non digits")
			}
		}

		query := fmt.Sprintf("SELECT (content_id) FROM content WHERE event_id='%d' AND tag='%s' ORDER BY content_id ASC;", eventID, a)
		err := db.QueryRow(query).Scan(&contentID)
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
			var labels []Label
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

			for labelRows.Next() {
				labelID := 0
				labelCreated := ""
				label := ""
				err := labelRows.Scan(&labelID, &contentID, &commentID, &labelCreated, &label)
				if err != nil {
					retErr = err
					return
				}

				labels = append(labels, Label{labelID, labelCreated, label})
			}

			comments = append(comments, Comment{commentID, commentVersion, commentCreated, commentModified, body, labels})
		}

		err = db.QueryRow(fmt.Sprintf("SELECT * FROM content WHERE event_id='%d' AND content_id='%d'", eventID, contentID)).Scan(&contentID, &eventID, &contentCreated, &status, &contentVersion, &title, &tag)
		if err != nil {
			retErr = err
			return
		}

		content := Content{contentID, contentCreated, contentVersion, status, title, tag, comments}

		err = db.QueryRow(fmt.Sprintf("SELECT * FROM events WHERE event_id=%d;", eventID)).Scan(&eventID, &name, &created)
		if err != nil {
			retErr = err
			return
		}

		response = ResponseEvent{Event{eventID, name, created, []Content{content}}}

	case int:
		contentID = a
		commentRows, err := db.Query(fmt.Sprintf("SELECT * FROM comment WHERE content_id='%d'", contentID))
		if err != nil {
			retErr = err
			return
		}
		defer commentRows.Close()

		for commentRows.Next() {
			var labels []Label
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

			for labelRows.Next() {
				labelID := 0
				labelCreated := ""
				label := ""
				err := labelRows.Scan(&labelID, &contentID, &commentID, &labelCreated, &label)
				if err != nil {
					retErr = err
					return
				}

				labels = append(labels, Label{labelID, labelCreated, label})
			}

			comments = append(comments, Comment{commentID, commentVersion, commentCreated, commentModified, body, labels})
		}

		err = db.QueryRow(fmt.Sprintf("SELECT * FROM content WHERE content_id='%d'", contentID)).Scan(&contentID, &eventID, &contentCreated, &status, &contentVersion, &title, &tag)
		if err != nil {
			retErr = err
			return
		}

		content := Content{contentID, contentCreated, contentVersion, status, title, tag, comments}

		err = db.QueryRow(fmt.Sprintf("SELECT * FROM events WHERE event_id=%d;", eventID)).Scan(&eventID, &name, &created)
		if err != nil {
			retErr = err
			return
		}

		response = ResponseEvent{Event{eventID, name, created, []Content{content}}}
	default:
		retErr = errors.New("ERROR - GET CONTENT - no valid input for var 'arg'; contentID (int) or tag (string)")
	}
	return
}

func getComment(db *sql.DB, contentID int, commentID int) (resp interface{}, retErr error) {
	var labels []Label
	commentVersion := 0
	commentCreated := ""
	commentModified := ""
	body := ""

	err := db.QueryRow(fmt.Sprintf("SELECT * FROM comment WHERE content_id='%d' AND comment_id='%d'", contentID, commentID)).Scan(&commentID, &contentID, &commentVersion, &commentCreated, &commentModified, &body)
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

	for labelRows.Next() {
		labelID := 0
		labelCreated := ""
		label := ""
		err := labelRows.Scan(&labelID, &contentID, &commentID, &labelCreated, &label)
		if err != nil {
			retErr = err
			return
		}

		labels = append(labels, Label{labelID, labelCreated, label})
	}

	comment := Comment{commentID, commentVersion, commentCreated, commentModified, body, labels}
	resp = CommentResponse{comment}
	return
}

func getLabel(db *sql.DB, commentID int, labelID int) (resp interface{}, retErr error) {
	labelCreated := ""
	label := ""
	contentID := 0

	err := db.QueryRow(fmt.Sprintf("SELECT * FROM labels WHERE label_id='%d' AND comment_id='%d'", labelID, commentID)).Scan(&labelID, &contentID, &commentID, &labelCreated, &label)
	if err != nil {
		retErr = err
		return
	}

	labels := Label{labelID, labelCreated, label}
	resp = LabelResponse{labels}
	return
}

/*
	tracking functions:
		increase content version - increaseContentVersion(db, contentID)
		increase comment version - increaseCommentversion(db, commentID)

*/
func increaseContentVersion(db *sql.DB, contentID int) {
	version := 0
	query := `SELECT version FROM content WHERE content_id=$1;`

	err := db.QueryRow(query, contentID).Scan(&version)
	if err != nil {
		Fatal.Println(err)
	}

	version++
	query = `UPDATE content SET version=$1 WHERE content_id=$2;`

	_, err = db.Exec(query, version, contentID)
	if err != nil {
		Fatal.Println(err)
	}
}

func increaseCommentVersion(db *sql.DB, commentID int) {
	version := 0
	query := `SELECT version FROM comment WHERE comment_id=$1;`

	err := db.QueryRow(query, commentID).Scan(&version)
	if err != nil {
		Fatal.Println(err)
	}

	version++
	query = `UPDATE comment SET version=$1 WHERE comment_id=$2;`

	_, err = db.Exec(query, version, commentID)
	if err != nil {
		Fatal.Println(err)
	}
}
