package main

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func buildRouter() *mux.Router {
	r := mux.NewRouter()

	// Command parsing

	r_cmd := r.PathPrefix("/cmd").Subrouter()
	for k := range nagiosCommands {
		// Skip if we haven't defined a name
		if nagiosCommands[k].Name == "" {
			continue
		}
		ncmd := nagiosCommands[k]
		r_cmd.HandleFunc("/"+ncmd.Name, func(w http.ResponseWriter, r *http.Request) {
			// Pull post body
			postraw, err := ioutil.ReadAll(r.Body)
			if err != nil {
				writeError(w, err)
				return
			}

			// Deserialize into default object
			var post map[string]interface{}
			err = json.Unmarshal(postraw, &post)
			if err != nil {
				writeError(w, err)
				return
			}

			// Check to see whether or not required params are present
			req := ncmd.Required
			for k2 := range req {
				_, present := post[req[k2]]
				if !present {
					writeError(w, errors.New("Missing required parameter "+req[k2]))
					return
				}
			}

			// Populate defaults
			args := map[string]interface{}{}
			for k3 := range ncmd.Defaults {
				args[k3] = ncmd.Defaults[k3]
			}

			// Populate vars
			for k4 := range post {
				args[k4] = post[k4]
			}

			// Execute
			cmd, err := FormCommand(k, args)
			if err != nil {
				writeError(w, err)
				return
			}
			err = writeCommand(cmd)
			if err != nil {
				writeError(w, err)
				return
			}
		}).Methods("POST")
	}

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

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, e error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(501)
	b, _ := json.Marshal(errorResponse{
		Error: e.Error(),
	})
	w.Write(b)
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
