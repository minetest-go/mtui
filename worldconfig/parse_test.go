package worldconfig

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSqlite(t *testing.T) {
	cfg, err := Parse("./testdata/world.mt.sqlite")
	assert.NoError(t, err)
	if cfg[CONFIG_AUTH_BACKEND] != BACKEND_SQLITE3 {
		t.Fatal("not sqlite3")
	}
}

func TestParsePostgres(t *testing.T) {
	cfg, err := Parse("./testdata/world.mt.postgres")
	assert.NoError(t, err)
	fmt.Println(cfg)
	if cfg[CONFIG_AUTH_BACKEND] != BACKEND_POSTGRES {
		t.Fatal("not postgres")
	}

	if cfg[CONFIG_PSQL_AUTH_CONNECTION] != "host=/var/run/postgresql user=postgres password=enter dbname=postgres" {
		t.Fatal("param err")
	}
}
