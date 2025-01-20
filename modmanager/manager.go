package modmanager

import (
	"fmt"
	"mtui/db"
	"mtui/types"
)

var modHandlers = map[types.SourceType]SourceTypeHandler{
	types.SourceTypeGIT: &GitModHandler{},
	types.SourceTypeCDB: &ContentDBModHandler{},
}

type ModManager struct {
	world_dir string
	repo      *db.ModRepository
}

func New(world_dir string, repo *db.ModRepository) *ModManager {
	return &ModManager{
		world_dir: world_dir,
		repo:      repo,
	}
}

func (m *ModManager) Create(mod *types.Mod) error {
	handler := modHandlers[mod.SourceType]
	err := handler.Create(m.world_dir, mod)
	if err != nil {
		return fmt.Errorf("error creating mod '%s': %v", mod.Name, err)
	}
	err = m.repo.Create(mod)
	if err != nil {
		return fmt.Errorf("error persisting mod '%s': %v", mod.Name, err)
	}
	return nil
}

func (m *ModManager) Update(mod *types.Mod, version string) error {
	handler := modHandlers[mod.SourceType]
	err := handler.Update(m.world_dir, mod, version)
	if err != nil {
		return fmt.Errorf("error updating mod '%s': %v", mod.Name, err)
	}
	err = m.repo.Update(mod)
	if err != nil {
		return fmt.Errorf("error persisting updated mod '%s': %v", mod.Name, err)
	}
	return nil
}

func (m *ModManager) Remove(mod *types.Mod) error {
	handler := modHandlers[mod.SourceType]
	err := handler.Remove(m.world_dir, mod)
	if err != nil {
		return fmt.Errorf("error remooving mod '%s': %v", mod.Name, err)
	}
	err = m.repo.Delete(mod.ID)
	if err != nil {
		return fmt.Errorf("error deleting mod entity '%s': %v", mod.Name, err)
	}
	return nil
}

func (m *ModManager) CheckUpdates() error {
	mods, err := m.repo.GetAll()
	if err != nil {
		return fmt.Errorf("error fetching mods from db: %v", err)
	}

	for _, mod := range mods {
		h := modHandlers[mod.SourceType]
		updated, err := h.CheckUpdate(m.world_dir, mod)
		if err != nil {
			return fmt.Errorf("update check failed for mod '%s': %v", mod.Name, err)
		}
		if updated {
			err = m.repo.Update(mod)
			if err != nil {
				return fmt.Errorf("error persiting update for mod '%s': %v", mod.Name, err)
			}
		}
	}

	return nil
}
