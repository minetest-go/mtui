package db

import "database/sql"

const PLAYERS_MIGRATE = `
pragma journal_mode = wal;
CREATE TABLE if not exists player (name VARCHAR(50) NOT NULL,pitch NUMERIC(11, 4) NOT NULL,yaw NUMERIC(11, 4) NOT NULL,posX NUMERIC(11, 4) NOT NULL,posY NUMERIC(11, 4) NOT NULL,posZ NUMERIC(11, 4) NOT NULL,hp INT NOT NULL,breath INT NOT NULL,creation_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,modification_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,PRIMARY KEY (name));
CREATE TABLE if not exists player_metadata (player VARCHAR(50) NOT NULL,metadata VARCHAR(256) NOT NULL,value TEXT,PRIMARY KEY(player, metadata),FOREIGN KEY (player) REFERENCES player (name) ON DELETE CASCADE );
CREATE TABLE if not exists player_inventories (layer VARCHAR(50) NOT NULL,inv_id INT NOT NULL,inv_width INT NOT NULL,inv_name TEXT NOT NULL DEFAULT '',inv_size INT NOT NULL,PRIMARY KEY(player, inv_id),   FOREIGN KEY (player) REFERENCES player (name) ON DELETE CASCADE );
CREATE TABLE if not exists player_inventory_items (   player VARCHAR(50) NOT NULL,inv_id INT NOT NULL,slot_id INT NOT NULL,item TEXT NOT NULL DEFAULT '',PRIMARY KEY(player, inv_id, slot_id),   FOREIGN KEY (player) REFERENCES player (name) ON DELETE CASCADE );
`

func NewPlayerRepository(filename string) (*PlayerRepository, error) {
	db_, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}

	_, err = db_.Exec(PLAYERS_MIGRATE)
	if err != nil {
		return nil, err
	}

	return &PlayerRepository{db: db_}, nil
}

type PlayerRepository struct {
	db *sql.DB
}
