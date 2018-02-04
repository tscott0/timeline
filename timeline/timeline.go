package timeline

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type Timeline struct {
	Events        []Event
	Options       Options
	Templates     *template.Template
	EarliestStart *time.Time
	LatestFinish  *time.Time

	// vis.js objects to inject into the template
	EventArray    template.JS
	OptionsObject template.JS
}

type Options struct {
	Start          *visJSTime `json:"start,omitempty"`
	End            *visJSTime `json:"end,omitempty"`
	FollowMouse    bool       `json:"followMouse,omitempty"`
	OverflowMethod string     `json:"overflowMethod,omitempty"`
}

type Event struct {
	ID      int        `json:"id"`
	Content string     `json:"content"`
	Start   *visJSTime `json:"start"`         // TODO: Convert to time.Time
	End     *visJSTime `json:"end,omitempty"` // TODO: Convert to time.Time
	Title   string     `json:"title"`
}

const padding float32 = 0.1

var homeTemplate *template.Template

func init() {
	templatesToParse := []string{"timeline/templates/home.html",
		"timeline/templates/timeline.js"}
	fmt.Println("Parsing templates...")
	for _, t := range templatesToParse {
		fmt.Printf("  %s\n", t)
	}

	homeTemplate = template.Must(template.ParseFiles(templatesToParse...))
}

func New() *Timeline {
	return &Timeline{
		Templates: homeTemplate,
	}
}

func (t *Timeline) WithOptions(o *Options) *Timeline {
	t.Options = *o
	return t
}

func (t *Timeline) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	eventData, err := json.MarshalIndent(t.Events, "", "   ")
	if err != nil {
		serveError(w, r, http.StatusInternalServerError, "Failed to marshal event data")
		return
	}

	t.EventArray = template.JS(eventData)

	optionsData, err := json.MarshalIndent(t.Options, "", "   ")
	if err != nil {
		serveError(w, r, http.StatusInternalServerError, "Failed to marshal event options")
		return
	}

	t.OptionsObject = template.JS(optionsData)

	w.WriteHeader(http.StatusOK)
	t.Templates.ExecuteTemplate(w, "home.html", t)
}

func serveError(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.WriteHeader(status)

	w.Write([]byte(fmt.Sprintf("%d %s - %s", status, http.StatusText(status), message)))
}

func (t *Timeline) AddEvent(content string, start time.Time) *Timeline {
	e := Event{
		ID:      len(t.Events),
		Content: content,
		Start:   &visJSTime{start},
		Title:   content,
	}
	t.Events = append(t.Events, e)

	t.updateBoundaries(&e)

	return t
}

func (t *Timeline) AddEventWithEnd(content string, start time.Time, end time.Time) *Timeline {
	e := Event{
		ID:      len(t.Events),
		Content: content,
		Start:   &visJSTime{start},
		End:     &visJSTime{end},
		Title:   content,
	}
	t.Events = append(t.Events, e)

	t.updateBoundaries(&e)

	return t
}

func (t *Timeline) updateBoundaries(e *Event) {
	if t.EarliestStart == nil || e.Start.Before(*t.EarliestStart) {
		t.EarliestStart = &e.Start.Time
	}

	if t.LatestFinish == nil || e.Start.After(*t.LatestFinish) {
		t.LatestFinish = &e.Start.Time
	}

	span := t.LatestFinish.Sub(*t.EarliestStart)
	paddingDuration := time.Duration(float32(span.Nanoseconds()) * padding)

	t.Options.Start = &visJSTime{t.EarliestStart.Add(paddingDuration * -1)}
	t.Options.End = &visJSTime{t.LatestFinish.Add(paddingDuration)}
}

func (t *Timeline) Reset() {
	*t = Timeline{
		Templates: t.Templates, // Parsed templates can persist after resets
	}
}
