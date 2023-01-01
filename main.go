package main

import (
	"context"
	"flag"
	"fmt"
	"mtui/app"
	"mtui/events"
	"mtui/jobs"
	"mtui/web"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {

	version_ptr := flag.Bool("version", false, "show the current version and exit")
	flag.Parse()

	if *version_ptr {
		v := app.Version
		if v == "" {
			v = "DEV"
		}
		fmt.Printf("mtui %s\n", v)
		return
	}

	if os.Getenv("LOGLEVEL") == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var err error
	world_dir := os.Getenv("WORLD_DIR")
	if world_dir == "" {
		// fall back to current directory
		world_dir, err = os.Getwd()
		if err != nil {
			panic(err)
		}
	}

	a, err := app.Create(world_dir)
	if err != nil {
		panic(err)
	}

	logrus.WithFields(logrus.Fields{
		"version": app.Version,
	}).Info("Starting mtui")

	err = web.Setup(a)
	if err != nil {
		panic(err)
	}

	err = events.Setup(a)
	if err != nil {
		panic(err)
	}

	// start jobs
	jobs.Start(a)

	server := &http.Server{Addr: ":8080", Handler: nil}

	go func() {
		fmt.Printf("Listening on port %d\n", 8080)
		err = server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	var captureSignal = make(chan os.Signal, 1)
	signal.Notify(captureSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-captureSignal
	server.Shutdown(context.Background())
}
