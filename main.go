package main

import (
	"fmt"
	"mtadmin/app"
	"mtadmin/mod"
	"mtadmin/web"
	"net/http"
	"os"
)

func main() {
	world_dir := os.Getenv("WORLD_DIR")

	a, err := app.Create(world_dir)
	if err != nil {
		panic(err)
	}

	err = web.Setup(a)
	if err != nil {
		panic(err)
	}

	err = mod.Install(world_dir + "/worldmods/mtadmin")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on port %d\n", 8080)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
