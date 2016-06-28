package services

import "github.com/StabbyCutyou/blunderbuss/models"

// MetricLoggingServiceConfig is
type MetricLoggingServiceConfig struct {
	Statsd StatsdClient
}

// MetricLoggingService is
type MetricLoggingService struct {
	statsd StatsdClient
}

// NewMetricLoggingService is
func NewMetricLoggingService(cfg *MetricLoggingServiceConfig) (*MetricLoggingService, error) {
	return &MetricLoggingService{
		statsd: cfg.Statsd,
	}, nil
}

// RecordEvent is
func (m *MetricLoggingService) RecordEvent(e *models.Event) error {
	return m.statsd.Incr(eventToCountKey(e), 1)
}

func eventToCountKey(e *models.Event) string {
	return e.Application + "." + e.Type + "." + e.Message
}

func eventToTimeSeriesKey(e *models.Event) string {
	return eventToCountKey(e) + "format the date here"
}
