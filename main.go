package main

import (
	"fmt"
	"mtadmin/app"
	"mtadmin/web"
	"net/http"
	"os"
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

	fmt.Printf("Listening on port %d\n", 8080)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
