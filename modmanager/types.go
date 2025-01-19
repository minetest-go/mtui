package modmanager

import (
	"mtui/db"
	"mtui/types"
)

type ModStatus struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}

type HandlerContext struct {
	WorldDir string
	Repo     *db.ModRepository
}

type SourceTypeHandler interface {
	Create(world_dir string, mod *types.Mod) error
	Update(world_dir string, mod *types.Mod, version string) error
	Remove(world_dir string, mod *types.Mod) error
	CheckUpdate(world_dir string, mod *types.Mod) (bool, error)
}
