package main

import (
	"context"
	"flag"
	"fmt"
	"mtui/app"
	"mtui/auth"
	"mtui/events"
	"mtui/jobs"
	"mtui/web"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	dbauth "github.com/minetest-go/mtdb/auth"

	"github.com/sirupsen/logrus"
)

var version_flag = flag.Bool("version", false, "show the current version and exit")
var help_flag = flag.Bool("help", false, "shows all options")
var create_admin = flag.String("create_admin", "", "creates a new admin-user (format: 'user:pass')")

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

	a, err := app.Create(world_dir)
	if err != nil {
		panic(err)
	}

	if *create_admin != "" {
		parts := strings.Split(*create_admin, ":")
		if len(parts) != 2 {
			panic("invalid format")
		}
		username := parts[0]
		password := parts[1]
		salt, verifier, err := auth.CreateAuth(username, password)
		if err != nil {
			panic(err)
		}

		str := auth.CreateDBPassword(salt, verifier)
		entry, err := a.DBContext.Auth.GetByUsername(username)
		if err != nil {
			panic(err)
		}

		entry = &dbauth.AuthEntry{
			Name:     username,
			Password: str,
		}
		err = a.DBContext.Auth.Create(entry)
		if err != nil {
			panic(err)
		}

		for _, priv := range []string{"server", "interact", "privs"} {
			err = a.DBContext.Privs.Create(&dbauth.PrivilegeEntry{
				ID:        *entry.ID,
				Privilege: priv,
			})

			if err != nil {
				panic(err)
			}
		}
		fmt.Printf("Created a new admin-user '%s'\n", username)
		return
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
