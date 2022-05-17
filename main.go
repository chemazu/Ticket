package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}
type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}
func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	// unmarshal the request body to the newEvent struct
	json.Unmarshal(reqBody, &newEvent)

	// add the new event to the events array
	events = append(events, newEvent)
	// send a response back to the client
	w.WriteHeader(http.StatusCreated)
	// encode the response and send it back
	json.NewEncoder(w).Encode(newEvent)
}
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}
func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	fmt.Fprintf(w, "The event with ID %v will be updated", eventID)
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}
func main() {
	// creating the router
	router := mux.NewRouter().StrictSlash(true)
	// e.g
	// r.HandleFunc("/products", ProductsHandler)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/newevent", createEvent).Methods("POST")
	router.HandleFunc("/updateevent", updateEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
