package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

var (
	fileToWatch *string
	pidFile     *string
)

func init() {
	fileToWatch = flag.String("fileToWatch", "config.conf", "the file to watch for changes")
	pidFile = flag.String("pidFile", "pidfile.pid", "the file containing the PID of the process to receive the SIGHUP")
}

func main() {
	flag.Parse()

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
				if event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Rename == fsnotify.Rename {
					// if the file was removed or renamed, linux inotify
					// removes the watch, so re-add it
					err = watcher.Add(*fileToWatch)
					if err != nil {
						log.Fatal(err)
					}
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Printf("change detected (%s)\n", event.Name)
					err := sendHup()
					if err != nil {
						log.Printf("failed to send SIGHUP: %v\n", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(*fileToWatch)
	if err != nil {
		log.Fatalf("%s: %v\n", *fileToWatch, err)
	}
	log.Printf("watching (%s)\n", *fileToWatch)
	<-done
}

func sendHup() error {
	content, err := ioutil.ReadFile(*pidFile)
	if err != nil {
		return err
	}

	pid, err := strconv.Atoi(strings.TrimSuffix(string(content), "\n"))
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Signal(syscall.SIGHUP)
	if err != nil {
		return err
	}
	log.Printf("sent SIGHUP to %d (%s)\n", pid, *pidFile)

	return nil
}
