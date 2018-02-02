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
	Start visJSTime `json:"start"`
	End   visJSTime `json:"end"`
}

type Event struct {
	ID      int       `json:"id"`
	Content string    `json:"content"`
	Start   visJSTime `json:"start"` // TODO: Convert to time.Time
}

var homeTemplate *template.Template

func init() {
	fmt.Println("Parsing home.html template")
	homeTemplate = template.Must(template.ParseFiles("timeline/templates/home.html",
		"timeline/templates/timeline.js"))
}

func New() *Timeline {
	return &Timeline{
		Templates: homeTemplate,
	}
}

func (t *Timeline) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	eventData, _ := json.MarshalIndent(t.Events, "", "   ") // TODO

	t.EventArray = template.JS(eventData)

	optionsData, _ := json.MarshalIndent(t.Options, "", "   ") // TODO

	t.OptionsObject = template.JS(optionsData)

	t.Templates.ExecuteTemplate(w, "home.html", t)
}

const padding float32 = 0.1

func (t *Timeline) AddEvent(content string, start time.Time) *Timeline {
	e := Event{ID: len(t.Events), Content: content, Start: visJSTime{start}}
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

	t.Options.Start = visJSTime{t.EarliestStart.Add(paddingDuration * -1)}
	t.Options.End = visJSTime{t.LatestFinish.Add(paddingDuration)}
}

func (t *Timeline) Reset() {
	*t = Timeline{
		Templates: t.Templates, // Parsed templates can persist after resets
	}
}
