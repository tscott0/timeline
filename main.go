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

	// Example: British monarchy since 1707
	// source of data https://en.wikipedia.org/w/index.php?title=List_of_British_monarchs&action=edit&section=4
	monarchy := timeline.New().
		WithOptions(&timeline.Options{
			FollowMouse:    true,
			OverflowMethod: "cap",
		}).
		AddEvent("Anne", time.Date(1707, time.January, 1, 0, 0, 0, 0, time.UTC)).
		AddEvent("George I", time.Date(1714, time.January, 1, 0, 0, 0, 0, time.UTC)).
		AddEvent("George II", time.Date(1727, time.January, 1, 0, 0, 0, 0, time.UTC)).
		AddEvent("George III", time.Date(1760, time.January, 1, 0, 0, 0, 0, time.UTC)).
		AddEvent("George IV", time.Date(1820, time.January, 1, 0, 0, 0, 0, time.UTC)).
		AddEvent("William IV", time.Date(1830, time.January, 1, 0, 0, 0, 0, time.UTC)).
		AddEvent("Victoria", time.Date(1837, time.January, 1, 0, 0, 0, 0, time.UTC)).
		AddEvent("Edward VII", time.Date(1901, time.January, 1, 0, 0, 0, 0, time.UTC)).
		AddEvent("George V", time.Date(1910, time.January, 1, 0, 0, 0, 0, time.UTC)).
		AddEvent("Edward VIII", time.Date(1936, time.January, 1, 0, 0, 0, 0, time.UTC)).
		AddEvent("George VI", time.Date(1936, time.January, 1, 0, 0, 0, 0, time.UTC))

	monarchy.AddEvent("Elizabeth II", time.Date(1952, time.January, 1, 0, 0, 0, 0, time.UTC))

	r.Handle("/monarchy", monarchy).Methods("GET")

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
