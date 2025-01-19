package modmanager

import (
	"mtui/db"
	"mtui/types"
)

var modHandlers = map[types.SourceType]SourceTypeHandler{
	types.SourceTypeGIT: &GitModHandler{},
	types.SourceTypeCDB: &ContentDBModHandler{},
}

type ModManager struct {
	world_dir string
}

func New(world_dir string, repo *db.ModRepository) *ModManager {
	return &ModManager{
		world_dir: world_dir,
	}
}

func (m *ModManager) Create(mod *types.Mod) error {
	handler := modHandlers[mod.SourceType]
	return handler.Create(m.world_dir, mod)
}

func (m *ModManager) Update(mod *types.Mod, version string) error {
	handler := modHandlers[mod.SourceType]
	return handler.Update(m.world_dir, mod, version)
}

func (m *ModManager) Remove(mod *types.Mod) error {
	handler := modHandlers[mod.SourceType]
	return handler.Remove(m.world_dir, mod)
}
