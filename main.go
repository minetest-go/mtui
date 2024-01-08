package main

import (
	"context"
	"flag"
	"fmt"
	"mtui/app"
	"mtui/events"
	"mtui/jobs"
	"mtui/types"
	"mtui/web"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

var version_flag = flag.Bool("version", false, "show the current version and exit")
var help_flag = flag.Bool("help", false, "shows all options")

func main() {
	flag.Parse()

	if *help_flag {
		flag.PrintDefaults()
		return
	}

	if *version_flag {
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

	cfg := types.NewConfig(world_dir)
	a, err := app.Create(cfg)
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
		logrus.WithFields(logrus.Fields{"port": 8080}).Info("Listening")
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
