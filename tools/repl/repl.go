package main

import (
	"log"

	"github.com/StabbyCutyou/blunderbuss/boot"
	"github.com/StabbyCutyou/blunderbuss/models"
	"github.com/StabbyCutyou/blunderbuss/services"
	"github.com/StabbyCutyou/instructor"
)

type locator struct {
	els services.IEventLoggingService
}

func main() {
	// Start by wiring up all of our external dependencies
	bp, err := boot.Boot()
	if err != nil {
		log.Fatal(err)
	}

	l := locator{
		els: bp.EventService,
	}
	i := instructor.New()
	i.RegisterFinder("Events", l.findEvents)
	i.RegisterFinder("NewEvent", newEvent)
	err = i.REPL()
	if err != nil {
		log.Fatal(err)
	}
}

func (l *locator) findEvents(id string) (interface{}, error) {
	return l.els.FindEvents(&services.EventSearchParams{Application: id})
}

func newEvent(id string) (interface{}, error) {
	return &models.Event{}, nil
}
