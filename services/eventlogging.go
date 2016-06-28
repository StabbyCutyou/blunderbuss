package services

import (
	"fmt"
	"time"

	"github.com/StabbyCutyou/blunderbuss/models"
	"github.com/jmoiron/sqlx"
)

// EventLoggingServiceConfig is
type EventLoggingServiceConfig struct {
	DB            *sqlx.DB
	MetricService IMetricLoggingService
}

// EventLoggingService is
type EventLoggingService struct {
	db            *sqlx.DB
	metricService IMetricLoggingService
}

// IEventLoggingService is
type IEventLoggingService interface {
	LogEvent(e *models.Event) error
	FindEvents(p *EventSearchParams) ([]models.Event, error)
}

// EventSearchParams is
type EventSearchParams struct {
	Application    string    `json:"application"`
	Type           string    `json:"type"`
	Message        string    `json:"message"`
	PartialMessage bool      `json:"partial_message"`
	Start          time.Time `json:"start"`
	End            time.Time `json:"end"`
}

const insertEventQuery = "INSERT INTO events (application, type, message, context, stack_trace, created_at) VALUES ($1, $2, $3, $4, $5, $6)"

// NewEventLoggingService is
func NewEventLoggingService(cfg *EventLoggingServiceConfig) (IEventLoggingService, error) {
	return &EventLoggingService{
		db:            cfg.DB,
		metricService: cfg.MetricService,
	}, nil
}

// LogEvent will
func (els *EventLoggingService) LogEvent(e *models.Event) error {
	if e == nil {
		return fmt.Errorf("Cannot log nil events")
	}

	args := make([]interface{}, 6)
	args[0] = e.Application
	args[1] = e.Type
	args[2] = e.Message
	args[3] = e.Context
	args[4] = e.StackTrace
	args[5] = e.CreatedAt

	if _, err := els.db.Exec(insertEventQuery, args...); err != nil {
		return err
	}
	if err := els.metricService.RecordEvent(e); err != nil {
		return err
	}
	return nil
}

// FindEvents will
func (els *EventLoggingService) FindEvents(p *EventSearchParams) ([]models.Event, error) {
	var evts []models.Event
	if p.Application == "" && p.Type == "" && p.Message == "" {
		return nil, fmt.Errorf("You must provide atleast one value to search")
	}
	query := "SELECT * FROM events WHERE "
	paramCount := 0
	needsAnd := false
	args := make([]interface{}, 0, 3)
	if p.Application != "" {
		paramCount++
		needsAnd = true
		query += fmt.Sprintf("application = $%d", paramCount)
		args = append(args, p.Application)
	}

	// Bucketting by time
	if !p.Start.IsZero() || !p.End.IsZero() {
		if needsAnd {
			query += " AND "
		}
		// It's going up by atleast one
		paramCount++
		if !p.Start.IsZero() && !p.End.IsZero() {
			// BETWEEN
			query += fmt.Sprintf("created_at BETWEEN $%d AND $%d", paramCount, paramCount+1)
			// We used 2 params, up it again
			paramCount++
		} else if !p.Start.IsZero() {
			// AFTER START
			query += fmt.Sprintf("created_at >= $%d", paramCount)
		} else if !p.End.IsZero() {
			// BEFORE END
			query += fmt.Sprintf("created_at <= $%d", paramCount)
		}
		needsAnd = true
	}

	if p.Type != "" {
		paramCount++
		if needsAnd {
			query += " AND "
		} else {
			needsAnd = true
		}
		query += fmt.Sprintf("type = $%d", paramCount)
		args = append(args, p.Type)
	}

	if p.Message != "" {
		paramCount++
		if p.PartialMessage {
			query += fmt.Sprintf("type LIKE $" + fmt.Sprintf("%%%d%%", paramCount))
		} else {
			query += fmt.Sprintf("type = $%d", paramCount)
		}
		args = append(args, p.Message)
	}
	err := els.db.Select(&evts, query, args...)
	if evts == nil {
		evts = make([]models.Event, 0)
	}
	return evts, err
}
