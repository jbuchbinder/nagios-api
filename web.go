package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func buildRouter() *mux.Router {
	r := mux.NewRouter()

	// Status file parsing

	r_status := r.PathPrefix("/status").Subrouter()

	r_hosts := r_status.PathPrefix("/hosts").Subrouter()
	r_hosts.HandleFunc("/", listHostHandler).Methods("GET")
	r_hosts.HandleFunc("/{host}", hostHandler).Methods("GET")

	r_services := r_status.PathPrefix("/services").Subrouter()
	r_services.HandleFunc("/{host}", listServiceHandler).Methods("GET")
	r_services.HandleFunc("/{host}/{service}", serviceHandler).Methods("GET")
	r_services.HandleFunc("/{host}/{service}/state", serviceStateHandler).Methods("GET")

	return r
}

func writeError(w http.ResponseWriter, e error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(501)
	w.Write([]byte(fmt.Sprintf("{\"error\":\"%s\"}", e)))
}

func writeOutput(w http.ResponseWriter, obj interface{}) {
	b, e := json.Marshal(obj)
	if e != nil {
		writeError(w, e)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(b)
}
