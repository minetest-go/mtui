package main

import (
	"fmt"
	"mtadmin/db"
	"mtadmin/web"
	"net/http"
	"os"
)

func main() {
	world_dir := os.Getenv("WORLD_DIR")

	repos, err := db.CreateRepositories(world_dir)
	if err != nil {
		panic(err)
	}

	err = web.Setup(repos)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on port %d\n", 8080)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}
