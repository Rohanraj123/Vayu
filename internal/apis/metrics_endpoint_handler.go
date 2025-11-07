package apis

import "net/http"

// healthyHandler
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}

// readinessHandler
func ReadyzHandler(w http.ResponseWriter, r *http.Request) {

}
