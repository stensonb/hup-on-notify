package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

/*
const fileToWatch = "../hup-catcher/hup-catcher.conf"
const pidFile = "../hup-catcher/hup-catcher.pid"
*/
const fileToWatch = "/etc/squid/squid.conf"
const pidFile = "/var/run/squid.pid"

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					sendHup()
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(fileToWatch)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("watching %s\n", fileToWatch)
	<-done
}

func sendHup() error {
	content, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return err
	}

	pid, err := strconv.Atoi(string(content))
	if err != nil {
		return err
	}

	err = syscall.Kill(pid, syscall.SIGHUP)
	if err != nil {
		return err
	}
	log.Printf("sent SIGHUP to %d\n", pid)

	return nil
}
