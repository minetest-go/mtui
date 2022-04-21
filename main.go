package main

import (
	"database/sql"
	"mtadmin/worldconfig"
	"os"

	_ "modernc.org/sqlite"
)

const AUTH_MIGRATE = `
pragma journal_mode = wal;
CREATE TABLE if not exists auth (id INTEGER PRIMARY KEY AUTOINCREMENT,name VARCHAR(32) UNIQUE,password VARCHAR(512),last_login INTEGER);
CREATE TABLE if not exists user_privileges (id INTEGER,privilege VARCHAR(32),PRIMARY KEY (id, privilege)CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES auth (id) ON DELETE CASCADE);
`

const PLAYERS_MIGRATE = `
pragma journal_mode = wal;
CREATE TABLE if not exists player (name VARCHAR(50) NOT NULL,pitch NUMERIC(11, 4) NOT NULL,yaw NUMERIC(11, 4) NOT NULL,posX NUMERIC(11, 4) NOT NULL,posY NUMERIC(11, 4) NOT NULL,posZ NUMERIC(11, 4) NOT NULL,hp INT NOT NULL,breath INT NOT NULL,creation_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,modification_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,PRIMARY KEY (name));
CREATE TABLE if not exists player_metadata (player VARCHAR(50) NOT NULL,metadata VARCHAR(256) NOT NULL,value TEXT,PRIMARY KEY(player, metadata),FOREIGN KEY (player) REFERENCES player (name) ON DELETE CASCADE );
CREATE TABLE if not exists player_inventories (layer VARCHAR(50) NOT NULL,inv_id INT NOT NULL,inv_width INT NOT NULL,inv_name TEXT NOT NULL DEFAULT '',inv_size INT NOT NULL,PRIMARY KEY(player, inv_id),   FOREIGN KEY (player) REFERENCES player (name) ON DELETE CASCADE );
CREATE TABLE if not exists player_inventory_items (   player VARCHAR(50) NOT NULL,inv_id INT NOT NULL,slot_id INT NOT NULL,item TEXT NOT NULL DEFAULT '',PRIMARY KEY(player, inv_id, slot_id),   FOREIGN KEY (player) REFERENCES player (name) ON DELETE CASCADE );
`

const MOD_STORAGE_MIGRATE = `
pragma journal_mode = wal;
CREATE TABLE if not exists entries (
	modname TEXT NOT NULL,
	key BLOB NOT NULL,
	value BLOB NOT NULL,
	PRIMARY KEY (modname, key)
);
`

const MAP_MIGRATE = `
pragma journal_mode = wal;
CREATE TABLE if not exists blocks (
	pos INT PRIMARY KEY,
	data BLOB
);
`

func main() {
	world_dir := os.Getenv("WORLD_DIR")

	cfg, err := worldconfig.Parse(world_dir + "/world.mt")
	if err != nil {
		panic(err)
	}

	var db_ *sql.DB
	auth_backend := cfg[worldconfig.CONFIG_AUTH_BACKEND]
	switch auth_backend {
	case worldconfig.BACKEND_SQLITE3:
		db_, err = sql.Open("sqlite", world_dir+"/auth.sqlite")
		if err != nil {
			panic(err)
		}
	default:
		panic("unsupported backend: " + auth_backend)
	}

	_, err = db_.Exec(AUTH_MIGRATE)
	if err != nil {
		panic(err)
	}

}
