package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Response struct {
	Status string
}

type LabelBody struct {
	Content    interface{} `json:"content"`
	CommentID  int         `json:"comment_id"`
	LabelValue string      `json:"label"`
}

type EventBody struct {
	Name string `json:"name"`
}

type CommentBody struct {
	Body string `json:"body"`
}

func apiCreateEvent(w http.ResponseWriter, r *http.Request) {

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	name := r.URL.Query().Get("name")
	if name == "" {
		log.Println("apiPostContent", err)
		http.Error(w, "ERROR - FAILED TO CREATE EVENT - NO NAME RECEIVED", 400)
	} else {
		resp, err := createEvent(db, name)
		if err != nil {
			log.Println("apiCreateEvent", err)
			http.Error(w, "ERROR - FAILED TO CREATE EVENT - BAD REQUEST", 400)
		} else {
			json.NewEncoder(w).Encode(resp)
		}
	}
}

func apiGetEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	eventID, err := strconv.Atoi(params["eventId"])
	if err != nil {
		log.Println("reqEvent", err)
		http.Error(w, fmt.Sprintf("ERROR - FAILED GET EVENT - %s", params["eventId"]), 400)
	}

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	event, err := getEvent(db, eventID)
	if err == nil {
		json.NewEncoder(w).Encode(event)
	} else {
		log.Println("reqEvent", err)
		http.Error(w, fmt.Sprintf("ERROR - FAILED GET EVENT - %s", params["eventId"]), 400)
	}

}

func apiGetEvents(w http.ResponseWriter, r *http.Request) {

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	events, err := getEvents(db, offset)
	if err == nil {
		json.NewEncoder(w).Encode(events)
	} else {
		log.Println("reqEvents", err)
		http.Error(w, "ERROR - FAILED GET EVENTS", 400)
	}

}

func apiGetEventContent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	eventID, err := strconv.Atoi(params["eventId"])
	if err != nil {
		log.Println("apiPostContent", err)
		http.Error(w, "ERROR - FAILED TO CREATE CONTENT - INVALID ID", 400)
	} else {
		title := r.URL.Query().Get("title")
		if title == "" {
			log.Println("apiPostContent", err)
			http.Error(w, "ERROR - FAILED TO CREATE CONTENT - NO TITLE RECEIVED", 400)
		} else {
			resp, err := createContent(db, eventID, title)
			if err != nil {
				log.Println("apiPostContent", err)
				http.Error(w, "ERROR - FAILED TO CREATE CONTENT - BAD REQUEST", 400)
			} else {
				json.NewEncoder(w).Encode(resp)
			}
		}
	}
}

func apiGetContent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var content interface{}

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	content, err = strconv.Atoi(params["contentId"])
	if err != nil {
		content = params["contentId"]
	}

	eventID, err := strconv.Atoi(params["eventId"])
	if err != nil {
		log.Println("apiGetContent", err)
		http.Error(w, "ERROR - FAILED GET CONTENT - INVALID EVENT ID", 400)
	}

	resp, err := getContent(db, eventID, content)
	if err == nil {
		json.NewEncoder(w).Encode(resp)
	} else {
		log.Println("apiGetContent", err)
		http.Error(w, "ERROR - FAILED GET CONTENT", 400)
	}

}

func apiGetComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = strconv.Atoi(params["eventId"])
	if err != nil {
		log.Println("apiGetComment", err)
		http.Error(w, "ERROR - FAILED TO GET COMMENT - INVALID ID", 400)
	} else {
		commentID, err := strconv.Atoi(params["commentId"])
		if err != nil {
			log.Println("apiGetComment", err)
			http.Error(w, "ERROR - FAILED TO GET COMMENT - INVALID COMMENT ID", 400)
		} else {
			contentID, err := strconv.Atoi(params["contentId"])
			if err != nil {
				log.Println("apiGetComment", err)
				http.Error(w, "ERROR - FAILED TO GET COMMENT - INVALID CONTENT ID", 400)
			} else {
				resp, err := getComment(db, contentID, commentID)
				if err != nil {
					log.Println("apiGetComment", err)
					http.Error(w, "ERROR - FAILED TO GET COMMENT - INVALID REQUEST", 400)
				} else {
					json.NewEncoder(w).Encode(resp)
				}
			}
		}
	}
}

func apiGetLabel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = strconv.Atoi(params["eventId"])
	if err != nil {
		log.Println("apiGetComment", err)
		http.Error(w, "ERROR - FAILED TO GET COMMENT - INVALID ID", 400)
	} else {
		commentID, err := strconv.Atoi(params["commentId"])
		if err != nil {
			log.Println("apiGetComment", err)
			http.Error(w, "ERROR - FAILED TO GET COMMENT - INVALID COMMENT ID", 400)
		} else {
			_, err := strconv.Atoi(params["contentId"])
			if err != nil {
				log.Println("apiGetComment", err)
				http.Error(w, "ERROR - FAILED TO GET COMMENT - INVALID CONTENT ID", 400)
			} else {
				labelID, err := strconv.Atoi(params["labelId"])
				if err != nil {
					log.Println("apiGetComment", err)
					http.Error(w, "ERROR - FAILED TO GET LABEL - INVALID LABEL ID", 400)
				} else {
					resp, err := getLabel(db, commentID, labelID)
					if err != nil {
						log.Println("apiGetComment", err)
						http.Error(w, "ERROR - FAILED TO GET COMMENT - INVALID REQUEST", 400)
					} else {
						json.NewEncoder(w).Encode(resp)
					}
				}
			}
		}
	}
}

func apiPostComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var content interface{}
	var commentBody CommentBody

	_ = json.NewDecoder(r.Body).Decode(&commentBody)

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	content, err = strconv.Atoi(params["contentId"])
	if err != nil {
		content = params["contentId"]
	}

	eventID, err := strconv.Atoi(params["eventId"])
	if err != nil {
		log.Println("apiPostcomment", err)
		http.Error(w, "ERROR - FAILED TO CREATE COMMENT - INVALID EVENT ID", 400)
	}

	if commentBody.Body == "" {
		log.Println("apiPostcomment", err)
		http.Error(w, "ERROR - FAILED TO CREATE COMMENT - NO BODY IN DATA", 400)
	} else {
		resp, err := createComment(db, eventID, content, commentBody.Body)
		if err == nil {
			json.NewEncoder(w).Encode(resp)
		} else {
			log.Println("apiPostcomment", err)
			http.Error(w, "ERROR - FAILED TO CREATE COMMENT - BAD REQUEST", 400)
		}
	}

}

func apiCreateLabel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var content interface{}

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	content, err = strconv.Atoi(params["contentId"])
	if err != nil {
		content = params["contentId"]
	}

	_, err = strconv.Atoi(params["eventId"])
	if err != nil {
		log.Println("apiPostContent", err)
		http.Error(w, "ERROR - FAILED TO CREATE LABEL - INVALID ID", 400)
	} else {
		commentID, err := strconv.Atoi(params["commentId"])
		if err != nil {
			log.Println("apiPostContent", err)
			http.Error(w, "ERROR - FAILED TO CREATE LABEL - INVALID ID", 400)
		} else {
			label := r.URL.Query().Get("label")
			if label == "" {
				log.Println("apiPostContent", err)
				http.Error(w, "ERROR - FAILED TO CREATE LABEL - NO LABEL RECEIVED", 400)
			} else {
				_, err := createLabel(db, content, commentID, label)
				if err != nil {
					log.Println("apiPostContent", err)
					http.Error(w, "ERROR - FAILED TO CREATE LABEL - BAD REQUEST", 400)
				} else {
					json.NewEncoder(w).Encode(Response{"SUCCESS"})
				}
			}
		}
	}

}

func apiDeleteLabel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	labelID, err := strconv.Atoi(params["labelId"])
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("ERROR - LABEL DELETE FAILED - LABEL_ID: %s - NOT CORRECT TYPE", params["id"]), 400)
	} else {
		err = deleteLabel(db, labelID)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("ERROR - LABEL DELETE FAILED - LABEL_ID: %d - FAILED TO FIND ID", labelID), 400)
		} else {
			json.NewEncoder(w).Encode(Response{fmt.Sprintf("SUCCESS - LABEL DELETE - LABEL_ID: %d", labelID)})
		}
	}
}

func apiIndexer(w http.ResponseWriter, r *http.Request) {
	db, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = IndexComments(db)
	if err != nil {
		log.Println("apiIndexer", err)
		http.Error(w, "ERROR - FAILED TO INDEX", 400)
	} else {
		json.NewEncoder(w).Encode(Response{"SUCCESS - INDEXED"})
	}
}

func main() {
	PORT := 8000
	LOGFILE := "Log_Views.log"
	router := mux.NewRouter()

	f, err := os.OpenFile(LOGFILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// events paginated route
	router.HandleFunc("/api/v1/events", apiGetEvents).Methods("GET") // get events 50 at a time

	// event routes
	router.HandleFunc("/api/v1/event", apiCreateEvent).Methods("GET")                                                                      // create an event
	router.HandleFunc("/api/v1/event/{eventId}", apiGetEvent).Methods("GET")                                                               // get event by ID
	router.HandleFunc("/api/v1/event/{eventId}/content", apiGetEventContent).Methods("GET")                                                // create content for an event
	router.HandleFunc("/api/v1/event/{eventId}/content/{contentId}", apiGetContent).Methods("GET")                                         // get content by ID or Tag
	router.HandleFunc("/api/v1/event/{eventId}/content/{contentId}/comment", apiPostComment).Methods("POST")                               // create comment
	router.HandleFunc("/api/v1/event/{eventId}/content/{contentId}/comment/{commentId}", apiGetComment).Methods("GET")                     // get comment
	router.HandleFunc("/api/v1/event/{eventId}/content/{contentId}/comment/{commentId}/label", apiCreateLabel).Methods("GET")              // create label for a comment
	router.HandleFunc("/api/v1/event/{eventId}/content/{contentId}/comment/{commentId}/label/{labelId}", apiGetLabel).Methods("GET")       // get label
	router.HandleFunc("/api/v1/event/{eventId}/content/{contentId}/comment/{commentId}/label/{labelId}", apiDeleteLabel).Methods("DELETE") // delete label

	router.HandleFunc("/api/v1/indexer", apiIndexer).Methods("GET")

	fmt.Println(fmt.Sprintf(" * Server hosted on http://127.0.0.1:%d", PORT))
	fmt.Println(fmt.Sprintf(" * Logging output -> %s", LOGFILE))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), router))
}
