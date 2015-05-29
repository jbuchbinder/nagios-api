package main

import (
	"flag"
	"fmt"
	ns "github.com/jbuchbinder/nagiosstatus"
	"log"
	"net/http"
	"time"
)

var (
	statusFilePath = flag.String("statusfile", "/var/log/nagios/status.dat", "Path to Nagios status.dat file")
	port           = flag.Int("port", 8888, "Listening port")
	status         *ns.NagiosStatus
	watch          *watcher
)

func main() {
	flag.Parse()

	log.Print("Starting service")

	log.Printf("Initial load of %s", *statusFilePath)
	loadStatusFile()

	log.Print("Starting watcher thread")
	watch = &watcher{StatusFile: *statusFilePath}
	err := watch.Start()
	if err != nil {
		panic(err)
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", *port),
		Handler:        buildRouter(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Starting HTTP server on port %d", *port)
	log.Fatal(s.ListenAndServe())
}
