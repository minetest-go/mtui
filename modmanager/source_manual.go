package modmanager

import (
	"errors"
	"fmt"
	"mtui/types"
	"os"
)

// ManualModHandler handles mods without an upstream source (for example legacy contentdb mods)
type ManualModHandler struct{}

func (h *ManualModHandler) Create(world_dir string, mod *types.Mod) error {
	return errors.New("manual mods can't be installed")
}

func (h *ManualModHandler) Update(world_dir string, mod *types.Mod, version string) error {
	return errors.New("manual mods can't be updated")
}

func (h *ManualModHandler) Remove(world_dir string, mod *types.Mod) error {
	if mod.Name == "" {
		return nil
	}
	dir := getDir(world_dir, mod)
	if dir == "" {
		return nil
	}
	return os.RemoveAll(dir)
}

func (h *ManualModHandler) CheckUpdate(world_dir string, mod *types.Mod) (bool, error) {
	return false, nil
}

func getHandler(source types.SourceType) (SourceTypeHandler, error) {
	h, ok := modHandlers[source]
	if !ok {
		return nil, fmt.Errorf("unsupported source type: '%s'", source)
	}
	return h, nil
}
