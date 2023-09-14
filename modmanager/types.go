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
	Create(ctx *HandlerContext, mod *types.Mod) error
	Update(ctx *HandlerContext, mod *types.Mod, version string) error
	Remove(ctx *HandlerContext, mod *types.Mod) error
	CheckUpdate(ctx *HandlerContext, mod *types.Mod) (bool, error)
}
