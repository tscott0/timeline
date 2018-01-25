package timeline

import (
	"encoding/json"
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

func New() *Timeline {
	return &Timeline{
		Templates: template.Must(template.ParseFiles("timeline/templates/home.html",
			"timeline/templates/timeline.js")),
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

func (t *Timeline) AddEvent(content string, start time.Time) *Timeline {
	e := Event{ID: len(t.Events), Content: content, Start: visJSTime{start}}

	if t.EarliestStart == nil || start.Before(*t.EarliestStart) {
		t.EarliestStart = &start
	}

	// TODO: Work out a sensible margin before and after the events, maybe a %
	t.Options.Start = visJSTime{*t.EarliestStart}

	if t.LatestFinish == nil || start.After(*t.LatestFinish) {
		t.LatestFinish = &start
	}

	// TODO: Work out a sensible margin before and after the events, maybe a %
	t.Options.End = visJSTime{*t.LatestFinish}

	t.Events = append(t.Events, e)

	return t
}

func (t *Timeline) Reset() {
	*t = Timeline{
		Templates: t.Templates, // Parsed templates can persist after resets
	}
}
