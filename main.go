package main

import (
	"context"
	"fmt"
	"mtui/app"
	"mtui/web"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
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

	err = web.Setup(a)
	if err != nil {
		panic(err)
	}

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
