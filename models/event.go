package models

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx/types"
)

// Event is an instance of a thing that happened
type Event struct {
	Application string         `db:"application"`
	Type        string         `db:"type"`
	Message     string         `db:"message"`
	Context     types.JSONText `db:"context"`
	StackTrace  string         `db:"stack_trace"`
	CreatedAt   time.Time      `db:"created_at"`
}

type eventScaffold struct {
	Application string                 `json:"application"`
	Type        string                 `json:"type"`
	Message     string                 `json:"message"`
	Context     map[string]interface{} `json:"context"`
	StackTrace  string                 `json:"stack_trace"`
	CreatedAt   int64                  `json:"created_at"`
}

// UnmarshalJSON is a custom unmarshaller
func (e *Event) UnmarshalJSON(b []byte) error {
	es := eventScaffold{}
	if err := json.Unmarshal(b, &es); err != nil {
		return err
	}
	ctxt, err := json.Marshal(es.Context)
	if err != nil {
		return err
	}
	e.Application = es.Application
	e.Type = es.Type
	e.Message = es.Message
	e.Context = ctxt
	e.StackTrace = es.StackTrace
	e.CreatedAt = time.Unix(es.CreatedAt, 0) // no nano sec at this time
	return nil
}

// MarshalJSON is a custom marshaller
func (e *Event) MarshalJSON() ([]byte, error) {
	ctxt := make(map[string]interface{})
	if err := json.Unmarshal(e.Context, &ctxt); err != nil {
		return nil, err
	}
	es := eventScaffold{
		Application: e.Application,
		Type:        e.Type,
		Message:     e.Message,
		Context:     ctxt,
		StackTrace:  e.StackTrace,
		CreatedAt:   e.CreatedAt.Unix(),
	}
	return json.Marshal(es)
}
