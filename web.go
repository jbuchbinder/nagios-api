package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func buildRouter() *mux.Router {
	r := mux.NewRouter()

	r_hosts := r.PathPrefix("/hosts").Subrouter()
	r_hosts.HandleFunc("/", listHostHandler).Methods("GET")
	r_hosts.HandleFunc("/{host}", hostHandler).Methods("GET")

	r_services := r.PathPrefix("/services").Subrouter()
	r_services.HandleFunc("/{host}", listServiceHandler).Methods("GET")
	r_services.HandleFunc("/{host}/{service}", serviceHandler).Methods("GET")
	r_services.HandleFunc("/{host}/{service}/state", serviceStateHandler).Methods("GET")

	return r
}

func listHostHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("listHostHandler: %s", r.RemoteAddr)
	keys := make([]string, len(status.HostStatus))
	iter := 0
	for host := range status.HostStatus {
		keys[iter] = host
		iter++
	}
	writeOutput(w, keys)
}

func hostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("hostHandler: %s: %s", r.RemoteAddr, vars["host"])
	writeOutput(w, status.HostStatus[vars["host"]])
}

func listServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("listServiceHandler: %s: %s", r.RemoteAddr, vars["host"])
	h := status.HostStatus[vars["host"]]
	keys := make([]string, len(h.ServiceStatus))
	iter := 0
	for svc := range h.ServiceStatus {
		keys[iter] = svc
		iter++
	}
	writeOutput(w, keys)
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("serviceHandler: %s: %s/%s", r.RemoteAddr, vars["host"], vars["service"])
	h := status.HostStatus[vars["host"]]
	writeOutput(w, h.ServiceStatus[vars["service"]])
}

func serviceStateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("serviceStateHandler: %s: %s/%s", r.RemoteAddr, vars["host"], vars["service"])
	h := status.HostStatus[vars["host"]]
	i, e := strconv.Atoi(h.ServiceStatus[vars["service"]].Current_state)
	if e != nil {
		writeError(w, e)
		return
	}
	writeOutput(w, i)
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
