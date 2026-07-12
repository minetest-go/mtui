package types

import (
	"reflect"
	"testing"
)

func TestMinetestCommand(t *testing.T) {
	want := []string{
		"--world", "/world",
		"--config", "/minetest.conf",
		"--info",
		"--logfile", "/logs/luanti.log",
	}
	cfg := Config{DockerMinetestLogLevel: "info", DockerMinetestLogfile: "/logs/luanti.log"}
	if got := cfg.MinetestCommand(); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected command: %v", got)
	}
}

func TestMinetestCommandDefaultsLogfileToWorld(t *testing.T) {
	want := []string{
		"--world", "/world",
		"--config", "/minetest.conf",
		"--logfile", "/world/debug.txt",
	}
	cfg := Config{DockerMinetestLogLevel: "action"}
	if got := cfg.MinetestCommand(); !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected command: %v", got)
	}
}
