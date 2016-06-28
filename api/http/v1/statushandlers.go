package httpv1

import (
	"encoding/json"
	"net/http"
)

func (h *HTTPApi) statusServer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
}
