package main

import (
	"flag"
	ns "github.com/jbuchbinder/nagiosstatus"
	"log"
	"net/http"
)

var (
	statusFilePath = flag.String("statusfile", "/var/log/nagios/status.dat", "Path to Nagios status.dat file")
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

	http.Handle("/", buildRouter())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
