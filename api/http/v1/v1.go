package httpv1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	//"github.com/facebookgo/grace/gracehttp"
	"github.com/StabbyCutyou/blunderbuss/models"
	"github.com/StabbyCutyou/blunderbuss/services"
	"github.com/gorilla/mux"
)

// HTTPApi represents the object used to govern http calls into the system
type HTTPApi struct {
	Config *Config
	Router *mux.Router
	Server *http.Server
}

// Config is the configuration for the HTTPApi struct
type Config struct {
	Version int
	Port    int
	Sha     string

	EventService *services.EventLoggingService
}

// New initializes a new http api
func New(config *Config) (*HTTPApi, error) {
	h := &HTTPApi{
		Config: config,
	}

	http.Handle("/", h.NewRouter())

	server := &http.Server{Handler: http.DefaultServeMux}
	h.Server = server
	return h, nil
}

// NewRouter builds the router for the server. We export this method for testing
func (h *HTTPApi) NewRouter() http.Handler {
	// Create a gorilla mux
	router := mux.NewRouter()
	// Strict Slash is documented here: http://www.gorillatoolkit.org/pkg/mux#Router.StrictSlash
	// It means we will try to match /path and /path/
	router.StrictSlash(true)

	v1Router := router.PathPrefix("/v1").Subrouter()
	statusRoutes := v1Router.PathPrefix("/status").Subrouter()
	// Define our health check under a modern route (status)
	statusRoutes.HandleFunc("/", h.statusServer).Methods("GET", "HEAD")

	v1Router.HandleFunc("/event", h.RecordEvent).Methods("PUT")
	v1Router.HandleFunc("/events", h.FindEvents).Methods("POST")

	// Serve our JSON Hyper Schemas as files directly
	// if we don't strip the prefix here, it will look for /api/http/v1/schemas/v1/schemas/{path}
	//s := http.StripPrefix("/v1/schemas/", http.FileServer(http.Dir("./api/http/v1/schemas/")))
	//v1Router.PathPrefix("/schemas").Handler(SchemaLogger(s))

	return router
}

// Listen is
func (h *HTTPApi) Listen() error {
	h.Server.Addr = fmt.Sprintf("0.0.0.0:%d", h.Config.Port)

	log.Printf("blunderbuss listening on %s\n", h.Server.Addr)
	log.Fatal(h.Server.ListenAndServe())

	return nil
}

func (h *HTTPApi) RecordEvent(w http.ResponseWriter, r *http.Request) {
	var e models.Event
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		resp, _ := json.Marshal(map[string]string{"error": err.Error(), "a": "a"})
		w.Write(resp)
		return
	}
	if err = r.Body.Close(); err != nil {
		w.WriteHeader(500)
		resp, _ := json.Marshal(map[string]string{"error": err.Error(), "b": "b"})
		w.Write(resp)
		return
	}

	if err = json.Unmarshal(b, &e); err != nil {
		w.WriteHeader(500)
		resp, _ := json.Marshal(map[string]string{"error": err.Error(), "c": "c"})
		w.Write(resp)
		return
	}
	if err = h.Config.EventService.LogEvent(&e); err != nil {
		w.WriteHeader(500)
		resp, _ := json.Marshal(map[string]string{"error": err.Error(), "d": "d"})
		w.Write(resp)
		return
	}
	w.WriteHeader(200)
	resp, _ := json.Marshal(map[string]string{"status": "ok"})
	w.Write(resp)
	return
}

func (h *HTTPApi) FindEvents(w http.ResponseWriter, r *http.Request) {
	var e services.EventSearchParams
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		resp, _ := json.Marshal(map[string]string{"error": err.Error(), "a": "a"})
		w.Write(resp)
		return
	}
	if err = r.Body.Close(); err != nil {
		w.WriteHeader(500)
		resp, _ := json.Marshal(map[string]string{"error": err.Error(), "b": "b"})
		w.Write(resp)
		return
	}

	if err = json.Unmarshal(b, &e); err != nil {
		w.WriteHeader(500)
		resp, _ := json.Marshal(map[string]string{"error": err.Error(), "c": "c"})
		w.Write(resp)
		return
	}
	evts, err := h.Config.EventService.FindEvents(&e)
	if err != nil {
		w.WriteHeader(500)
		resp, _ := json.Marshal(map[string]string{"error": err.Error(), "d": "d"})
		w.Write(resp)
		return
	}
	w.WriteHeader(200)
	resp, _ := json.Marshal(evts)
	w.Write(resp)
}
