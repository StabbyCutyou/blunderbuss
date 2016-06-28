// Package middleware is where http request wrapping middleware for the service
// should be defined. There is a main Middleware Manager where you should register
// middleware to be called in a chain. You can also define any dependencies you need
// the middleware stack to have on the manager, so that they can be injected cleanly
package middleware

import (
	"net/http"

	"github.com/StabbyCutyou/blunderbuss/services"
)

// Manager is
type Manager struct {
	StatsClient services.StatsdClient
}

// Change th to accept in any of the dependencies, and assign it to a property
// on the Manager

// NewManager is
func NewManager(s services.StatsdClient) (*Manager, error) {
	return &Manager{StatsClient: s}, nil
}

// RequestHandler is
type RequestHandler func(http.ResponseWriter, *http.Request)

// Run will invoke the middleware in whatever order we decide
func (m *Manager) Run(rh RequestHandler) RequestHandler {
	// We could construct any kind of advanced, conditional middleware chaining
	// we might need from here
	return rh
}
