package main

import (
	"mtadmin/db"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	world_dir := os.Getenv("WORLD_DIR")

	_, err := db.NewAuthRepository(world_dir + "/auth.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = db.NewMapRepository(world_dir + "/map.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = db.NewModStorageRepository(world_dir + "/mod_storage.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = db.NewPlayerRepository(world_dir + "/players.sqlite")
	if err != nil {
		panic(err)
	}

}
