package main

import (
	ns "github.com/jbuchbinder/nagiosstatus"
	"gopkg.in/fsnotify.v1"
	"log"
)

func loadStatusFile() {
	log.Print("loadStatusFile()")
	nsp := ns.NagiosStatusParser{}
	nsp.SetStatusFile(*statusFilePath)
	status = nsp.Parse()
}

type watcher struct {
	StatusFile string

	w       *fsnotify.Watcher
	running bool
}

func (self *watcher) Start() error {
	log.Print("watcher.Start()")
	self.running = true

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Print(err)
		return err
	}

	// Persist to object
	self.w = w

	log.Print("watcher: Spinning up event loop")
	go func() {
		log.Print("watcher: Event loop started")
		for {
			if !self.running {
				break
			}

			select {
			case ev := <-self.w.Events:
				log.Println("watcher: Event: ", ev)
				if ev.Op != 0 {
					self.ProcessEvent(&ev)
				}
				break
			case <-self.w.Errors:
				log.Println("watcher: Error: ", err)
				break
			}
		}
	}()

	err = self.w.Add(self.StatusFile)
	if err != nil {
		self.w.Close()
		return err
	}

	return nil
}

func (self *watcher) opToText(op fsnotify.Op) string {
	switch op {
	case fsnotify.Create:
		return "create"
	case fsnotify.Write:
		return "write"
	case fsnotify.Remove:
		return "remove"
	case fsnotify.Rename:
		return "rename"
	case fsnotify.Chmod:
		return "chmod"
	default:
		return "unknown"
	}
}

func (self *watcher) ProcessEvent(event *fsnotify.Event) {
	log.Printf("watcher.ProcessEvent %s : %s", event.Name, self.opToText(event.Op))
	loadStatusFile()
}

func (self *watcher) Stop() {
	log.Print("watcher.Stop()")
	if self.w != nil {
		log.Print("watcher: Closing")
		self.w.Close()
	}
}
