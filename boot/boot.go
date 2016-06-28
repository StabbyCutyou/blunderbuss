// Package boot provides a simple Boot() function where app-boot code would go
// you can then call it from main, or from any tools or other places where you
// want to "boot" the app with all it's configuration and dependencies
package boot

import (
	"github.com/StabbyCutyou/blunderbuss/api/http/v1"
	"github.com/StabbyCutyou/blunderbuss/config"
	"github.com/StabbyCutyou/blunderbuss/services"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Payload contains all the artifacts of a successfully booted system
type Payload struct {
	MetricService *services.MetricLoggingService
	EventService  *services.EventLoggingService
	HTTPServer    *httpv1.HTTPApi
}

// Boot will boot the application, and return an error if something went wrong
func Boot() (*Payload, error) {

	// There are some ordering considerations for initializing the various packages.
	// You need to initialize global configs because those are used to configure the other
	// packages. Then you need to configure logging because if any package fails to initialize
	// they need to at *least* log to the standard logger. Then you have to initialize
	// monitoring next because you want to report any fatal errors to the application monitoring
	// client.

	globalCfg, err := config.ReadFromEnv()
	if err != nil {
		return nil, err
	}

	db, err := openDB(globalCfg)
	if err != nil {
		return nil, err
	}

	statsd := services.NewStatsdClient(&services.StatsdConfig{
		Prefix:        globalCfg.StatsdPrefix,
		Address:       globalCfg.StatsdAddress,
		FlushInterval: int(globalCfg.StatsdInterval.Seconds()),
		Type:          "statsd",
		Enabled:       false,
	})
	if err != nil {
		return nil, err
	}

	metricService, err := services.NewMetricLoggingService(&services.MetricLoggingServiceConfig{
		Statsd: statsd,
	})
	if err != nil {
		return nil, err
	}

	eventService, err := services.NewEventLoggingService(&services.EventLoggingServiceConfig{
		DB:            db,
		MetricService: metricService,
	})
	if err != nil {
		return nil, err
	}

	httpServer, err := httpv1.New(&httpv1.Config{
		Version:      globalCfg.HTTPApiVersion,
		Port:         globalCfg.HTTPPort,
		Sha:          "",
		EventService: eventService,
	})
	if err != nil {
		return nil, err
	}
	return &Payload{
		EventService:  eventService,
		MetricService: metricService,
		HTTPServer:    httpServer,
	}, nil
}

func openDB(cfg *config.Config) (*sqlx.DB, error) {
	return sqlx.Open("postgres", cfg.DBConnString)
}
