package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/tscott0/timeline/timeline"
)

var t2 *timeline.Timeline
var t3 *timeline.Timeline

func main() {
	r := mux.NewRouter()

	// t1: Basic timeline
	t1 := timeline.New().
		AddEvent("Item 1", time.Date(2014, time.April, 20, 12, 0, 0, 0, time.UTC)).
		AddEvent("Item 2", time.Date(2014, time.April, 20, 12, 13, 0, 0, time.UTC)).
		AddEvent("Item 3", time.Date(2014, time.April, 20, 12, 7, 0, 0, time.UTC))

	t1.AddEvent("Item 4", time.Date(2014, time.April, 20, 12, 30, 0, 0, time.UTC))

	r.Handle("/t1", t1).Methods("GET")

	// t2: Timeline wrapped in some middleware
	t2 = timeline.New()
	r.Handle("/t2/{id:[0-9A-Fa-f]{6}}", SimpleMiddleware(t2)).
		Methods("GET")

	// t3: Timeline wrapped in some more complex middleware
	t3 = timeline.New()
	r.Handle("/t3", ComplexMiddleware(t3)).
		Methods("GET")

	http.ListenAndServe(":8080", r)
}

func SimpleMiddleware(handler http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {

		// Timelines can persist between requests or be manually reset
		// which removes any Events that have been added
		t2.Reset()

		// Perform any routing and application-specific logic before
		// adding events to the timeline
		vars := mux.Vars(r)
		id := strings.ToUpper(vars["id"])
		fmt.Printf("Middleware ID: %v\n", id)

		t2.AddEvent("Item 1", time.Date(2014, time.April, 20, 12, 0, 0, 0, time.UTC))
		t2.AddEvent("Item 2", time.Date(2014, time.April, 20, 12, 13, 0, 0, time.UTC))
		t2.AddEvent("Item 3", time.Date(2014, time.April, 20, 12, 7, 0, 0, time.UTC))
		t2.AddEvent("Item 4", time.Date(2014, time.April, 20, 12, 30, 0, 0, time.UTC))
		t2.AddEvent("Item 5", time.Date(2014, time.April, 20, 12, 50, 0, 0, time.UTC))

		t2.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func ComplexMiddleware(t *timeline.Timeline) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		t.AddEvent(":D", time.Now()).ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}
